package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAnLc646JkqSxHSxKaqqTDUZbYnijnFAdL0CbyRUnM4s5XeSiP
Bezr+7STzjkwCuGfRDvvpg1K4T9Nn6lthYD15AELQ5OhH9/Ze/PKuQxDDMtucFGB
nTJWxRg7MxewnCuQ6n2Cx739L+43JEMqlFNLn+fF0OB+Rqg/fLzLWbh70YRmlfDM
EcDMC87EVEsmQC//l1jkOrHQVflcilaWPfZbCgDmmlUxfOrMeYbXGx8Xs6VBXo4n
CdpemC/qE3WirHJpeltYEdN+/HBIDyMJss7JyZB0kIKwCcqV1bIGgV2rnTjrFL2u
LTXVSUS9Mgxq62KrlHuQtoiJlpmv8QIW2BODfQIDAQABAoIBACa3gd4BHbtJzCc4
5msoH3UFvmh8lHI3RzyakpoZnHugHK47HfKJ0NczipyVNlBJ424ZHKC6gfhaClRa
qEsmTBlTRLQiQdk9FV7xIPFMnTgI/jTTfiEW8abp0J5TyPccMEYuUeQFBQsVYBwe
V+OjHHjZ6t4qyCeuo1iuz9JPqI9fMhCz0r2/X/6Q7flyQ/F7r49PguywLhLkCbtE
zE/rvZwpYQvg9++UuTfIVxQsigA5PMuFDz/n0XyyCcGJ+Ch8EU17GILiTnFAkqSO
iOKUehsqDRUcU5el3BLgtm2O1RXCqGd71Gmyoq4Fu0FrahvfXfX3pWxIEKCApj0p
zz5PvmUCgYEA262QElLc0/zwBuK6cqDaTmUH3myZ1owTHVnqKme0qwn3THP/Kct6
jOcp+tLoay2ixm1alLRvegHy1EhYsJeORGslOAUmhGPsgkUBXT7uIZtNs1INC/Oi
xBrakuLbogNIIfYAvW9iUbsFSU7JODliGz/dFcaxvurz63jfzYjw6HsCgYEAtqCg
zTNGeidtGJLgMcjMNnAjKpQw8go/2uzWfWBFreojm74bqGfj0OKMcFCRZWiAiwas
6Str+kmf8h01nbCUvnejcb8f16z5aPWsRKJXU0u28qbOrcd8xwusD4ukGcBFZ6Q2
g8cfWsCRjl6ENJ6oUuYGizElyUEbYCpnsPtJzmcCgYEAjdytl0evl65WCvxLz06U
699Oj5KuXeCjT2cLU0sZXwLWkqat9v2SLH/zmiitMtmLrnxb7IABJVcwy2nU7GVS
2Fgg9uZMk148E3wgf2juOwGh0dWA22EAkYeN8yFRGHTqFhRZMfxGD+WoakjYpNhZ
xKMfULq5ekMcNcofLQnsGRUCgYBfeyuXHS/Dvck0B9ZfMPRTod1A7amJYgJwm1Ko
yiSkAL4NNx+OtIJPO6LhNb5OnoxWI29TmPgjK0sMcmkNwLyDuFkjpyEmybC8R3WB
jL7LNdK4mq2D/cAm8NtMZV2ueO/Qd/JogzrJX9S58oB8YlbuwIS7UT7IMdn2NTVx
OnAkEQKBgBo1L9sDa9xCuM2INOecXJqR8boQFWmNPzj2605bSEdZJJjBO2JNYWEe
GUTGX8Ab7+x61FRyN1xk0oey/3ubmfVkZ2D3BXdylLT9npqwCgGhuhRD9nto9vED
cqM5YANpEt0aCfr7EzvD6CYfa1alsbgm5xd40k0adAOzJOdA4YGG
-----END RSA PRIVATE KEY-----`

func TestGenerateToken(t *testing.T) {
	key,err:=jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("cannot parse private key:%v",err)
	}
	g := NewJWTTokenGen("coolcar/auth", key)
	g.nowFunc = func() time.Time {
		return time.Unix(1516239022,0)
	}
	tkn,err:=g.GenerateToken("1234567890",2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate token: %v",err)
	}
	want:=`eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.Fyth6ojy_3oXeEDOd-D2PisHH2sd13xxOGBIjySHoLacHK7z8FFs_BYtNyhbBS_e-3EqTgBdAE_3pp1JJBTtzYGqLvgu6aMs1RdkNZuU2g1OmV8QoS1qlt_S84QrexHj_3nwHtlgJg_ZLbW1pWeuQUCtgY5h801bjKgkQlzx458dMnMHFAumG5uWUF6cb5LOKUk1dKhUwVsX4fnYAT7C5E42M92KZZBqPtkVT38HW5SuZNCmAU35VFxxefYH4gWOIssx9VoWstH4YorwGxbN25pM2TWvm9blWfKHfhn-PUIo5PxKmNQYveFZfuTbqxdRQy2IOc394gdNP5hce8LtDA`
	if tkn!=want {
		t.Errorf("wrong token generated. want:%q; got:%q",want,tkn)
	}
}
