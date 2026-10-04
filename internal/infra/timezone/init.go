package timezone

import "os"

func init() {
	_ = os.Setenv("TZ", "Asia/Ho_Chi_Minh")
}
