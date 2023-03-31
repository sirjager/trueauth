package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/utils"
)

const TASK_SEND_VERIFY_EMAIL = "task:send-verify-email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}

	task := asynq.NewTask(TASK_SEND_VERIFY_EMAIL, jsonPayload, opts...)
	_, err = distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	distributor.logger.Info().
		Str("task", TASK_SEND_VERIFY_EMAIL).
		Msg("task enqueued")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		// if errors.Is(err, sql.ErrNoRows) {
		// 	return fmt.Errorf("user does not exists: %w", asynq.SkipRetry)
		// }
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	// generate a random code
	sixDigitCode := utils.RandomNumberAsString(6)
	codeExpiresAt := time.Now().Add(time.Minute * 10)
	emailRecordParams := sqlc.CreateEmailRecordParams{
		UserID:        user.ID,
		Email:         user.Email,
		Verified:      user.Verified,
		Code:          sixDigitCode,
		CodeExpiresAt: codeExpiresAt,
	}

	// send email
	emailEntry, err := processor.store.CreateEmailRecord(ctx, emailRecordParams)
	if err != nil {
		return fmt.Errorf("failed to create email entry: %w", err)
	}

	//TODO: send an actual email
	processor.logger.Info().
		Str("task", TASK_SEND_VERIFY_EMAIL).
		Str("email", emailEntry.Email).
		Msg("task processed successfully")

	return nil
}
