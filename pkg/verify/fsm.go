package verify

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type stateMap map[string]func(m *StateParams) error

func (s stateMap) get(m *StateParams) (func(m *StateParams) error, bool) {
	if m.User == nil || m.Guild == nil {
		return nil, false
	}
	out, ok := s[m.User.ID+m.Guild.ID]
	return out, ok
}

func (s stateMap) set(m *StateParams, p func(m *StateParams) error) {
	s[m.User.ID+m.Guild.ID] = p
}

var (
	state   = stateMap{}
	session *discordgo.Session
)

func Init(s *discordgo.Session) {
	session = s
}

type StateParams struct {
	User  *discordgo.User
	Guild *discordgo.Guild
}

// Transition from one state to another base on the decision tree.
func Transition(s *StateParams) error {
	if handler, ok := state.get(s); ok {
		handler(s)
	} else if err := createLink(s); err != nil {
		return err
	}
	return nil
}

func createLink(m *StateParams) error {
	session.ChannelMessageSend(m.User.ID, fmt.Sprintf("To register for %s, click on the following link: %s", m.Guild.Name, "link"))
	state.set(m, handleSuiteAuth)
	return nil
}

func handleSuiteAuth(m *StateParams) error {
	return nil
}
