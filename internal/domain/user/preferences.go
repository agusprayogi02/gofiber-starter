package user

import (
	"encoding/json"

	"gorm.io/gorm"
)

// UserPreferences stores user-specific settings and preferences
type UserPreferences struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	UserID    uint           `gorm:"uniqueIndex;not null"`
	User      User           `gorm:"foreignKey:UserID"`
	Data      string         `gorm:"type:jsonb"` // JSON data for flexible preferences
	DeletedAt gorm.DeletedAt `gorm:"index"`
	gorm.Model
}

// PreferencesData represents the structure of preferences JSON data
type PreferencesData struct {
	// Notification preferences
	EmailNotifications bool `json:"email_notifications"`
	PushNotifications  bool `json:"push_notifications"`
	SMSNotifications   bool `json:"sms_notifications"`

	// Privacy preferences
	ProfileVisibility string `json:"profile_visibility"` // public, private, friends
	ShowEmail         bool   `json:"show_email"`
	ShowOnlineStatus  bool   `json:"show_online_status"`

	// UI preferences
	Theme    string `json:"theme"`    // light, dark, auto
	Language string `json:"language"` // en, id, etc
	Timezone string `json:"timezone"`

	// Other preferences (can be extended)
	Custom map[string]interface{} `json:"custom,omitempty"`
}

// GetData returns parsed preferences data
func (p *UserPreferences) GetData() (*PreferencesData, error) {
	if p.Data == "" {
		return &PreferencesData{
			EmailNotifications: true,
			PushNotifications:  true,
			SMSNotifications:   false,
			ProfileVisibility:  "public",
			ShowEmail:          false,
			ShowOnlineStatus:   true,
			Theme:              "auto",
			Language:           "en",
			Timezone:           "UTC",
			Custom:             make(map[string]interface{}),
		}, nil
	}

	var data PreferencesData
	if err := json.Unmarshal([]byte(p.Data), &data); err != nil {
		return nil, err
	}

	if data.Custom == nil {
		data.Custom = make(map[string]interface{})
	}

	return &data, nil
}

// SetData sets preferences data from PreferencesData struct
func (p *UserPreferences) SetData(data *PreferencesData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	p.Data = string(jsonData)
	return nil
}

// UpdateField updates a specific field in preferences
func (p *UserPreferences) UpdateField(key string, value interface{}) error {
	data, err := p.GetData()
	if err != nil {
		return err
	}

	// Update specific fields based on key
	switch key {
	case "email_notifications":
		if v, ok := value.(bool); ok {
			data.EmailNotifications = v
		}
	case "push_notifications":
		if v, ok := value.(bool); ok {
			data.PushNotifications = v
		}
	case "sms_notifications":
		if v, ok := value.(bool); ok {
			data.SMSNotifications = v
		}
	case "profile_visibility":
		if v, ok := value.(string); ok {
			data.ProfileVisibility = v
		}
	case "show_email":
		if v, ok := value.(bool); ok {
			data.ShowEmail = v
		}
	case "show_online_status":
		if v, ok := value.(bool); ok {
			data.ShowOnlineStatus = v
		}
	case "theme":
		if v, ok := value.(string); ok {
			data.Theme = v
		}
	case "language":
		if v, ok := value.(string); ok {
			data.Language = v
		}
	case "timezone":
		if v, ok := value.(string); ok {
			data.Timezone = v
		}
	default:
		// Store in custom map
		if data.Custom == nil {
			data.Custom = make(map[string]interface{})
		}
		data.Custom[key] = value
	}

	return p.SetData(data)
}
