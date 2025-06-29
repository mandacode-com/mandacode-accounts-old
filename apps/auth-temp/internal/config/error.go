package config

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InvalidConfigError(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}
