package auth

import (
	"errors"
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

func ParseToken(tokenString string, conf config.Jwt) (*authMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.SecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("malformed token: %v", err)
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired: %v", err)
		} else {
			return nil, fmt.Errorf("invalid token: %v", err)
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token claims: %v", err)
	}
}

func RepeatToken(conf config.Jwt, claims *entities.UserClaims, exp int64) string {
	obj := &auth{
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
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SingToken()
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
