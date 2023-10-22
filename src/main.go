// main.go

package main

import (
	"fmt"
	"math/rand"
	"spysat/config"
	"spysat/datastore"
	"spysat/logger"
	"spysat/operation"
	"spysat/probe"
	"time"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	config.LoadConfig()
	logger.SetLevel(config.Config.LogLevel)
	logger.SetFormat(config.Config.LogFormat)

	router = gin.New()
	router.Use(gin.LoggerWithFormatter(logger.ConsoleLogFormatter))
	router.Use(gin.Recovery())

	logger.Infof("", "Running with port: %d", config.Config.Port)

	// router.LoadHTMLGlob("templates/*")
	initializeRoutes()

	rand.Seed(time.Now().UnixNano())
	datastore.Init()
	operation.LoadOperation()
	for oName, o := range operation.Operations.Observers {
		datastore.AddObserver(oName, o)
	}

	logger.Info("", "Creating probes")
	for pName, p := range operation.Operations.Probes {
		logger.Debugf("", "Checking probe %s", pName)
		probeArgs := make(map[string]map[string]map[string]map[string]interface{})
		for oName, o := range operation.Operations.Observers {
			logger.Debugf("", "Checking observer %s", oName)
			for sName, s := range o.Streams {
				logger.Debugf("", "Checking stream %s", sName)
				if s.Probe == pName {
					logger.Infof("Associating probe %s with stream %s/%s", pName, oName, sName)
					if _, ok := probeArgs[o.Group]; !ok {
						probeArgs[o.Group] = make(map[string]map[string]map[string]interface{})
					}
					if _, ok := probeArgs[o.Group][oName]; !ok {
						probeArgs[o.Group][oName] = make(map[string]map[string]interface{})
					}
					probeArgs[o.Group][oName][sName] = s.Arguments

				}
			}
		}
		logger.Tracef("", "Probe args: %v", probeArgs)
		go probe.Run(p, probeArgs)
	}
	routerPort := fmt.Sprintf(":%d", config.Config.Port)
	router.Run(routerPort)
}
