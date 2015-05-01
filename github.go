package githubbot

import (
	"fmt"
	"time"

	"github.com/phayes/hookserve/hookserve"
)

func (s *Session) hookServer(port int, secret string) {
	server := hookserve.NewServer()
	server.Port = port
	server.Secret = secret
	s.logger.Debug("Starting webhook server...")
	server.GoListenAndServe()
	s.logger.Debug("...started.")
	for {
		select {
		case event := <-server.Events:
			URL := fmt.Sprintf("github.com/%s/%s/commits/%s", event.Owner, event.Repo, event.Commit)
			s.logger.Debug("Received webhook event of type: %s", event.Type)
			if event.Type != "push" {
				continue
			}
			s.sendMessage(fmt.Sprintf("[%s | %s] Commit pushed. %s", event.Repo, event.Branch, URL), "")
		case <-time.After(time.Duration(30) * time.Second):
			s.logger.Debug("Timed out without an event.")
		}
	}
}
