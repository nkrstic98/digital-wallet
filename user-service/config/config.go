package config

type Config struct {
	Web      WebConfig
	Database DatabaseConfig
}

type WebConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Host string
	Port string
	DB   string
	User string
	Pwd  string
}
