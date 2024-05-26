package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/sirjager/trueauth/pkg/mail"
)

const TaskSendUserDeletionCode = "task:sendEmailUserDeletion"

type PayloadUserDeletionCode struct {
	Token string `json:"token"`
}

func (d *RedisTaskDistributor) DistributeTaskSendUserDeletionCode(
	ctx context.Context,
	payload PayloadUserDeletionCode,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}
	task := asynq.NewTask(TaskSendUserDeletionCode, jsonPayload, opts...)
	if _, err := d.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}
	d.logger.Info().Str("task", TaskSendUserDeletionCode).Msg("task enqueued")
	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendUserDeletionCode(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadUserDeletionCode
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// NOTE: if token is invalid, we wont even retry the task
	tokenPayload, err := p.tokens.VerifyToken(payload.Token)
	if err != nil {
		return fmt.Errorf(err.Error(), asynq.SkipRetry)
	}

	email := mail.Mail{To: []string{tokenPayload.Payload.UserEmail}}
	email.Subject = "Account Deletion Requested"
	email.Body = fmt.Sprintf(`
	User account deletion requested.<br><br>
	Use the code to complete the process.<br>
	Client IP  : <b>%s</b> <br>
	User Agent : <b>%s</b> <br>
	Use Code   : <b>%s</b> <br>
	Valid Till : <b>%s</b> <br><br>
	<b>If you did not make this request, please ignore this email.</b>`,
		tokenPayload.Payload.ClientIP,
		tokenPayload.Payload.UserAgent,
		tokenPayload.Payload.Code,
		tokenPayload.ExpiresAt.Format("2006-01-02 15:04:05"),
	)

	if err = p.mailer.SendMail(email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	p.logger.Info().Str("task", TaskSendUserDeletionCode).Msg("task processed successfully")
	return nil
}
