package main

import (
	"flag"
	"os"

	"github.com/Sirupsen/logrus"
	gb "github.com/fireside-chat/githubbot"
)

var roomName string
var password string
var verbose bool

func init() {
	flag.StringVar(&roomName, "room", "test", "room for the bot to join.")
	flag.StringVar(&password, "pass", "", "optional password for the bot to join.")
	flag.BoolVar(&verbose, "v", false, "Toggle whether debug statements are displayed.")
}

func main() {
	flag.Parse()
	logFile, err := os.OpenFile("fireside.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Error creating log file: %s\n", err)
	}
	defer logFile.Close()
	logger := logrus.New()
	logger.Out = logFile
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logger.Debugln("Flags parsed. Creating session...")
	s, err := gb.NewSession(roomName, password, logger)
	if err != nil {
		logger.Fatalf("Fatal error: creating session: %s", err)
	}
	logger.Debugln("Session created.")
	s.Run()
}
