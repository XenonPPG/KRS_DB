package models

type Config struct {
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresDB       string `env:"POSTGRES_DB"`

	GrpcNetwork string `env:"GRPC_NETWORK"`
	GrpcPort    string `env:"GRPC_PORT"`
}
