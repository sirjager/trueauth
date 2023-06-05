package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirjager/trueauth/internal/db/sqlc"
	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/utils"
)

const TASK_SEND_EMAIL_VERIFICATION = "task:send-email-verification"

type PayloadSendEmailVerification struct {
	CreateUserParams    sqlc.CreateUserParams `json:"create_user_params"`
	ClientIP            string                `json:"clientip"`
	UserAgent           string                `json:"useragent"`
	VerificationCode    string                `json:"verification_code"`
	VerificationCodeTTL time.Duration         `json:"verification_code_ttl"`
	CreatedAt           time.Time             `json:"payload_created_at"`
}

func (d *RedisTaskDistributor) DistributeTaskSendEmailVerification(ctx context.Context, payload PayloadSendEmailVerification, opts ...asynq.Option) (err error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}
	task := asynq.NewTask(TASK_SEND_EMAIL_VERIFICATION, jsonPayload, opts...)
	if _, err = d.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}
	d.logr.Info().Str("task", TASK_SEND_EMAIL_VERIFICATION).Msg("task enqueued")
	return
}

func (p *RedisTaskProcessor) ProcessTaskSendEmailVerification(ctx context.Context, task *asynq.Task) (err error) {
	var payload PayloadSendEmailVerification
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	//
	payload.VerificationCode = utils.UUID_XID()
	verificationCodeTimeToLive := payload.VerificationCodeTTL

	mail := mail.Mail{To: []string{payload.CreateUserParams.Email}}
	mail.Subject = "Complete Registration Process"
	template := p.emailTemplate

	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_HEADING, "Please confirm<br>your email address")
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_SUBHEADING, "Thank your for joining us !")
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_BODY, fmt.Sprintf(`Verify your email to complete registration.<br>
	It ensures accurate contact info, enhances account security, and enables communication.<br>
	Get important alerts and reset password if needed.<br><br>
	Request from<br>IP Address: <strong>%s</strong><br>UserAgent: <strong>%s</strong><br><br>
	Verification Code is only valid for %s<br>
	`, payload.ClientIP, payload.UserAgent, verificationCodeTimeToLive.String()))

	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_ACTION, payload.VerificationCode)
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_TEAMNAME, "TrueAuth")

	mail.Body = template

	if err := p.mailer.SendMail(mail); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if err = p.redis.Set(ctx, payload, verificationCodeTimeToLive, utils.PendingRegistrationKey(payload.CreateUserParams.Email)); err != nil {
		return err
	}

	p.logr.Info().Str("task", TASK_SEND_EMAIL_VERIFICATION).Msg("task processed successfully")

	return
}
