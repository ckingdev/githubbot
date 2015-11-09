package githubbot

import (
	"fmt"
	"strconv"

	"github.com/cpalone/dronehook"
	"github.com/cpalone/travishook"
)

func (s *Session) ciHandler() error {
	// spin off travis server
	tServer := travishook.NewServer(8085, "/travishook")
	s.logger.Info("Starting travis server...")
	tServer.GoListenAndServe()

	//spin off drone server
	dServer := dronehook.NewServer(8082, "/dronehook")
	s.logger.Info("Starting drone server...")
	dServer.GoListenAndServe()

	for {
		select {
		case p := <-tServer.Out:
			var parent string
			parent, ok := s.commitParent[p.Commit]
			if !ok {
				parent = ""
			}
			var emoji string
			if p.StatusMessage == "Passed" || p.StatusMessage == "Fixed" {
				emoji = ":white_check_mark:"
			} else {
				emoji = ":no_entry:"
			}
			s.sendMessage(fmt.Sprintf(
				"%s [ travis-ci.org | %s | Branch: %s ] (%s)",
				emoji, p.Repository.Name, p.Branch, p.BuildURL),
				parent, strconv.Itoa(s.msgID))
			s.msgID++
		case p := <-dServer.Out:
			var parent string
			parent, ok := s.commitParent[p.Build.Commit]
			if !ok {
				parent = ""
			}
			var emoji string
			if p.Build.Status == "success" {
				emoji = ":white_check_mark:"
			} else {
				emoji = ":no_entry:"
			}
			// https://drone.in.euphoria.io/euphoria-io/heim/77
			url := fmt.Sprintf("drone.in.euphoria.io/%s/%s/%s",
				p.Repo.Name,
				p.Build.Branch,
				p.Build.Number,
			)
			str := fmt.Sprintf("%s [ drone.io | %s | Branch: %s ] (%s)",
				emoji,
				p.Repo.Name,
				p.Build.Branch,
				url,
			)
			s.sendMessage(str,
				parent, strconv.Itoa(s.msgID))
			s.msgID++
		}
	}
}
