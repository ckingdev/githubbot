package githubbot

import (
	"fmt"

	"github.com/cpalone/dronehook"
)

func (s *Session) droneServer(port int) {
	server := dronehook.NewServer(port, "/dronehook")
	server.GoListenAndServe()
	for {
		p := <-server.Out
		var emoji string
		if p.Commit.Status == "Success" {
			emoji = ":white_check_mark:"
		} else {
			emoji = ":no_entry:"
		}
		str := fmt.Sprintf("%s [ drone.io | Branch: %s | %s ] %s | %s",
			emoji,
			p.Repository.Name,
			p.Commit.Branch,
			p.Commit.Message,
			p.Commit.Status,
		)
		s.sendMessage(str, "")
	}
}
