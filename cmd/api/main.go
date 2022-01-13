package main

import (
	"app/config"
	"app/internal/server"
	"app/pkg/db/mysql"
	"log"
	"os"
)

func main() {
	cfg := &config.Config{}

	cfg.Server.Host = os.Getenv("MY_SERVER_HOST")

	cfg.MySQL.Host = os.Getenv("MYSQL_HOST")
	cfg.MySQL.Port = os.Getenv("MYSQL_PORT")
	cfg.MySQL.User = os.Getenv("MYSQL_USER")
	cfg.MySQL.Pass = os.Getenv("MYSQL_PASS")
	cfg.MySQL.DB = os.Getenv("MYSQL_DB")

	log.Println(cfg)
	db, err := mysql.NewMysqlDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("server is starting...")
	s := server.NewServer(db, cfg)
	log.Fatal(s.Run())
}
