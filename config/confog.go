package config

type Config struct {
	MySQL  MySQLConfig
	Server ServerConfig
}

type ServerConfig struct {
	Host string
}

type MySQLConfig struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}
