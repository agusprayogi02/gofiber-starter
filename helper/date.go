package helper

import (
	"time"

	"starter-gofiber/variables"
)

func TimeNow() string {
	return time.Now().Format(variables.FORMAT_TIME)
}
