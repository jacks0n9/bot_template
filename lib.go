package bot_template

import (
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

type BotCommand struct {
	Handler interactionHandler
	Options BotCommandOptions
}
type BotComponent struct {
	Handler interactionHandler
	Options BotComponentOptions
}
type BotComponentOptions struct {
	// Applies to any error message when clicking a button
	GeneralErrorMessage disgord.CreateInteractionResponse
	// Only one user can click the button
	UserLockedTo disgord.Snowflake
	// Sent when somebody that isn't userlockedto clicks the button
	OutsideInteractionErrorMessage disgord.CreateInteractionResponse
}
type BotCommandOptions struct {
	// Only people who have this perm or above can use command
	RequiredPermission disgord.PermissionBit
	// Sent when you don't have perms
	PermissionErrorMessage disgord.CreateInteractionResponse
	// When any error occurs when running command
	GeneralErrorMessage disgord.CreateInteractionResponse
}
type BotConfig struct {
	DefaultPermissionErrorMessage         disgord.CreateInteractionResponse
	DefaultGeneralErrorMessage            disgord.CreateInteractionResponse
	DefaultOutsideInteractionErrorMessage disgord.CreateInteractionResponse
}
type Bot struct {
	Client                  *disgord.Client
	Config                  BotConfig
	Commands                []*disgord.CreateApplicationCommand
	CommandHandlers         map[string]BotCommand
	ActiveComponentHandlers map[string]BotComponent
}

func NewBotWithConfig(client_config disgord.Config) Bot {
	new := NewBot()
	new.Client = disgord.New(client_config)
	return new
}
func NewBot() Bot {
	return Bot{
		Commands:                []*disgord.CreateApplicationCommand{},
		CommandHandlers:         map[string]BotCommand{},
		ActiveComponentHandlers: map[string]BotComponent{},
	}
}
func (b *Bot) Run() error {
	client := b.Client
	var err error
	defer func() {
		err = client.Gateway().StayConnectedUntilInterrupted()
	}()
	client.Gateway().BotReady(func() {
		user, _ := client.CurrentUser().Get()
		for _, command := range b.Commands {
			err := client.ApplicationCommand(0).Global().Create(command)
			if err != nil {
				client.Logger().Error(err)
			}
		}
		client.Logger().Info(fmt.Sprintf("Logged in as %s#%s ", user.Username, user.Discriminator))
	})
	client.Gateway().InteractionCreate(func(s disgord.Session, h *disgord.InteractionCreate) {
		empty := disgord.CreateInteractionResponse{}
		errResp := empty
		switch h.Type {
		case disgord.InteractionApplicationCommand:
			{
				cmd := b.CommandHandlers[h.Data.Name]
				perms, _ := h.Member.GetPermissions(context.Background(), s)
				if !perms.Contains(cmd.Options.RequiredPermission) {
					if msg := cmd.Options.PermissionErrorMessage; msg != empty {
						errResp = msg
					} else {
						errResp = b.Config.DefaultPermissionErrorMessage
					}
					break
				}
				err := cmd.Handler(s, h)
				if err != nil {
					b.Client.Logger().Error(fmt.Sprintf("Error occured running command %s: %s", h.Data.Name, err))
				}
				if msg := cmd.Options.GeneralErrorMessage; msg != empty {
					errResp = msg
				} else {
					errResp = b.Config.DefaultGeneralErrorMessage
				}

			}
		case disgord.InteractionMessageComponent:
			{
				comp := b.ActiveComponentHandlers[h.Data.CustomID]
				if comp.Options.UserLockedTo != h.User.ID {
					if msg := comp.Options.OutsideInteractionErrorMessage; msg != empty {
						errResp = msg
					} else {
						errResp = b.Config.DefaultOutsideInteractionErrorMessage
					}
					break
				}
				err := comp.Handler(s, h)
				if err != nil {
					if msg := comp.Options.GeneralErrorMessage; msg != empty {
						errResp = msg
					} else {
						errResp = b.Config.DefaultGeneralErrorMessage
					}
					break
				}
			}
		}
		s.SendInteractionResponse(context.Background(), h, &errResp)
	})
	if err != nil {
		return fmt.Errorf("error while connected to gateway: %s", err)
	}
	return nil
}
