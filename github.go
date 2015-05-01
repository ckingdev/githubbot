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
		switch et.Type {
		case gohook.PingEventType:
			continue
		case gohook.PushEventType:
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
