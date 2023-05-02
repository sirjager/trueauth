package service

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	err_invalid_identity = "invalid identity, it can either be email, username or id(uuid)"
)

func unAuthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err.Error())
}

func unknownIPError() error {
	return status.Error(codes.PermissionDenied, "activity from unknown ip address is strictly prohibited")
}

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentsError(violations []*errdetails.BadRequest_FieldViolation) error {
	invalidReq := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")
	statusDetails, err := statusInvalid.WithDetails(invalidReq)
	if err != nil {
		return statusInvalid.Err()
	}
	return statusDetails.Err()
}
