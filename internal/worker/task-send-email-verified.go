package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/sirjager/trueauth/pkg/mail"
)

const TASK_SEND_EMAIL_VERIFIED = "task:send-email-verified"

type PayloadSendEmailVerified struct {
	Email string `json:"email"`
}

func (d *RedisTaskDistributor) DistributeTaskSendEmailVerified(ctx context.Context, payload PayloadSendEmailVerified, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}

	task := asynq.NewTask(TASK_SEND_EMAIL_VERIFIED, jsonPayload, opts...)
	_, err = d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	d.logr.Info().Str("task", TASK_SEND_EMAIL_VERIFIED).Msg("task enqueued")

	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendEmailVerified(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendEmailVerified
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	email := mail.Mail{To: []string{payload.Email}}
	email.Subject = "Registration Complete!"

	template := p.emailTemplate

	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_HEADING, "Registration Complete!")
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_SUBHEADING, "Welcome to our community")
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_BODY, `Thank you for joining us!<br><br>
	Should you require any assistance or have any questions, please feel free to reach out to our support team. 
	We are here to help!<br><br> Once again, thank you for choosing to be a part of our community.<br>
	We are excited to have you on board.`)
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_ACTION, fmt.Sprintf("%s Verified", payload.Email))
	template = strings.ReplaceAll(template, EMAIL_TEMPLATE_TEAMNAME, "TrueAuth")
	email.Body = template
	if err := p.mailer.SendMail(email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	p.logr.Info().Str("task", TASK_SEND_EMAIL_VERIFIED).Msg("task processed successfully")

	return nil
}
