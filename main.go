package main

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/gofunct/chronic/config"
	"log"
	"os"
	"os/user"
)

func init() {
	u, _ := user.Current()
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller, "user", u)
	log.SetOutput(kitlog.NewStdlibAdapter(logger))
}

func main() {

	var run = config.New()
	if err := run.Init(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := run.Cron.AddFunc("@hourly", func() {
		if err := run.Write(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("starting cron...")
	defer run.Cron.Stop()
	run.Cron.Run()
	run.Cron.Start()
}
