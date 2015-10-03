package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ivpusic/golog"
)

var (
	logger    = golog.GetLogger("github.com/ivpusic/rerun")
	TEST_MODE = false
)

func main() {
	conf, err := loadConfiguration()
	if err != nil {
		logger.Panicf("Error while loading configuration! %s", err.Error())
	}

	// setup logger level
	if *verbose {
		logger.Level = golog.DEBUG
	} else {
		logger.Level = golog.INFO
	}

	pm := &processManager{
		conf: conf,
	}

	w := &watcher{pm: pm}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// on ctrl+c remove build files
	go func(conf *config) {
		<-sigs

		err := os.Remove(conf.build)
		if err != nil && !os.IsNotExist(err) {
			logger.Warnf("Build file not removed! %s", err.Error())
		}

		os.Exit(0)
	}(conf)

	w.start()
}
