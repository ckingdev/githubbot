package githubbot

import (
	"fmt"

	"github.com/cpalone/gohook"
)

func (s *Session) hookServer(port int, secret string) {
	server := gohook.NewServer(port, secret, "/postreceive")
	s.logger.Info("Starting webhook server...")
	server.GoListenAndServe()
	s.logger.Info("...started.")
	for {
		et := <-server.EventAndTypes
		s.logger.Infof("Received hook event of type '%s'.", et.Type)
		switch et.Type {
		case gohook.PingEventType:
			continue
		case gohook.CommitCommentEventType:
			payload, ok := et.Event.(*gohook.CommitCommentEvent)
			if !ok {
				panic("Malformed *CommitCommentEvent.")
			}
			// TODO: can we get the branch here?
			msg := fmt.Sprintf("[ %s ] Comment on commit: %s (%s)",
				payload.Repository.Name,
				payload.Comment.Body,
				payload.Comment.HTMLURL,
			)
			s.sendMessage(msg, "")
		case gohook.CreateEventType:
			payload, ok := et.Event.(*gohook.CreateEvent)
			if !ok {
				panic("Malformed *CreateEvent.")
			}
			msg := fmt.Sprintf("[ %s | Branch/Tag: %s] Created.",
				payload.Repository.Name,
				payload.RefType,
			)
			s.sendMessage(msg, "")
		case gohook.DeleteEventType:
			payload, ok := et.Event.(*gohook.DeleteEvent)
			if !ok {
				panic("Malformed *DeleteEvent.")
			}
			msg := fmt.Sprintf("[ %s | Branch/Tag: %s] Deleted.",
				payload.Repository,
				payload.RefType,
			)
			s.sendMessage(msg, "")
		case gohook.IssueCommentEventType:
			payload, ok := et.Event.(*gohook.IssueCommentEvent)
			if !ok {
				panic("Malformed *CommitCommentEvent.")
			}
			msg := fmt.Sprintf("[ %s | Issue: %s ] Comment: %s (%s)",
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
			msg := fmt.Sprintf("[ %s | Issue: %s ] Action: %s. (%s)",
				payload.Repository.Name,
				payload.Issue.Title,
				payload.Action,
				payload.Issue.HTMLURL,
			)
			s.sendMessage(msg, "")
		case gohook.PullRequestEventType:
			payload, ok := et.Event.(*gohook.PullRequestEvent)
			if !ok {
				panic("Malformed *PullRequestEvent.")
			}
			msg := fmt.Sprintf("[ %s | PR: %s ] %s",
				payload.Repository.Name,
				payload.PullRequest.Title,
				payload.Action,
			)
			s.sendMessage(msg, "")
		case gohook.PullRequestReviewCommentEventType:
			payload, ok := et.Event.(*gohook.PullRequestReviewCommentEvent)
			if !ok {
				panic("Malformed *PullRequestReviewCommentEvent.")
			}
			msg := fmt.Sprintf("[ %s | PR: %s ] Comment: %s: %s",
				payload.Repository.Name,
				payload.PullRequest.Title,
				payload.Sender.Login,
				payload.Comment.Body,
			)
			s.sendMessage(msg, "")
		case gohook.PushEventType:
			s.logger.Info("Entering PushEventType case.")
			payload, ok := et.Event.(*gohook.PushEvent)
			if !ok {
				panic("Malformed *PushEvent.")
			}
			msg := fmt.Sprintf(":repeat: [ %s | Branch: %s ] Commit: %s (%s)",
				payload.Repository.Name,
				payload.Ref[11:], // this discards "refs/heads/"
				payload.HeadCommit.Message,
				payload.HeadCommit.URL,
			)
			s.sendMessage(msg, "")
		case gohook.RepositoryEventType:
			payload, ok := et.Event.(*gohook.RepositoryEvent)
			if !ok {
				panic("Malformed *RepositoryEvent.")
			}
			msg := fmt.Sprintf("[ Repository: %s ] Action: created. ",
				payload.Repository.Name,
			)
			s.sendMessage(msg, "")
		}

	}
}
