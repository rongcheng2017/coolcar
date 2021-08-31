package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const PublicKey = `-----BEGIN PUBLIC KEY-----
MIIBITANBgkqhkiG9w0BAQEFAAOCAQ4AMIIBCQKCAQBdEyX1Pk1qD9XI6saaFU+r
zYgW/OuZEaPqCxJqNyLD4eiFCtyu8diYmJz0fZ8kGXowV4QRSKguCavcNhS6CJy0
l+68MVdtv14faK9jqrH7T0ClCrmr8vcYz8un1s0QjGZWVOD7kpM3+P60GGuTABKO
rl9UC2AYSymgOggXibn0PxIj5ocQkPIaXuF9apuauit5tNQUeHR43ifBi8oQ0L26
jKlKyOrBJIB0bK/HBErBJ/rdW4/CYUN44q2k8D77yLqk6MKNmYxwsUy81tynI7yA
XfXNXN/m1epBloOTov01FvyXr1lsYhUwYGbLjCN5QXijroUHMQEFvn9dT8CgSeif
AgMBAAE=
-----END PUBLIC KEY-----`

func TestVerigy(t *testing.T) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(PublicKey))
	if err != nil {
		t.Fatalf("cannot parse public key: %v", err)
	}

	v := JWTTokenVerifier{
		PublicKey: pubKey,
	}

	cases := []struct {
		name    string
		tkn     string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name:    "valid_token",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.ME8AQfq8lh1efZrhUFxZmBbokSy-DYZzuzQwI_uHpYMq_NXCU4nq6MObydUJpEk2Ki2KgZZzZI6RMx5BD5njagd7e44nzeY-NzQMlO5OhLlUwu-buFhCxtwtAFH5YRa4DWz-6cSKTmxg5FZ-S6bB-fgl4GZOu8drGAoyk9r9di8pqLhlIMYgTQ48RG6eWtSMbiNiRiNS7Kh_uyMgjKl6Lp1f-RcZs6It_WUeer4czxIUBPh7lEuuZr1pcjcHscROJctax2h-ToGcp3OLOeNgeV9jPqoifS6Cwjmw4jnrylaE9a0Zp02Z8vMYPHq5LWSvisHBQOcR3P6eVIDcfkg8xg",
			now:     time.Unix(1516239122, 0),
			want:    "1234567890",
			wantErr: false,
		},
		{
			name:    "token_expired",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.ME8AQfq8lh1efZrhUFxZmBbokSy-DYZzuzQwI_uHpYMq_NXCU4nq6MObydUJpEk2Ki2KgZZzZI6RMx5BD5njagd7e44nzeY-NzQMlO5OhLlUwu-buFhCxtwtAFH5YRa4DWz-6cSKTmxg5FZ-S6bB-fgl4GZOu8drGAoyk9r9di8pqLhlIMYgTQ48RG6eWtSMbiNiRiNS7Kh_uyMgjKl6Lp1f-RcZs6It_WUeer4czxIUBPh7lEuuZr1pcjcHscROJctax2h-ToGcp3OLOeNgeV9jPqoifS6Cwjmw4jnrylaE9a0Zp02Z8vMYPHq5LWSvisHBQOcR3P6eVIDcfkg8xg",
			now:     time.Unix(1517239122, 0),
			want:    "1234567890",
			wantErr: true,
		},
		{
			name:    "bad_token",
			tkn:     "bad_token",
			now:     time.Unix(1516239122, 0),
			want:    "1234567890",
			wantErr: true,
		},
		{
			name:    "wrong_signature",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.",
			now:     time.Unix(1516239122, 0),
			want:    "1234567890",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			//修改test now time
			jwt.TimeFunc = func() time.Time {
				return c.now
			}
			accountID, err := v.Verify(c.tkn)
			// if err != nil {
			// 	t.Errorf("verification failed:%v", err)
			// }

			if !c.wantErr && err != nil {
				t.Errorf("verification failed:%v", err)
			}
			if c.wantErr && err == nil {
				t.Errorf("want error; got no error")
			}

			if accountID != c.want {
				t.Errorf("wrong account id, want:%q,got: %q", c.want, accountID)
			}

		})

	}

}
