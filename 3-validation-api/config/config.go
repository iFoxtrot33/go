package config

type Config struct {
	SMTP struct {
		From     string `json:"from"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
	} `json:"smtp"`
	ServerAddress string `json:"server_address"`
}

func NewConfig() (*Config, error) {
	return &Config{
		SMTP: struct {
			From     string `json:"from"`
			Password string `json:"password"`
			Host     string `json:"host"`
			Port     string `json:"port"`
		}{
			From:     "your-email@example.com",
			Password: "your-password",
			Host:     "smtp.gmail.com",
			Port:     "587",
		},
		ServerAddress: "localhost:8081",
	}, nil
}
