package config

type Config struct {
	MySQL MySQLConfig
}

type MySQLConfig struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}
