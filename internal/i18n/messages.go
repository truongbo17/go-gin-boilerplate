package i18n

import (
	"strings"

	"github.com/truongbo17/go-gin-boilerplate/config"
)

func Supported(language string) bool {
	for _, available := range config.LanguageAvailable {
		if language == available {
			return true
		}
	}
	return false
}

var messages = map[string]map[string]string{
	"en": {
		"error.app":   "An error occurred, please try again later.",
		"success.200": "Success.",
		"success.201": "Created successfully.",
		"success.202": "Accepted, processing asynchronously.",
		"success.204": "No content.",
		"success.206": "Partial content.",
		"error.400":   "Bad request.",
		"error.401":   "Unauthorized.",
		"error.402":   "Payment required.",
		"error.403":   "Forbidden.",
		"error.404":   "Not found.",
		"error.405":   "Method not allowed.",
		"error.406":   "Not acceptable.",
		"error.407":   "Proxy authentication required.",
		"error.408":   "Request timeout.",
		"error.409":   "Conflict.",
		"error.410":   "Gone.",
		"error.411":   "Length required.",
		"error.412":   "Precondition failed.",
		"error.413":   "Payload too large.",
		"error.414":   "URI too long.",
		"error.415":   "Unsupported media type.",
		"error.416":   "Range not satisfiable.",
		"error.417":   "Expectation failed.",
		"error.422":   "Unprocessable entity.",
		"error.429":   "Too many requests.",
		"error.500":   "Internal server error.",
		"error.501":   "Not implemented.",
		"error.502":   "Bad gateway.",
		"error.503":   "Service unavailable.",
		"error.504":   "Gateway timeout.",
		"error.505":   "HTTP version not supported.",

		"validation.required":               "This field is required.",
		"validation.length":                 "The length must be between {min} and {max} characters.",
		"validation.min_length":             "The minimum length must be {min} characters.",
		"validation.max_length":             "The maximum length must be {max} characters.",
		"validation.email":                  "Invalid email format.",
		"validation.url":                    "Invalid URL format.",
		"validation.in":                     "Invalid value. Must be one of {values}.",
		"validation.numeric":                "Must be a numeric value.",
		"validation.int":                    "Must be an integer.",
		"validation.float":                  "Must be a floating-point number.",
		"validation.min":                    "Value must be at least {min}.",
		"validation.max":                    "Value must be at most {max}.",
		"validation.range":                  "Value must be between {min} and {max}.",
		"validation.regex":                  "Invalid format.",
		"validation.phone":                  "Invalid phone number.",
		"validation.uuid":                   "Invalid UUID format.",
		"validation.credit_card":            "Invalid credit card number.",
		"validation.isbn":                   "Invalid ISBN format.",
		"validation.unique":                 "This value must be unique.",
		"validation.file":                   "Invalid file format.",
		"validation.invalid_file_extension": "File {file} is invalid. Only {ext} are allowed.",
		"validation.max_file_size":          "File {file} too large, maximum size is {max_size}.",
		"validation.dir":                    "Invalid directory path.",
		"validation.eq_field":               "Must be equal to {field}.",
		"validation.ne_field":               "Must not be equal to {field}.",
		"validation.confirmation":           "Must match the confirmation field.",
		"validation.password_uppercase":     "Must contain at least one uppercase letter.",
		"validation.password_lowercase":     "Must contain at least one lowercase letter.",
		"validation.password_digit":         "Must contain at least one digit.",
		"validation.password_special":       "Must contain at least one special character (!@#$%^&*(),.?\":{}|<>).",
		"validation.password_mismatch":      "Confirm password must match the new password.",
		"validation.datetime.format":        "Invalid date time format.",
		"validation.datetime.before":        "Date cannot be in future.",
		"validation.datetime.after":         "Date cannot be in pass.",
	},
	"vi": {
		"error.app":   "Có lỗi xảy ra, vui lòng thử lại sau.",
		"success.200": "Thành công.",
		"success.201": "Tạo thành công.",
		"success.202": "Đã chấp nhận, đang xử lý bất đồng bộ.",
		"success.204": "Không có nội dung.",
		"success.206": "Nội dung một phần.",
		"error.400":   "Yêu cầu không hợp lệ.",
		"error.401":   "Chưa xác thực.",
		"error.402":   "Yêu cầu thanh toán.",
		"error.403":   "Không có quyền truy cập.",
		"error.404":   "Không tìm thấy.",
		"error.405":   "Phương thức không được phép.",
		"error.406":   "Không thể chấp nhận.",
		"error.407":   "Yêu cầu xác thực proxy.",
		"error.408":   "Hết thời gian yêu cầu.",
		"error.409":   "Xung đột dữ liệu.",
		"error.410":   "Tài nguyên không còn tồn tại.",
		"error.411":   "Yêu cầu độ dài nội dung.",
		"error.412":   "Điều kiện tiên quyết không thỏa mãn.",
		"error.413":   "Dữ liệu gửi lên quá lớn.",
		"error.414":   "Đường dẫn quá dài.",
		"error.415":   "Loại dữ liệu không được hỗ trợ.",
		"error.416":   "Phạm vi yêu cầu không hợp lệ.",
		"error.417":   "Yêu cầu không đáp ứng được kỳ vọng.",
		"error.422":   "Dữ liệu không hợp lệ.",
		"error.429":   "Quá nhiều yêu cầu gửi đi.",
		"error.500":   "Lỗi máy chủ nội bộ.",
		"error.501":   "Chưa được hỗ trợ.",
		"error.502":   "Cổng kết nối lỗi.",
		"error.503":   "Dịch vụ tạm thời không hoạt động.",
		"error.504":   "Hết thời gian chờ cổng.",
		"error.505":   "Phiên bản HTTP không được hỗ trợ.",

		"validation.required":               "Là trường bắt buộc.",
		"validation.length":                 "Độ dài phải từ {min} đến {max} ký tự.",
		"validation.min_length":             "Độ dài tối thiểu là {min} ký tự.",
		"validation.max_length":             "Độ dài tối đa là {max} ký tự.",
		"validation.email":                  "Định dạng email không hợp lệ.",
		"validation.url":                    "Định dạng URL không hợp lệ.",
		"validation.in":                     "Giá trị không hợp lệ. Phải là một trong các giá trị sau: {values}.",
		"validation.numeric":                "Phải là số.",
		"validation.int":                    "Phải là số nguyên.",
		"validation.float":                  "Phải là số thực.",
		"validation.min":                    "Giá trị tối thiểu là {min}.",
		"validation.max":                    "Giá trị tối đa là {max}.",
		"validation.range":                  "Giá trị phải nằm trong khoảng từ {min} đến {max}.",
		"validation.regex":                  "Định dạng không hợp lệ.",
		"validation.phone":                  "Số điện thoại không hợp lệ.",
		"validation.uuid":                   "Định dạng UUID không hợp lệ.",
		"validation.credit_card":            "Số thẻ tín dụng không hợp lệ.",
		"validation.isbn":                   "Định dạng ISBN không hợp lệ.",
		"validation.unique":                 "Giá trị này phải là duy nhất.",
		"validation.file":                   "Định dạng tệp không hợp lệ.",
		"validation.invalid_file_extension": "Tập tin {file} không hợp lệ. Chỉ chấp nhận {ext}.",
		"validation.max_file_size":          "Tập tin {file} quá lớn, tối đa kích thước là {max_size}.",
		"validation.dir":                    "Đường dẫn thư mục không hợp lệ.",
		"validation.eq_field":               "Phải bằng với {field}.",
		"validation.ne_field":               "Phải khác với {field}.",
		"validation.confirmation":           "Phải khớp với trường xác nhận.",
		"validation.password_uppercase":     "Phải chứa ít nhất một chữ cái viết hoa.",
		"validation.password_lowercase":     "Phải chứa ít nhất một chữ cái viết thường.",
		"validation.password_digit":         "Phải chứa ít nhất một chữ số.",
		"validation.password_special":       "Phải chứa ít nhất một ký tự đặc biệt (!@#$%^&*(),.?\":{}|<>).",
		"validation.password_mismatch":      "Mật khẩu xác nhận phải trùng với mật khẩu mới.",
		"validation.datetime.format":        "Định dạng ngày không hợp lệ.",
		"validation.datetime.before":        "Ngày không thể ở tương lai.",
		"validation.datetime.after":         "Ngày không thể ở quá khứ.",
	},
}

func GetMessage(lang, key string, params map[string]string) string {
	msg, exists := messages[lang][key]
	if !exists {
		return key
	}

	for k, v := range params {
		msg = strings.ReplaceAll(msg, "{"+k+"}", v)
	}
	return msg

}
