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
		case gohook.CommitCommentEventType:
			payload, ok := et.Event.(*gohook.CommitCommentEvent)
			if !ok {
				panic("Malformed *CommitCommentEvent.")
			}
			msg := fmt.Sprintf("[ %s ] Comment on commit: %s (%s)",
				payload.Repository.Name,
				payload.Comment.Body,
				payload.Comment.HTMLURL,
			)
			s.sendMessage(msg, "")
		case gohook.IssueCommentEventType:
			payload, ok := et.Event.(*gohook.IssueCommentEvent)
			if !ok {
				panic("Malformed *CommitCommentEvent.")
			}
			msg := fmt.Sprintf("[ %s ] Comment on issue '%s': %s (%s)",
				payload.Repository.Name,
				payload.Issue.Title,
				payload.Comment.Body,
				payload.Comment.HTMLURL,
			)
			s.sendMessage(msg, "")
		case gohook.IssuesEventType:
			payload, ok := et.Event.(*gohook.IssuesEvent)
			if !ok {
				panic("Malformed *IssuesEvent.")
			}
			msg := fmt.Sprintf("[ %s ] Issue '%s' was %s. (%s)",
				payload.Repository.Name,
				payload.Issue.Title,
				payload.Action,
				payload.Issue.HTMLURL,
			)
			s.sendMessage(msg, "")
		}
	}
}
