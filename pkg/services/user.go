package services

import (
	"context"
	"fmt"
	"time"

	"tea-logistics/pkg/config"
	"tea-logistics/pkg/models"
	"tea-logistics/pkg/repository"

	"github.com/golang-jwt/jwt"
)

/*
 * ユーザーサービス
 * ユーザー関連のビジネスロジックを実装する
 */

// UserService ユーザーサービス
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService ユーザーサービスを作成する
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Login ユーザーログイン
func (s *UserService) Login(ctx context.Context, req *models.LoginRequest) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if !user.CheckPassword(req.Password) {
		return "", fmt.Errorf("パスワードが正しくありません")
	}

	// JWTトークンの生成
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("トークン生成エラー: %v", err)
	}

	return tokenString, nil
}

// Register ユーザー登録
func (s *UserService) Register(ctx context.Context, req *models.RegisterRequest) error {
	// メールアドレスの重複チェック
	_, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return fmt.Errorf("このメールアドレスは既に登録されています")
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     req.Role,
		Status:   models.UserStatusActive,
	}

	if err := user.HashPassword(); err != nil {
		return fmt.Errorf("パスワードのハッシュ化に失敗しました: %v", err)
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

// GetProfile プロフィール取得
func (s *UserService) GetProfile(ctx context.Context, userID int64) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// UpdateProfile プロフィール更新
func (s *UserService) UpdateProfile(ctx context.Context, userID int64, req *models.UpdateProfileRequest) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Name = req.Name
	return s.repo.UpdateUser(ctx, user)
}

// ChangePassword パスワード変更
func (s *UserService) ChangePassword(ctx context.Context, userID int64, req *models.ChangePasswordRequest) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.CheckPassword(req.OldPassword) {
		return fmt.Errorf("現在のパスワードが正しくありません")
	}

	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return fmt.Errorf("パスワードのハッシュ化に失敗しました: %v", err)
	}

	return s.repo.UpdatePassword(ctx, userID, user.Password)
}
