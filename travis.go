package githubbot

import (
	"fmt"
	"strconv"

	"github.com/cpalone/travishook"
)

func (s *Session) travisServer(port int) {
	server := travishook.NewServer(port, "/travishook")
	server.GoListenAndServe()
	for {
		p := <-server.Out
		var emoji string
		if p.StatusMessage == "Passed" || p.StatusMessage == "Fixed" {
			emoji = ":white_check_mark:"
		} else {
			emoji = ":no_entry:"
		}

		fmt.Printf("Received payload with status: %s\n", p.StatusMessage)
		s.sendMessage(fmt.Sprintf(
			"%s [ travis.ci | Branch: %s | %s ] %s | %s.",
			emoji, p.Repository.Name, p.Branch, p.Message, p.StatusMessage),
			s.commitMsgID,
			strconv.Itoa(s.msgID))
		s.msgID++
	}
}
