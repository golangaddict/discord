package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strings"
	"time"
)

const (
	ColorSuccess = 3066993
	ColorWarning = 16776960
	ColorError   = 15158332
	ColorInfo    = 3447003
)

const (
	Success = iota
	Warning
	Error
	Info
)

const EmptyChar = '\u200B'

type Role struct {
	ID   string
	Name string
}

type RawIdentifier string

func (id RawIdentifier) String() string {
	return regexp.MustCompile("[<#@&>]").ReplaceAllString(string(id), "")
}

type RawEmoji string

func (e RawEmoji) String() string {
	d := strings.Split(string(e), ":")
	if len(d) < 2 {
		return string(e)
	}

	return fmt.Sprintf(":%s:%s", d[1], d[2][:len(d[2])-1])
}

func (e RawEmoji) Format() string {
	log.Println(e)
	return fmt.Sprintf("<%s>", e)
}

func Notify(s *discordgo.Session, channelID string, messageType int, title string, message string) error {
	embed := NewEmbed(messageType)
	embed.Title = title
	embed.Description = message

	return SendChannelMessageEmbed(s, channelID, embed)
}

func NewEmbed(embedType int) *discordgo.MessageEmbed {
	var icon, title string
	var color int
	switch embedType {
	case Success:
		title = "Success"
		color = ColorSuccess
		icon = "https://www.pinclipart.com/picdir/middle/565-5650948_green-tick-check-mark-icon-simple-style-vector.png"
	case Warning:
		title = "Warning"
		color = ColorWarning
		icon = "https://openclipart.org/image/2000px/29833"
	case Error:
		title = "Error"
		color = ColorError
		icon = "https://cdn1.iconfinder.com/data/icons/color-bold-style/21/08-512.png"
	case Info:
		title = "Information"
		color = ColorInfo
		icon = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS354QGgS7hd6qERv5mHR-0MVM5xuZdEmrK3jZLv185tcglYm7tLf8T13MJc0LIid7g8e8&usqp=CAU"
	}

	return &discordgo.MessageEmbed{
		Color: color,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    title,
			IconURL: icon,
		},
	}
}

func ClearChannel(s *discordgo.Session, channelID string) error {
	for {
		cm, err := s.ChannelMessages(channelID, 100, "", "", "")
		if err != nil {
			return err
		}

		if len(cm) == 0 {
			break
		}

		var mes []string
		for _, m := range cm {
			mes = append(mes, m.ID)
		}

		if err := s.ChannelMessagesBulkDelete(channelID, mes); err != nil {
			return err
		}

		if len(cm) < 100 {
			break
		}

		time.Sleep(time.Second * 3)
	}

	return nil
}

func SendChannelMessage(session *discordgo.Session, channelID, message string) error {
	_, err := session.ChannelMessageSend(channelID, message)
	if err != nil {
		return fmt.Errorf("error sending channel message: %w", err)
	}

	return nil
}

func SendChannelMessageEmbed(session *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) error {
	_, err := session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		return fmt.Errorf("error sending channel message embed: %w", err)
	}

	return nil
}

func SendChannelMessageEmbedSimple(session *discordgo.Session, channelID, title, description string) error {
	_, err := session.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{Title: title, Description: description})
	if err != nil {
		return fmt.Errorf("error sending channel message embed: %w", err)
	}

	return nil
}
