package server

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/hash"
	"github.com/sirjager/trueauth/pkg/utils"
	"github.com/sirjager/trueauth/pkg/validator"
	rpc "github.com/sirjager/trueauth/rpc"
)

func (s *Server) Signup(ctx context.Context, req *rpc.SignupRequest) (*rpc.SignupResponse, error) {
	// returns invalid requests
	if violations := validateSignupRequest(req); violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	hashingSalt := hash.RandomSalt(len(req.GetPassword()))
	hashedPassword, err := s.hasher.Hash(hashingSalt, req.GetPassword())
	if err != nil {
		return nil, status.Errorf(_internal, "failed to hash password: %s", err.Error())
	}

	createUserParams := db.CreateUserParams{
		ID:        utils.XIDNew().Bytes(),
		Email:     req.GetEmail(),
		Username:  req.GetUsername(),
		HashSalt:  hashingSalt,
		HashPass:  hashedPassword,
		Firstname: req.GetFirstname(),
		Lastname:  req.GetLastname(),
	}

	user, err := s.store.CreateUser(ctx, createUserParams)
	if err != nil {
		if db.ErrorCode(err) == db.ErrUniqueViolation.Code {
			return nil, status.Errorf(_conflict, uniqueViolationError(err))
		}
		return nil, status.Errorf(_internal, err.Error())
	}

	profile, err := user.Profile()
	if err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	return &rpc.SignupResponse{User: publicProfile(profile)}, nil
}

func validateSignupRequest(
	req *rpc.SignupRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {
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

func uniqueViolationError(err error) string {
	key := strings.Split(err.Error(), "_users_")[1]
	key = strings.Split(key, "_key")[0]
	return fmt.Errorf("%s already exists", key).Error()
}
