package auth

import (
	"fmt"
	"math"
	"time"

	"github.com/Witthaya22/golang-backend-itctc/config"
	"github.com/Witthaya22/golang-backend-itctc/entities"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type auth struct {
	mapClaims *authMapClaims
	conf      config.Jwt
}

type authMapClaims struct {
	Claims *entities.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type IAuth interface {
	SingToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *auth) SingToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString([]byte(a.conf.SecretKey))
	return ss
}

func NewAuth(tokenType TokenType, conf config.Jwt, claims *entities.UserClaims) (IAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(conf, claims), nil
	case Refresh:
		return newRefreshToken(conf, claims), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func newAccessToken(conf config.Jwt, claims *entities.UserClaims) IAuth {
	return &auth{
		conf: conf,
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:  "backend-itctc",
				Subject: "access-token",
				Audience: []string{
					"admin",
					"user",
				},
				ExpiresAt: jwtTimeDurationCal(int(conf.AccessExpires)),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(conf config.Jwt, claims *entities.UserClaims) IAuth {
	return &auth{
		conf: conf,
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:  "backend-itctc",
				Subject: "refresh-token",
				Audience: []string{
					"admin",
					"user",
				},
				ExpiresAt: jwtTimeDurationCal(int(conf.RefreshExpires)),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}
