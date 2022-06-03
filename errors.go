package discord

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type relayError struct {
	e string
}

func (e *relayError) Error() string {
	return e.e
}

func NewRelayError(err string) *relayError {
	return &relayError{err}
}

func defaultErrorHandler(s *discordgo.Session, err error) {
	log.Println(err)
}

func defaultCommandErrorHandler(s *discordgo.Session, m *discordgo.Message, err error) {
	if e := Notify(s, m.ChannelID, Error, "", err.Error()); e != nil {
		log.Println(e)
	}
}
