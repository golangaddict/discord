package discord

// TODO: change docs to slice for ordering

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

type CommandFunc func(s *discordgo.Session, m *discordgo.Message, args []string) error

type Command struct {
	Name       string
	Usage      string
	Permission *CommandPermission
	Prefix     string
	Func       CommandFunc
}

func NewCommand(name, usage string, f CommandFunc) *Command {
	return &Command{
		Name:  name,
		Usage: usage,
		Func:  f,
	}
}

func (c *Command) SetPrefix(prefix string) *Command {
	c.Prefix = prefix

	return c
}

func (c *Command) SetFunc(f CommandFunc) *Command {
	c.Func = f

	return c
}

func (c *Command) SetPerms(perms int64) *Command {
	if c.Permission == nil {
		c.Permission = new(CommandPermission)
	}

	c.Permission.Perms = perms

	return c
}

func (c *Command) SetPermRoleID(id string) *Command {
	if c.Permission == nil {
		c.Permission = new(CommandPermission)
	}

	c.Permission.Role.ID = id

	return c
}

func (c *Command) SetPermRoleName(name string) *Command {
	if c.Permission == nil {
		c.Permission = new(CommandPermission)
	}

	c.Permission.Role.Name = name

	return c
}

type CommandPermission struct {
	Perms int64
	Role  Role
}

func (c CommandPermission) Authorized(s *discordgo.Session, member *discordgo.Member) (bool, error) {
	for _, roleID := range member.Roles {
		if roleID == c.Role.ID {
			return true, nil
		}

		role, err := s.State.Role(member.GuildID, roleID)
		if err != nil {
			return false, err
		}

		if role.Name == c.Role.Name {
			return true, nil
		}

		if role.Permissions&c.Perms != 0 {
			return true, nil
		}
	}

	return false, nil
}

type CommandListener struct {
	ErrorHandler func(*discordgo.Session, *discordgo.Message, error)

	commands []*Command
}

func NewCommandListner(s *discordgo.Session, prefix string) *CommandListener {
	cl := &CommandListener{
		ErrorHandler: defaultCommandErrorHandler,
	}

	cl.Add(NewCommand("commands", "list commands", cl.Commands))

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if strings.HasPrefix(m.Content, prefix) {
			r := csv.NewReader(strings.NewReader(m.Content))
			r.Comma = ' '
			data, err := r.Read()
			if err != nil {
				log.Println(err)
				return
			}

			commandName := strings.TrimPrefix(data[0], prefix)
			var args []string
			if len(data) > 1 {
				args = data[1:]
			}

			cmd := cl.Find(commandName)
			if cmd == nil {
				return
			}

			if cmd.Permission != nil {
				// set member guild id because it is not set by this event.
				m.Member.GuildID = m.GuildID
				ok, err := cmd.Permission.Authorized(s, m.Member)
				if err != nil {
					cl.ErrorHandler(s, m.Message, err)
					return
				}

				if !ok {
					cl.ErrorHandler(s, m.Message, errors.New("You do not have sufficient permissions to use this command."))
					return
				}
			}

			if err := cmd.Func(s, m.Message, args); err != nil {
				if _, ok := err.(*relayError); ok {
					if err := Notify(s, m.ChannelID, Error, "", err.Error()); err != nil {
						cl.ErrorHandler(s, m.Message, err)
					}
					return
				}

				infoLog.Printf("command error: sender[%s] | name[%s] | args[%s] | error[%s]\n", m.Author.Username, commandName, args, err)
			}
		}
	})

	return cl
}

func (cl *CommandListener) Add(cmd *Command) {
	cl.commands = append(cl.commands, cmd)
}

func (cl *CommandListener) Commands(s *discordgo.Session, m *discordgo.Message, args []string) error {
	var sb strings.Builder
	for _, cmd := range cl.commands {
		sb.WriteString(fmt.Sprintf("```::%s -> %s\n```", cmd.Name, cmd.Usage))
	}

	return Notify(s, m.ChannelID, Info, "Commands", sb.String())
}

func (cl *CommandListener) Find(name string) *Command {
	for _, c := range cl.commands {
		if c.Name == name {
			return c
		}
	}

	return nil
}
