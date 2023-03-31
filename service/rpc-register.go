package service

import (
	"context"
	"strings"

	"github.com/lib/pq"
	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/validator/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s *TrueAuthService) Register(ctx context.Context, req *rpc.RegisterRequest) (*rpc.RegisterResponse, error) {
	violations := validateRegisterRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	return &rpc.RegisterResponse{}, nil
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
