package model

// Config is a struct to hold the configuration data for the program
type Config struct {
	APIKey   string `json:"api-key"`
	SavePath string `json:"save-path"`
}
