package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/utils"
	"github.com/sirjager/trueauth/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TrueAuthService) Register(ctx context.Context, req *rpc.RegisterRequest) (*rpc.RegisterResponse, error) {
	violations := validateRegisterRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err.Error())
	}

	// extract metadata like client-ip and user-agent
	meta := s.extractMetadata(ctx)

	params := sqlc.CreateAccountParams{
		Email:      req.GetEmail(),
		Username:   req.GetUsername(),
		Password:   hashedPassword,
		Firstname:  req.GetFirstname(),
		Lastname:   req.GetLastname(),
		AllowedIps: []string{meta.ClientIp},
	}
	account, err := s.store.CreateAccount(ctx, params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s already exists", unique_violation(pqErr)))
			}
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &rpc.RegisterResponse{Account: publicProfile(account)}, nil
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

func unique_violation(err *pq.Error) string {
	s := strings.Split(err.Detail, "=")[0]
	return strings.Split(strings.Split(s, "(")[1], ")")[0]
}
