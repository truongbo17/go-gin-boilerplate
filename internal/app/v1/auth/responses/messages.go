package responses

import "github.com/truongbo17/go-gin-boilerplate/internal/i18n"

var messages = map[string]map[string]string{
	"en": {
		"auth.login_failed":             "An error occurred during login, please try again.",
		"auth.user_not_found":           "User not found.",
		"auth.wrong_password":           "Incorrect password.",
		"auth.generate_token_failed":    "Failed to generate authentication token.",
		"auth.user_exists":              "User already exists.",
		"auth.register_failed":          "An error occurred during register, please try again.",
		"auth.reset_requested":          "If the account exists, a reset email will be sent.",
		"auth.reset_invalid":            "Invalid or expired reset token.",
		"auth.reset_unavailable":        "Password reset is temporarily unavailable.",
		"auth.reset_failed":             "Unable to reset password, please try again.",
		"auth.logout_failed":            "An error occurred during logout, please try again.",
		"user_list.internal_error":      "An error occurred, please try again.",
		"err.user_change_pass":          "Password incorrect, please try again later.",
		"err.refresh_token":             "Unable to get new token, please login again.",
		"error.find_user_failed":        "Find user Failed.",
		"error.user_not_found":          "User not found.",
		"error.update_user_failed":      "Update User Failed.",
		"role.internal_error":           "An error occurred, please try again.",
		"role.not_found":                "Role not found.",
		"role.create_failed":            "Failed to create role.",
		"role.update_failed":            "Failed to update role.",
		"role.delete_failed":            "Failed to delete role.",
		"role.assign_permission_failed": "Failed to assign permission to role.",
		"role.assign_user_failed":       "Failed to assign role to user.",
		"permission.internal_error":     "An error occurred, please try again.",
		"permission.not_found":          "Permission not found.",
		"permission.create_failed":      "Failed to create permission.",
		"permission.update_failed":      "Failed to update permission.",
		"permission.delete_failed":      "Failed to delete permission.",
	},
	"vi": {
		"auth.login_failed":             "Có lỗi xảy ra khi đăng nhập, vui lòng thử lại.",
		"auth.register_failed":          "Có lỗi xảy ra khi đăng ký, vui lòng thử lại.",
		"auth.reset_requested":          "Nếu tài khoản tồn tại, email đặt lại mật khẩu sẽ được gửi.",
		"auth.reset_invalid":            "Mã đặt lại mật khẩu không hợp lệ hoặc đã hết hạn.",
		"auth.reset_unavailable":        "Tạm thời không thể yêu cầu đặt lại mật khẩu.",
		"auth.reset_failed":             "Không thể đặt lại mật khẩu, vui lòng thử lại.",
		"auth.logout_failed":            "Có lỗi xảy ra khi đăng xuất, vui lòng thử lại.",
		"auth.user_not_found":           "Người dùng không tồn tại.",
		"auth.wrong_password":           "Mật khẩu không chính xác.",
		"auth.generate_token_failed":    "Lỗi khi tạo token xác thực.",
		"auth.user_exists":              "Người dùng đã tồn tại.",
		"user_list.internal_error":      "Có lỗi xảy ra, vui lòng thử lại.",
		"err.refresh_token":             "Không thể lấy token mới, vui lòng đăng nhập lại.",
		"err.user_change_pass":          "Mật khẩu cũ không khớp, vui lòng thử lại sau.",
		"error.find_user_failed":        "Không tìm thấy user.",
		"error.user_not_found":          "user không tồn tại.",
		"error.update_user_failed":      "Có lỗi xảy ra khi cập nhật user.",
		"role.internal_error":           "Có lỗi xảy ra, vui lòng thử lại.",
		"role.not_found":                "Role không tồn tại.",
		"role.create_failed":            "Lỗi khi tạo role.",
		"role.update_failed":            "Lỗi khi cập nhật role.",
		"role.delete_failed":            "Lỗi khi xóa role.",
		"role.assign_permission_failed": "Lỗi khi gán quyền cho role.",
		"role.assign_user_failed":       "Lỗi khi gán role cho user.",
		"permission.internal_error":     "Có lỗi xảy ra, vui lòng thử lại.",
		"permission.not_found":          "Permission không tồn tại.",
		"permission.create_failed":      "Lỗi khi tạo permission.",
		"permission.update_failed":      "Lỗi khi cập nhật permission.",
		"permission.delete_failed":      "Lỗi khi xóa permission.",
	},
}

func authMessage(language, key string) string {
	if message, ok := messages[language][key]; ok {
		return message
	}
	return i18n.GetMessage(language, key, nil)
}
