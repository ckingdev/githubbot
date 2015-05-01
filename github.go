package githubbot

import (
	"fmt"

	"github.com/cpalone/gohook"
)

func (s *Session) hookServer(port int, secret string) {
	server := gohook.NewServer(port, secret, "/postreceive")
	s.logger.Debug("Starting webhook server...")
	server.GoListenAndServe()
	s.logger.Debug("...started.")
	for {
		et := <-server.EventAndTypes
		s.logger.Infof("Received hook event of type '%s'.", et.Type)
		switch et.Type {
		case gohook.PingEventType:
			continue
		case gohook.PushEventType:
			s.logger.Debug("Entering PushEventType case.")
			payload, ok := et.Event.(*gohook.PushEvent)
			if !ok {
				panic("Malformed *PushEvent.")
			}
			msg := fmt.Sprintf("[ %s | %s ] Commit: %s (%s)",
				payload.Repository.Name,
				payload.Ref[11:], // this discards "refs/heads/"
				payload.HeadCommit.Message,
				payload.HeadCommit.URL,
			)
			s.sendMessage(msg, "")
		}
	}
}
