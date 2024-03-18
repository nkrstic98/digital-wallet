package config

var DefaultConfig = Config{
	Web: WebConfig{
		Host: "localhost",
		Port: "8080",
	},
	Database: DatabaseConfig{
		Host: "localhost",
		Port: "5432",
		DB:   "users",
		User: "dw_role",
		Pwd:  "5tBsPvvXUBDw25zt",
	},
}
