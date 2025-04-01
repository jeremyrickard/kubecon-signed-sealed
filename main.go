package main

import (
	_ "embed"
	"os"

	"github.com/jeremyrickard/kubecon-2025-ssd/cmd"
	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
