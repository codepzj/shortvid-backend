package service

import (
	"context"
	"errors"
	"time"
	"shortvid-backend/app/user-service/internal/conf"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ClaimsContextKey = "claims"
)

type JwtService struct {
	jwtConf *conf.Jwt
}

func NewJwtService(jwtConf *conf.Jwt) *JwtService {
	return &JwtService{jwtConf: jwtConf}
}

type JwtCustomClaims struct {
	UID       int
	SessionID string
	jwt.RegisteredClaims
}

// GenerateAccessToken 创建access_token
func (s *JwtService) GenerateAccessToken(UID int, sessionID string) (string, error) {
	expiresAt := time.Now().Add(s.GetTokenExpiration())
	claims := JwtCustomClaims{
		UID:       UID,
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
func (s *JwtService) GenerateRefreshToken(UID int, sessionID string) (string, error) {
	expiresAt := time.Now().Add(s.GetRefreshTokenExpiration())
	claims := JwtCustomClaims{
		UID:       UID,
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
	if err != nil || !token.Valid {
		return nil, errors.New("token invalid")
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

// 从ctx中设置claims
func (s *JwtService) SetClaimsFromContext(ctx context.Context, claims *JwtCustomClaims) context.Context {
	ctxWithClaims := context.WithValue(ctx, ClaimsContextKey, claims)
	return ctxWithClaims
}

// 从ctx中获取claims
func (s *JwtService) GetClaimsFromContext(ctx context.Context) (*JwtCustomClaims, error) {
	claims, ok := ctx.Value(ClaimsContextKey).(*JwtCustomClaims)
	if !ok {
		return nil, errors.New("no claims in context")
	}
	return claims, nil
}

// 获取token过期时间
func (s *JwtService) GetTokenExpiration() time.Duration {
	return s.jwtConf.AccessTokenExpiresIn.AsDuration()
}

// 获取refresh_token过期时间
func (s *JwtService) GetRefreshTokenExpiration() time.Duration {
	return s.jwtConf.RefreshTokenExpiresIn.AsDuration()
}
