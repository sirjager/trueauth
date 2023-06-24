package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rpc "github.com/sirjager/trueauth/stubs/go"

	"github.com/sirjager/trueauth/internal/db/sqlc"
	"github.com/sirjager/trueauth/internal/worker"

	"github.com/sirjager/trueauth/pkg/crypt"
	"github.com/sirjager/trueauth/pkg/utils"
	"github.com/sirjager/trueauth/pkg/validator"
)

func (s *CoreService) Register(ctx context.Context, req *rpc.RegisterRequest) (*rpc.RegisterResponse, error) {
	violations := validateRegisterRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	bytes, err := s.redis.Get(ctx, utils.PendingRegistrationKey(req.GetEmail()))
	if err == nil {
		var payload worker.PayloadSendEmailVerification
		if err = json.Unmarshal(bytes, &payload); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		tryAfter := time.Until(payload.CreatedAt.Add(s.Config.VerifyTokenCooldown))
		return &rpc.RegisterResponse{Message: fmt.Sprintf("complete pending registration using email verification code sent to your email or start fresh registration after %s", tryAfter)}, nil
	}

	// check if user already exists
	user, err := s.store.ReadUserByIdentity(ctx, sqlc.ReadUserByIdentityParams{Email: req.GetEmail(), Username: req.GetUsername()})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.Internal, err.Error())
	}

  

	if user.Email == req.GetEmail() {
		return nil, status.Errorf(codes.AlreadyExists, "email already exists")
	}
	if user.Username == req.Username {
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	}

	hashedPassword, err := crypt.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err.Error())
	}
	// extract metadata like client-ip and user-agent
	meta := s.extractMetadata(ctx)
	if err != nil && err != redis.Nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// this means no pending registraion, and we can start new registration
	// Here there will no pending, means no code has been expired or record for this email does not exists
	// now we will save create user params in redis, and send email verification code to user's email
	// and return message: verify email to complete registration process

	// first we send email verification code to user's email
	// and then we save create user params in redis
	// so when user will use email verification code, we will fetch create params and just create user with email verified status to true
	createUserParams := sqlc.CreateUserParams{
		ID:        []byte(utils.UUID_XID()),
		Email:     req.GetEmail(),
		Username:  req.GetUsername(),
		Password:  hashedPassword,
		Firstname: req.GetFirstname(),
		Lastname:  req.GetLastname(),
		// AllowedIps:    []string{meta.ClientIP},
		EmailVerified: false,
	}

	payload := worker.PayloadSendEmailVerification{
		CreateUserParams:    createUserParams,
		ClientIP:            meta.ClientIP,
		UserAgent:           meta.UserAgent,
		VerificationCodeTTL: s.Config.VerifyTokenTTL,
		CreatedAt:           time.Now(),
	}

	opts := []asynq.Option{asynq.Group(worker.QUEUE_CRITICAL), asynq.MaxRetry(3), asynq.ProcessIn(time.Second)}
	if err = s.taskDistributor.DistributeTaskSendEmailVerification(ctx, payload, opts...); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &rpc.RegisterResponse{Message: "check your email to complete registration"}, nil
}

func validateRegisterRequest(req *rpc.RegisterRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if req.GetFirstname() != "" {
		if err := validator.ValidateName(req.GetFirstname()); err != nil {
			violations = append(violations, fieldViolation("firstname", err))
		}
	}
	if req.GetLastname() != "" {
		if err := validator.ValidateName(req.GetLastname()); err != nil {
			violations = append(violations, fieldViolation("lastname", err))
		}
	}

	return violations
}
