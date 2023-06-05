package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/sirjager/trueauth/pkg/utils"
)

const TASK_SEND_EMAIL_ALLOW_IP = "task:send-email-allowip"

type PayloadSendEmailAllowIP struct {
	Email           string    `json:"email"`
	AllowIP         string    `json:"allowip"`
	UserAgent       string    `json:"useragent"`
	Timestamp       time.Time `json:"timestamp"`
	Code            string    `json:"allowip_code"`
	LastEmailSentAt time.Time `json:"last_email_sent_at"`
}

func (d *RedisTaskDistributor) DistributeTaskSendEmailAllowIP(ctx context.Context, payload PayloadSendEmailAllowIP, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}
	task := asynq.NewTask(TASK_SEND_EMAIL_ALLOW_IP, jsonPayload, opts...)
	_, err = d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}
	d.logr.Info().Str("task", TASK_SEND_EMAIL_ALLOW_IP).Msg("task enqueued")
	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendEmailAllowIP(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendEmailAllowIP
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	payload.Code = utils.UUID_XID()
	payload.LastEmailSentAt = time.Now()
	expiration := p.config.AllowIPTokenTTL

	mail := mail.Mail{To: []string{payload.Email}}
	mail.Subject = "Security Alert: Unrecognized Login Request from Unknown IP Address"
	template := p.emailTemplate
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_HEADING, "Security Alert<br>Unrecognized Login")
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_SUBHEADING, "")
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_BODY, fmt.Sprintf(`We have detected a login attempt on your account<br>from an unknown IP address.<br><br>
	<b>Details of the login attempt</b>: <br>
	<b>Timestamp</b><br>%s<br>
	<b>IP Address</b><br>%s<br>
	<b>User Agent</b><br>%s<br><br>
	<strong>If this login attempt was unauthorized or not initiated by you, please take the following steps immediately to secure your account:</strong><br>
	1. Change your account password. <br><br>
	<strong>If you believe this login attempt was made by you, and you would like to allow access from this IP address, please use the following code to authorize it:</strong><br>
	Code is only valid for %s`,
		payload.Timestamp.Local().String(), payload.AllowIP, payload.UserAgent, expiration.String()))

	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_ACTION, payload.Code)
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_TEAMNAME, "TrueAuth")

	mail.Body = template

	if err := p.mailer.SendMail(mail); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if err := p.redis.Set(ctx, payload, expiration, utils.AllowIPKey(payload.Email, payload.AllowIP)); err != nil {
		return err
	}

	return nil
}
