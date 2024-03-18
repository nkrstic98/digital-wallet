package config

var DefaultConfig = Config{
	Web: WebConfig{
		Host: "localhost",
		Port: "8000",
	},
	Database: DatabaseConfig{
		Host: "localhost",
		Port: "5432",
		DB:   "transactions",
		User: "dw_role",
		Pwd:  "5tBsPvvXUBDw25zt",
	},
	Kafka: KafkaConfig{
		Host: "localhost",
		Port: "9092",
	},
	Nats: NatsConfig{
		URL: "nats://127.0.0.1:4222",
	},
}
