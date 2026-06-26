package model

type HealthResponse struct {
	Name               string   `json:"name"`
	Environment        string   `json:"environment"`
	Database           string   `json:"database"`
	DefaultLanguage    string   `json:"default_language"`
	SupportedLanguages []string `json:"supported_languages"`
	Timestamp          string   `json:"timestamp"`
}
