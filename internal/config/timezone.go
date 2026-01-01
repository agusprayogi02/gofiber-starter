package config

import (
	"fmt"
	"time"
)

func LoadTimezone() {
	location, err := time.LoadLocation(ENV.TIMEZONE)
	if err != nil {
		panic(fmt.Sprintf("Error loading location: %v", err))
	}
	time.Local = location
}

// GetTimezone returns configured timezone location
func GetTimezone() *time.Location {
	return time.Local
}
