package config

type Config struct {
	Web      WebConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
	Nats     NatsConfig
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

type KafkaConfig struct {
	Host string
	Port string
}

type NatsConfig struct {
	URL string
}
