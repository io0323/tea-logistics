package models

/*
 * ロールモデル
 * ユーザーの権限を定義する
 */

// Role ユーザーロール
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleOperator Role = "operator"
	RoleViewer   Role = "viewer"
)

// IsValidRole ロールが有効かどうかを確認する
func IsValidRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleManager, RoleOperator, RoleViewer:
		return true
	default:
		return false
	}
}

// HasPermission 指定されたロールの権限を持っているかどうかを確認する
func HasPermission(role Role, required Role) bool {
	switch role {
	case RoleAdmin:
		return true
	case RoleManager:
		return required != RoleAdmin
	case RoleOperator:
		return required == RoleOperator || required == RoleViewer
	case RoleViewer:
		return required == RoleViewer
	default:
		return false
	}
}
