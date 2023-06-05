package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sirjager/trueauth/pkg/utils"
)

const TASK_CLEAR_COMPLETED_VERIFICATIONS = "task:clear-completed-verifications"

type PayloadClearCompletedVerfications struct {
	Email string `json:"email"`
}

func (d *RedisTaskDistributor) DistributeTaskClearCompletedVerifications(ctx context.Context, payload PayloadClearCompletedVerfications, opts ...asynq.Option) (err error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}
	task := asynq.NewTask(TASK_CLEAR_COMPLETED_VERIFICATIONS, jsonPayload, opts...)
	if _, err = d.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}
	d.logr.Info().Str("task", TASK_CLEAR_COMPLETED_VERIFICATIONS).Msg("task enqueued")
	return
}

func (p *RedisTaskProcessor) ProcessTaskClearCompletedVerifications(ctx context.Context, task *asynq.Task) (err error) {
	var payload PayloadClearCompletedVerfications
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	if err = p.redis.Del(ctx, utils.PendingRegistrationKey(payload.Email)); err != nil {
		return err
	}
	p.logr.Info().Str("task", TASK_CLEAR_COMPLETED_VERIFICATIONS).Msg("task processed successfully")

	return
}
