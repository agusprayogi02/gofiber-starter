package helper

import (
	"starter-gofiber/config"
	"time"
)

func TimeNow() string {
	return time.Now().Format(config.FORMAT_TIME)
}
