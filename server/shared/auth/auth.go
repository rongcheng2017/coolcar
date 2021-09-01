package auth

import (
	"context"
	"coolcar/shared/auth/token"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	bearerPrefix        = "Bearer "
)

type tokenVerifier interface {
	Verify(token string) (string, error)
}

type interceptor struct {
	verifier tokenVerifier
}

// Interceptor creates a grpc auth interceptor
func Inteceptor(publicKeyFile string) (grpc.UnaryServerInterceptor, error) {
	f, err := os.Open(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open public key file: %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key: %v", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}
	i := &interceptor{
		verifier: &token.JWTTokenVerifier{
			PublicKey: pubKey,
		},
	}
	return i.HandlerReq, nil
}

func (i *interceptor) HandlerReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	tkn, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	aid, err := i.verifier.Verify((tkn))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token not vaild %v", err)
	}

	return handler(ContextWithAccountID(ctx,AccountID(aid)), req)
}
func tokenFromContext(c context.Context) (string, error) {
	m, ok := metadata.FromIncomingContext(c)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}
	tkn := ""
	for _, v := range m[authorizationHeader] {
		if strings.HasPrefix(v, bearerPrefix) {
			tkn = v[len(bearerPrefix):]
		}
	}
	if tkn == "" {
		return "", status.Error(codes.Unauthenticated, "")
	}
	return tkn, nil
}

type accountIDKey struct{}

//AccountID defines account id object.
type AccountID string

func (a AccountID) String() string{
	return string(a)
}

func ContextWithAccountID(c context.Context, aid AccountID) context.Context {
	return context.WithValue(c, accountIDKey{}, aid)
}

// Returns unauthenticated error if no account id is available
func AccountIDFromContext(c context.Context) (AccountID, error) {
	v := c.Value(accountIDKey{})
	aid,ok := v.(AccountID)
	if !ok{
		return "",status.Error(codes.Unauthenticated,"")
	}
	return aid,nil
}