package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"testing"
)

func TestEventListener_Add(t *testing.T) {
	l := NewEventListener(nil)
	l.Add(func(s *discordgo.Session, e *discordgo.GuildCreate) error {
		return errors.New("this is a fake error")
	})
}
