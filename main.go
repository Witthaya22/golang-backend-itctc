package main

import (
	"github.com/Witthaya22/golang-backend-itctc/config"
	"github.com/Witthaya22/golang-backend-itctc/databases"
	"github.com/Witthaya22/golang-backend-itctc/modules/servers"
)

func main() {
	conf := config.ConfigGeting()
	db := databases.NewPostgresDb(conf.Database)

	servers.NewServer(conf, db.ConnectionGetting()).Start()
}
