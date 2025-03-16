package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"tea-logistics/pkg/auth"
	"tea-logistics/pkg/models"

	"github.com/gin-gonic/gin"
)

/*
 * 認証ハンドラ
 * ログインやユーザー登録のリクエストを処理する
 */

// AuthHandler 認証ハンドラ構造体
type AuthHandler struct {
	db *sql.DB
}

// NewAuthHandler 認証ハンドラを作成する
func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Login ログイン処理
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	// ユーザーの検索
	var user models.User
	var passwordHash string
	err := h.db.QueryRow(
		"SELECT id, username, email, password_hash, name, role, status FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.Name, &user.Role, &user.Status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザー名またはパスワードが正しくありません"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースエラー"})
		return
	}

	// パスワードの検証
	if !auth.CheckPassword(req.Password, passwordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザー名またはパスワードが正しくありません"})
		return
	}

	// ユーザーのステータスチェック
	if user.Status != models.UserStatusActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "アカウントが無効です"})
		return
	}

	// トークンの生成
	token, expiresIn, err := auth.GenerateToken(&user, os.Getenv("JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成エラー"})
		return
	}

	// アクセストークンの保存
	_, err = h.db.Exec(
		"INSERT INTO access_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		user.ID,
		token,
		time.Unix(expiresIn, 0),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン保存エラー"})
		return
	}

	// レスポンスの返却
	c.JSON(http.StatusOK, models.LoginResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		User:      user,
	})
}

// Register ユーザー登録処理
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	// パスワードのハッシュ化
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワード処理エラー"})
		return
	}

	// ユーザーの作成
	var user models.User
	err = h.db.QueryRow(
		`INSERT INTO users (username, email, password_hash, name, role, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`,
		req.Username,
		req.Email,
		passwordHash,
		req.Name,
		req.Role,
		models.UserStatusActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー登録エラー"})
		return
	}

	// ユーザー情報をレスポンスに設定
	user.Username = req.Username
	user.Email = req.Email
	user.Name = req.Name
	user.Role = req.Role
	user.Status = models.UserStatusActive

	c.JSON(http.StatusCreated, gin.H{"message": "ユーザーを登録しました"})
}
