package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// type AuthServiceServer interface {
// 	Login(context.Context, *LoginRequest) (*LoginResponse, error)
// 	mustEmbedUnimplementedAuthServiceServer()
// }

type Service struct {
	OpenIDResolver OpenIDResolver
	Logger         zap.Logger
	authpb.UnimplementedAuthServiceServer
}

//OpenIDResolver resolves an authorization code to an open id.
type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("received code", zap.String("code", req.Code))

	openID,err:=s.OpenIDResolver.Resolve((req.Code))
	if err != nil {
		return nil,status.Errorf(codes.Unavailable,"cannot resolve openid: %v",err)
	}

	return &authpb.LoginResponse{
		AccessToken: "token for " + openID,
		ExpiresIn:   7200,
	}, nil

}
