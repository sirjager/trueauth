package server

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sirjager/trueauth/db/db"
	rpc "github.com/sirjager/trueauth/rpc"
)

func publicProfile(user db.User) *rpc.User {
	return &rpc.User{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Verified:  user.Verified,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
