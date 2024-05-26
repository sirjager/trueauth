package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hibiken/asynq"

	"github.com/sirjager/trueauth/pkg/mail"
)

const TaskSendEmailVerification = "task:sendEmailVerification"

type PayloadEmailVerification struct {
	Token string `json:"token"`
}

func (d *RedisTaskDistributor) DistributeTaskSendEmailVerification(
	ctx context.Context,
	payload PayloadEmailVerification,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}
	task := asynq.NewTask(TaskSendEmailVerification, jsonPayload, opts...)
	if _, err := d.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}
	d.logger.Info().Str("task", TaskSendEmailVerification).Msg("task enqueued")
	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendEmailVerification(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadEmailVerification
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// NOTE: if token is invalid, we wont even retry the task
	tokenPayload, err := p.tokens.VerifyToken(payload.Token)
	if err != nil {
		return fmt.Errorf(err.Error(), asynq.SkipRetry)
	}

	email := mail.Mail{To: []string{tokenPayload.Payload.UserEmail}}
	email.Subject = "Complete Your Registration"
	templ := p.config.Mail.Template
	templ = strings.ReplaceAll(templ, TemplateHeading, "Welcome to our community")
	templ = strings.ReplaceAll(templ, TemplateSubHeading, "Complete Your Registration")
	callback := fmt.Sprintf("%s?email=%s&code=%s",
		p.config.Auth.CallbackURL, tokenPayload.Payload.UserEmail, tokenPayload.Payload.Code)

	message := fmt.Sprintf(`Thank you for joining us!<br><br>
	To complete your registration and unlock the full features of our platform,<br>
	please click on the link below to verify your email address.<br><br>
	Verifcation link: <a href="%s">Cick here to verify</a><br><br>
	Verification link expires at: <b>%s</b><br>
	after which you may need to request a new verification link.<br><br>
	Should you require any assistance or have any questions, please feel free to reach out to our support team.
	We are here to help!<br><br>
	Once again, thank you for choosing to be a part of our community.<br>
	We are excited to have you on board.`, callback, tokenPayload.ExpiresAt.Format("2006-01-02 15:04:05"))
	templ = strings.ReplaceAll(templ, TemplateBody, message)
	templ = strings.ReplaceAll(templ, TemplateAppName, p.config.Server.AppName)
	templ = strings.ReplaceAll(templ, TemplateAction, "Verify Email")
	templ = strings.ReplaceAll(templ, TemplateActionLink, callback)
	email.Body = templ

	if err = p.mailer.SendMail(email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	p.logger.Info().Str("task", TaskSendEmailVerification).Msg("task processed successfully")
	return nil
}
