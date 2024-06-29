package server

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirjager/gopkg/utils"
	"google.golang.org/grpc/status"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/internal/tokens"
	rpc "github.com/sirjager/trueauth/rpc"
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

	dbuser, err := s.store.ReadUserByEmail(ctx, req.GetEmail())
	if err != nil {
		// NOTE: Here we can also return ErrEmailNotRegistered, but
		// we dont want to disclose if user exists or not, so we simply return email sent.
		if errors.Is(err, db.ErrRecordNotFound) {
			return &rpc.VerifyResponse{Message: "check your inbox for further instructions"}, nil
		}
		return nil, status.Errorf(_internal, err.Error())
	}

	if dbuser.Verified {
		return &rpc.VerifyResponse{User: publicProfile(dbuser)}, nil
	}

	meta := s.extractMetadata(ctx)

	if len(req.GetCode()) == 0 {
		if time.Since(dbuser.LastEmailVerify) < s.config.Auth.VerifyTokenCooldown {
			tryAfter := s.config.Auth.VerifyTokenCooldown - time.Since(dbuser.LastEmailVerify)
			return nil, status.Errorf(
				_aborted,
				"email verification has been requested recently, please try again after %s",
				tryAfter,
			)
		}

		code := utils.RandomNumberAsString(verifyCodeDigitsCount)
		tokenParams := tokens.PayloadData{
			Code:      code,
			UserID:    dbuser.ID,
			UserEmail: dbuser.Email,
			ClientIP:  meta.clientIP(),
			UserAgent: meta.userAgent(),
		}
		token, _, tokenErr := s.tokens.CreateToken(tokenParams, s.config.Auth.VerifyTokenExpDur)
		if tokenErr != nil {
			return nil, status.Errorf(_internal, "failed to create token, %s", tokenErr.Error())
		}

		taskParams := worker.PayloadEmailVerificationCode{Token: token}
		taskOptions := []asynq.Option{
			asynq.MaxRetry(5),
			asynq.Group(worker.QueueLow),
			asynq.ProcessIn(time.Millisecond * time.Duration(utils.RandomInt(100, 600))),
		}
		if err = s.tasks.DistributeTaskSendEmailVerificationCode(ctx, taskParams, taskOptions...); err != nil {
			return nil, status.Errorf(_internal, "send email verification failed, %s", err.Error())
		}

		updateParam := db.UpdateUserEmailVerificationTokenParams{
			ID:               dbuser.ID,
			LastEmailVerify:  time.Now(),
			TokenEmailVerify: token,
		}
		if uErr := s.store.UpdateUserEmailVerificationToken(ctx, updateParam); uErr != nil {
			return nil, status.Errorf(_internal, uErr.Error())
		}

		return &rpc.VerifyResponse{Message: "check your inbox for further instructions"}, nil
	}

	tokenPayload, err := s.tokens.VerifyToken(dbuser.TokenEmailVerify)
	if err != nil {
		return nil, status.Errorf(_aborted, "invalid code: %s", err.Error())
	}

	if !bytes.Equal(tokenPayload.Payload.UserID, dbuser.ID) {
		return nil, status.Errorf(_aborted, errInvaidCode)
	}

	if tokenPayload.Payload.Code != req.GetCode() {
		return nil, status.Errorf(_aborted, errInvaidCode)
	}

	if tokenPayload.Payload.UserEmail != dbuser.Email {
		return nil, status.Errorf(_aborted, errInvaidCode)
	}

	params := db.UpdateUserEmailVerifiedParams{Verified: true, TokenEmailVerify: "", ID: dbuser.ID}
	dbuser, err = s.store.UpdateUserEmailVerified(ctx, params)
	if err != nil {
		return nil, status.Errorf(_internal, "failed to update verified status, %s", err.Error())
	}

	return &rpc.VerifyResponse{User: publicProfile(dbuser)}, nil
}
