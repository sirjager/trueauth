package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sirjager/trueauth/mail"
)

const TASK_SEND_EMAIL_VERIFIED = "task:send-verify-email"

type PayloadSendEmailVerified struct {
	Email string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendEmailVerified(ctx context.Context, payload PayloadSendEmailVerified, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}

	task := asynq.NewTask(TASK_SEND_EMAIL_VERIFIED, jsonPayload, opts...)
	_, err = distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	distributor.logger.Info().Str("task", TASK_SEND_EMAIL_VERIFIED).Msg("task enqueued")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendEmailVerified(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendEmailVerified
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	//TODO: send an actual email
	email := mail.Mail{To: []string{payload.Email}}
	email.Subject = "Email verified"
	email.Body = fmt.Sprintf(`Hello <br/>
	Your email <b>%s</b> is successfully verified. <br/><br/>
	Thank You`, payload.Email)

	if err := processor.mailer.SendMail(email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	processor.logger.Info().Str("task", TASK_SEND_EMAIL_VERIFIED).Msg("task processed successfully")

	return nil
}
