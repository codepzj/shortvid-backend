package service

import (
	"errors"
	"shortvid-backend/app/shortvid-service/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	jwtConf *conf.Jwt
	logger  log.Logger
}

func NewJwtService(jwtConf *conf.Jwt, logger log.Logger) *JwtService {
	return &JwtService{jwtConf: jwtConf, logger: logger}
}

type JwtCustomClaims struct {
	UserUID   int32
	SessionID string
	jwt.RegisteredClaims
}

// GenerateAccessToken 创建access_token
func (s *JwtService) GenerateAccessToken(userUID int32, sessionID string) (string, error) {
	expiresAt := time.Now().Add(s.GetTokenExpiration())
	claims := JwtCustomClaims{
		UserUID:   userUID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.jwtConf.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConf.SecretKey))
}

// GenerateRefreshToken 创建refresh_token
func (s *JwtService) GenerateRefreshToken(userUID int32, sessionID string) (string, error) {
	expiresAt := time.Now().Add(s.GetRefreshTokenExpiration())
	claims := JwtCustomClaims{
		UserUID:   userUID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.jwtConf.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConf.SecretKey))
}

// ValidateToken 将token解析为claims
func (s *JwtService) ValidateToken(jwtStr string) (*JwtCustomClaims, error) {
	claims := new(JwtCustomClaims)
	token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (any, error) {
		return []byte(s.jwtConf.SecretKey), nil
	})
	if err != nil {
		// 如果过期, 返回自定义错误
		if !token.Valid {
			return nil, errors.New("token expired")
		}
		return nil, err
	}
	return claims, nil
}

// // RefreshAccessToken 刷新access_token
// func (s *JwtService) RefreshAccessToken(ctx *gin.Context) (string, string, *JwtCustomClaims, error) {
// 	refresh_token, err := ctx.Cookie("refresh_token")
// 	if err != nil {
// 		return "", "", nil, errors.New("refresh_token invalid")
// 	}
// 	claims, err := s.ValidateToken(refresh_token)
// 	if err != nil {
// 		return "", "", nil, err
// 	}

// 	access_token, err := s.GenerateAccessToken(claims.UserId, claims.SessionId)
// 	if err != nil {
// 		return "", "", nil, err
// 	}

// 	// 如果refresh_token的剩余时间小于1小时, 续签refresh_token
// 	ok := time.Until(claims.ExpiresAt.Time) < 1*time.Hour
// 	if ok {
// 		newRefreshToken, err := s.GenerateRefreshToken(claims.UserId, claims.SessionId)
// 		if err != nil {
// 			return "", "", nil, err
// 		}
// 		refresh_token = newRefreshToken
// 	}

// 	return access_token, refresh_token, claims, nil
// }

func (s *JwtService) GetTokenExpiration() time.Duration {
	return s.jwtConf.AccessTokenExpiresIn.AsDuration()
}

func (s *JwtService) GetRefreshTokenExpiration() time.Duration {
	return s.jwtConf.RefreshTokenExpiresIn.AsDuration()
}
