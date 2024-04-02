package web

import (
	"launchpad/config"
	"launchpad/internal/server"
)

func Run() {
	server.Serv(Router(), config.Configure.Web.Address)
}
