package server

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/hibiken/asynq"
	"google.golang.org/grpc/status"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
	rpc "github.com/sirjager/trueauth/stubs"
	"github.com/sirjager/trueauth/worker"
)

const (
	verifyCodeDigitsCount = 8
	errInvaidCode         = "invalid code"
)

func (s *Server) Verify(ctx context.Context, req *rpc.VerifyRequest) (*rpc.VerifyResponse, error) {
	if len(req.GetCode()) > 0 && (len(req.GetCode()) != verifyCodeDigitsCount) {
		return nil, status.Errorf(_aborted, errInvaidCode)
	}

	user, err := s.store.ReadUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(_notFound, errEmailNotRegistered)
		}
		return nil, status.Errorf(_internal, err.Error())
	}

	profile, err := user.Profile()
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	if user.Verified {
		return &rpc.VerifyResponse{User: publicProfile(profile)}, nil
	}

	if len(req.GetCode()) == 0 {
		if time.Since(user.LastEmailVerify) < s.config.Auth.VerifyTokenCooldown {
			tryAfter := s.config.Auth.VerifyTokenCooldown - time.Since(user.LastEmailVerify)
			return nil, status.Errorf(
				_aborted,
				"email verification has been requested recently, please try again after %s",
				tryAfter,
			)
		}

		code := utils.RandomNumberAsString(verifyCodeDigitsCount)
		tokenParams := tokens.PayloadData{UserEmail: user.Email, UserID: user.ID, Code: code}
		token, _, tokenErr := s.tokens.CreateToken(tokenParams, s.config.Auth.VerifyTokenExpDur)
		if tokenErr != nil {
			return nil, status.Errorf(_internal, "failed to create token, %s", tokenErr.Error())
		}

		taskParams := worker.PayloadEmailVerification{Token: token}
		taskOptions := []asynq.Option{
			asynq.MaxRetry(5),
			asynq.Group(worker.QueueLow),
			asynq.ProcessIn(time.Millisecond * time.Duration(utils.RandomInt(100, 600))),
		}
		if err = s.tasks.DistributeTaskSendEmailVerification(ctx, taskParams, taskOptions...); err != nil {
			return nil, status.Errorf(_internal, "send email verification failed, %s", err.Error())
		}

		updateParam := db.UpdateUserTokenEmailVerifyParams{
			ID:               user.ID,
			LastEmailVerify:  time.Now(),
			TokenEmailVerify: token,
		}
		if uErr := s.store.UpdateUserTokenEmailVerify(ctx, updateParam); uErr != nil {
			return nil, status.Errorf(_internal, uErr.Error())
		}

		return &rpc.VerifyResponse{Message: "check your inbox for further instructions"}, nil
	}

	tokenPayload, err := s.tokens.VerifyToken(user.TokenEmailVerify)
	if err != nil {
		return nil, status.Errorf(_aborted, "invalid code: %s", err.Error())
	}

	if !bytes.Equal(tokenPayload.Payload.UserID, user.ID) {
		return nil, status.Errorf(_aborted, errInvaidCode)
	}

	if tokenPayload.Payload.Code != req.GetCode() {
		return nil, status.Errorf(_aborted, errInvaidCode)
	}

	if tokenPayload.Payload.UserEmail != user.Email {
		return nil, status.Errorf(_aborted, errInvaidCode)
	}

	params := db.UpdateUserVerifiedParams{Verified: true, TokenEmailVerify: "", ID: user.ID}
	user, err = s.store.UpdateUserVerified(ctx, params)
	if err != nil {
		return nil, status.Errorf(_internal, "failed to update verified status, %s", err.Error())
	}

	profile, err = user.Profile()
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	return &rpc.VerifyResponse{User: publicProfile(profile)}, nil
}
