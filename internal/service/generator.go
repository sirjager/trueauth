package service

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sirjager/trueauth/internal/db/sqlc"

	rpc "github.com/sirjager/trueauth/stubs/go"
)

func publicProfile(user sqlc.User) *rpc.User {
	_user := &rpc.User{
		Id:        string(user.ID),
		Email:     user.Email,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Verified:  user.EmailVerified,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
	return _user
}
