package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"

	"go.uber.org/zap"
)

// type AuthServiceServer interface {
// 	Login(context.Context, *LoginRequest) (*LoginResponse, error)
// 	mustEmbedUnimplementedAuthServiceServer()
// }

type Service struct {
	Logger zap.Logger
	authpb.UnimplementedAuthServiceServer
}

func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse,error) {
	s.Logger.Info("received code", zap.String("code", req.Code))

	return &authpb.LoginResponse{
		AccessToken: "token for "+req.Code,
		ExpiresIn: 7200,
	}, nil

}
