package main

import (
	"log"

	srvConfig "github.com/CHESSComputing/golib/config"
	docdb "github.com/CHESSComputing/golib/docdb"
	server "github.com/CHESSComputing/golib/server"
	"github.com/gin-gonic/gin"
)

// Verbose defines verbosity level
var Verbose int

// metaDB object
var metaDB docdb.DocDB

// helper function to setup our router
func setupRouter() *gin.Engine {
	routes := []server.Route{
		{Method: "POST", Path: "/record", Handler: PostHandler, Authorized: true, Scope: "write"},
	}
	r := server.Router(routes, nil, "static", srvConfig.Config.ProxyWriter.WebServer)
	return r
}

// Server defines our HTTP server
func Server() {
	var err error

	// init docdb
	metaDB, err = docdb.InitializeDocDB(srvConfig.Config.ProxyWriter.MongoDB.DBUri)
	if err != nil {
		log.Fatal(err)
	}

	// init Verbose
	Verbose = srvConfig.Config.ProxyWriter.WebServer.Verbose

	// setup web router and start the service
	r := setupRouter()
	webServer := srvConfig.Config.ProxyWriter.WebServer
	server.StartServer(r, webServer)
}
