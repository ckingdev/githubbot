package main

import (
	"flag"
	"os"

	"github.com/Sirupsen/logrus"
	gb "github.com/cpalone/githubbot"
)

var roomName string
var password string
var verbose bool
var port int
var secret string

func init() {
	flag.StringVar(&roomName, "room", "test", "room for the bot to join.")
	flag.StringVar(&password, "pass", "", "optional password for the bot to join.")
	flag.BoolVar(&verbose, "v", true, "Toggle whether debug statements are displayed.")
	flag.IntVar(&port, "port", 8081, "Specify the port to listen on for webhook events.")
	flag.StringVar(&secret, "secret", "", "Secret string used to encrypt webhook events.")
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
	logger.Infoln("Flags parsed. Creating session...")
	s, err := gb.NewSession(roomName, password, port, secret, logger)
	if err != nil {
		logger.Fatalf("Fatal error: creating session: %s", err)
	}
	logger.Infoln("Session created.")
	s.Run()
}
