package main

//go:generate rice embed-go

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/mmcloughlin/professor"

	"github.com/prologic/sm"
)

func main() {
	var (
		version bool
		debug   bool

		dburi string
		bind  string
	)

	flag.BoolVar(&version, "v", false, "display version information")
	flag.BoolVar(&debug, "d", false, "enable debug logging")

	flag.StringVar(&dburi, "dburi", "memory://", "database to use")
	flag.StringVar(&bind, "bind", "0.0.0.0:8000", "[int]:<port> to bind to")

	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if version {
		fmt.Printf("sm v%s", sm.FullVersion())
		os.Exit(0)
	}

	if debug {
		go professor.Launch(":6060")
	}

	opts := &sm.Options{}

	metrics := sm.InitMetrics("sm")

	db, err := sm.InitDB(dburi)
	if err != nil {
		log.Errorf("error initializing database: %s", err)
		os.Exit(1)
	}
	defer db.Close()

	server := sm.NewServer(bind, opts)
	server.AddRoute("GET", "/metrics", metrics.Handler())

	log.Infof("sm %s listening on %s", sm.FullVersion(), bind)
	server.ListenAndServe()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	log.Infof("shuting down...")
	server.Shutdown()
}
