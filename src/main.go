package main

import (
	"flag"
	"github.com/Lyt99/iop-statistics/config"
	"github.com/Lyt99/iop-statistics/router"
	"github.com/Lyt99/iop-statistics/tool"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	migrateStats   bool
	migrateFormula bool
	indexFormula   bool
)

func main() {
	flag.BoolVar(&migrateStats, "migrate-stats", false, "convert stats")
	flag.BoolVar(&migrateFormula, "migrate-formula", false, "convert formula")
	flag.BoolVar(&indexFormula, "index-formula", false, "index formula")

	flag.Parse()

	ifPar := false

	if migrateStats {
		tool.MigrateStats()
		ifPar = true
	}

	if migrateFormula {
		tool.MigrateFormula()
		ifPar = true
	}
	if indexFormula {
		tool.IndexFormula()
		ifPar = true
	}

	if ifPar {
		return
	}

	startWeb()
}

func startWeb() {
	// TODO: Log
	// TODO: 是不是可以直接缓存每个页面
	e := gin.Default()

	// Router
	g := e.Group("/stats")
	router.RouteStatistics(g)

	log.Fatal(e.Run(config.GlobalConfig.Addr))
}
