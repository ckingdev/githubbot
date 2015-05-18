package githubbot

import (
	"github.com/cpalone/dronehook"
)

func (s *Session) droneServer(port int) {
	server := dronehook.NewServer(port, "/dronehook")
	server.GoListenAndServe()
	for {
		p := <-server.Out
		s.sendMessage(p.String(), "")
	}
}
