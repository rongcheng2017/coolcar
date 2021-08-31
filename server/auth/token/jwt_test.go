package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQBdEyX1Pk1qD9XI6saaFU+rzYgW/OuZEaPqCxJqNyLD4eiFCtyu
8diYmJz0fZ8kGXowV4QRSKguCavcNhS6CJy0l+68MVdtv14faK9jqrH7T0ClCrmr
8vcYz8un1s0QjGZWVOD7kpM3+P60GGuTABKOrl9UC2AYSymgOggXibn0PxIj5ocQ
kPIaXuF9apuauit5tNQUeHR43ifBi8oQ0L26jKlKyOrBJIB0bK/HBErBJ/rdW4/C
YUN44q2k8D77yLqk6MKNmYxwsUy81tynI7yAXfXNXN/m1epBloOTov01FvyXr1ls
YhUwYGbLjCN5QXijroUHMQEFvn9dT8CgSeifAgMBAAECggEAHoXNCOPpl0KiJUV4
48bhXcIFQySmTohObM48t8BlEj/fdOHfBTAii8hgkH0x1zDTtU697L0bCh350jma
CngQi9jXXbgAp+j+AObfvZuYyoVu+bDOuujux6A9dUkn7qlcVit3rIig5tYtEPqB
LK/1Zf0hHfqtBMqBWB4v0ShFZydyGVwRqgRU6AjbFgvJx5x9wI0Hq9vtzpsOMERg
yvIl1gCFEvFkiROmTAxfq8hQt0pvsDKCNL8VRKhzmyjYg3KQkwPAZiDJ/9D5lUGH
QmSrdFiX83FJAYHG/6ebTpFF88jOIVWypX9mg3Jm4nlNQYpYueYVCV+VybjLRQBZ
3W0FwQKBgQCqfEIyx1HXovwooiAfr9GO31BPTahGxRajkZyq4DYyQvvkhCRcu/DM
/8gPKSnuNQbU8TsM+CtpZaYcaGtRRfQMVzLWD/3H7BKf+nIQ6Zo4Da8TJPbbuRAx
bHREsLNCviKqQd6lWG3QWGUONh3/frON0eFihQwcjjGXKVtTv6c55wKBgQCLwrn4
1Eziw6Uyv0nqFNt+MVk5V826j2pcE+NLCCpISPCkhrh6EzUjSJYxsho/lWIBiZHU
zAFq68P5G8gEezzVf0Iarz22sY/XIywvzBHnolQFiLMN0U6BsILJv2LJUUgPUekM
8lbxYdt9sHFuaj2PShK1yZrrsY1S9JcyQI80iQKBgGB1XZ8NVykCdlknIbXL7G1B
vFaiQYuJB34UbOfhY8icTZjFiy1MyLm0HqU1TRwRtIPW2OpFn4pKkOmRyuZ5BdPV
olWrRpNO5lrNgKxA/5inZV8XkvROiPLtwfr7XvFsUoCyNB6pIbi3yrV3uRFNxpl/
Hl53mJqveS9lnt6LmToRAoGAPo1E1w2N6+BMy825sz7qjixgFr4podoWbGeqTya0
Ze3fZoO1hU2bdtNCBbQE83hUiQOddXRpHgWvjIrWlsrhi1yNpYvRPzdxfYSMfkgD
q3yHxoJMQV7wmDL8FnfGKvxqGBE9EUJVj2uQ5UxXOGfsbXllrl8xK1QoQHygPymN
7qECgYEApAxqISrqrwG9ltZ+Va5A6nE5wlphr88yKKfb6iIQbiCDQwe5R2KsVwX4
vlFr6r5EZT3FXPyczyEz46OT38YSfRLeuTYKRPi26t77meHFrSdcGVwbF2nOJybw
QTlZiwQbcj5AJ9q/pk6VSB2wSrBll2qK2PK1QD8TmWyngcKKSCM=
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
	want:=`eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.ME8AQfq8lh1efZrhUFxZmBbokSy-DYZzuzQwI_uHpYMq_NXCU4nq6MObydUJpEk2Ki2KgZZzZI6RMx5BD5njagd7e44nzeY-NzQMlO5OhLlUwu-buFhCxtwtAFH5YRa4DWz-6cSKTmxg5FZ-S6bB-fgl4GZOu8drGAoyk9r9di8pqLhlIMYgTQ48RG6eWtSMbiNiRiNS7Kh_uyMgjKl6Lp1f-RcZs6It_WUeer4czxIUBPh7lEuuZr1pcjcHscROJctax2h-ToGcp3OLOeNgeV9jPqoifS6Cwjmw4jnrylaE9a0Zp02Z8vMYPHq5LWSvisHBQOcR3P6eVIDcfkg8xg`
	if tkn!=want {
		t.Errorf("wrong token generated. want:%q; got:%q",want,tkn)
	}
}
