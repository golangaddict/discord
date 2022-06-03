package discord

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"reflect"
)

type EventListener struct {
	ErrorHandler func(*discordgo.Session, error)

	sess   *discordgo.Session
	events map[any]any
}

func NewEventListener(sess *discordgo.Session) *EventListener {
	l := &EventListener{
		ErrorHandler: defaultErrorHandler,
		sess:         sess,
		events:       make(map[any]any),
	}

	sess.AddHandler(func(s *discordgo.Session, a any) {
		for k, v := range l.events {
			if k == reflect.TypeOf(a) {
				res := reflect.ValueOf(v).Call([]reflect.Value{reflect.ValueOf(s), reflect.ValueOf(a)})
				if len(res) > 0 {
					if err, ok := res[0].Interface().(error); ok && err != nil {
						log.Println(res[0].Interface())
					}
				}
			}
		}
	})

	return l
}

func (l *EventListener) Add(f any) {
	t := reflect.TypeOf(f)
	l.events[t.In(1)] = f
}
