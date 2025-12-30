package variables

import (
	"strings"
	"time"
)

var (
	STATIC_PATH = "/storage"
	POST_PATH   = "/post/"
	FORMAT_TIME = time.RFC3339
	ADMIN_ROLE  = "admin"
	USER_ROLE   = "user"
	USER_ID     = "user_id"
)

func GenerateStatic(paths []string) string {
	return STATIC_PATH + strings.Join(paths, "")
}
