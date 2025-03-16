package auth

import (
	"fmt"
	"time"

	"tea-logistics/pkg/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

/*
 * 認証ユーティリティパッケージ
 * パスワードのハッシュ化やJWTトークンの生成を行う
 */

const (
	// トークンの有効期限（24時間）
	tokenExpiration = 24 * time.Hour
)

// HashPassword パスワードをハッシュ化する
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("パスワードのハッシュ化エラー: %v", err)
	}
	return string(bytes), nil
}

// CheckPassword パスワードが正しいか確認する
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken JWTトークンを生成する
func GenerateToken(user *models.User, secret string) (string, int64, error) {
	expiresAt := time.Now().Add(tokenExpiration)

	claims := models.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, fmt.Errorf("トークン生成エラー: %v", err)
	}

	return tokenString, expiresAt.Unix(), nil
}

// ValidateToken トークンを検証する
func ValidateToken(tokenString, secret string) (*models.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("予期しない署名方式: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("トークン検証エラー: %v", err)
	}

	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("無効なトークン")
}
