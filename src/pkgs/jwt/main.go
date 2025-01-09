package jwt

import (
	"github.com/felix-Asante/pennyPilot-go-api/src/utils"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWTService struct {
	JWTAuth *jwtauth.JWTAuth
}

func NewJWTService() *JWTService {
	return &JWTService{
		JWTAuth: InitAuthToken(),
	}
}
func InitAuthToken() *jwtauth.JWTAuth {
	jwtAuth := jwtauth.New("HS256", []byte(utils.GetEnv("JWT_SECRET")), nil)
	return jwtAuth
}

func (s *JWTService) GetJWTAuth() *jwtauth.JWTAuth {
	return s.JWTAuth
}

func (s *JWTService) Encode(data map[string]interface{}) (jwt.Token, string, error) {
	return s.JWTAuth.Encode(data)
}
