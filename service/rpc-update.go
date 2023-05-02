package service

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/db/sqlc"
	"github.com/sirjager/trueauth/utils"
	"github.com/sirjager/trueauth/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TrueAuthService) Update(ctx context.Context, req *rpc.UpdateRequest) (*rpc.UpdateResponse, error) {
	violations := validateUpdateRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}

	shouldUpdate := false

	updateParams := sqlc.Update_UserParams{
		ID:        authorized.User.ID,
		Username:  authorized.User.Username,
		Password:  authorized.User.Password,
		Firstname: authorized.User.Firstname,
		Lastname:  authorized.User.Lastname,
	}

	if req.GetFirstname() != "" {
		updateParams.Firstname = req.GetFirstname()
		shouldUpdate = true
	}
	if req.GetLastname() != "" {
		updateParams.Lastname = req.GetLastname()
		shouldUpdate = true
	}
	if req.GetUsername() != "" {
		updateParams.Username = req.GetUsername()
		shouldUpdate = true
	}
	if req.GetPassword() != "" {
		hashedPassword, err := utils.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err.Error())
		}
		updateParams.Password = hashedPassword
		shouldUpdate = true
	}

	if shouldUpdate {
		authorized.User, err = s.store.Update_User(ctx, updateParams)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "unique_violation":
					return nil, status.Errorf(codes.Internal, fmt.Sprintf("%s already exists", unique_violation(pqErr)))
				}
			}
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return &rpc.UpdateResponse{Account: publicProfile(authorized.User)}, nil
	}

	return &rpc.UpdateResponse{Account: publicProfile(authorized.User)}, nil
}

func validateUpdateRequest(req *rpc.UpdateRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.GetPassword() != "" {
		if err := validator.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	if req.GetUsername() != "" {
		if err := validator.ValidateUsername(req.GetUsername()); err != nil {
			violations = append(violations, fieldViolation("username", err))
		}
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

	return
}
