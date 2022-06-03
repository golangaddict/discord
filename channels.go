package discord

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type ChannelListener struct {
	ErrorHandler func(*discordgo.Session, error)

	channels map[string]func(*discordgo.Session, *discordgo.Message) error
}

func NewChannelListener(sess *discordgo.Session) *ChannelListener {
	cl := &ChannelListener{
		ErrorHandler: defaultErrorHandler,
		channels:     make(map[string]func(*discordgo.Session, *discordgo.Message) error),
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		ch, err := s.Channel(m.ChannelID)
		if err != nil {
			cl.ErrorHandler(s, err)
			return
		}

		callback, ok := cl.channels[ch.ID]
		if !ok {
			callback, ok = cl.channels[ch.Name]
		}

		if ok {
			if err := callback(s, m.Message); err != nil {
				log.Println(err)
			}
		}
	})

	return cl
}

func (l *ChannelListener) AddChannel(identifier string, f func(*discordgo.Session, *discordgo.Message) error) {
	l.channels[identifier] = f
}
