package service

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sirjager/trueauth/db/sqlc"

	rpc "github.com/sirjager/rpcs/trueauth/go"
)

func publicProfile(user sqlc.User) *rpc.Account {
	_user := &rpc.Account{
		Id:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
	return _user
}
