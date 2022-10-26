package bot_template

import (
	"errors"
	"fmt"

	"github.com/andersfylling/disgord"
)

type BotCommand struct {
	Handler InteractionHandler
	Options BotCommandOptions
}
type BotComponent struct {
	Handler InteractionHandler
	Options BotComponentOptions
}
type BotComponentOptions struct {
	// Applies to any error message when clicking a button
	GeneralErrorMessage disgord.CreateInteractionResponse
	// Only one user can click the button
	UserLockedTo disgord.Snowflake
	// Sent when somebody that isn't userlockedto clicks the button
	OutsideInteractionErrorMessage disgord.CreateInteractionResponse
	// Choose if the button will only respond to one click (won't disable)
	OneClickOnly bool
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

func NewBotWithConfig(clientConfig BotConfig) Bot {
	newBot := NewBotWithDefault()
	newBot.Config = clientConfig
	return newBot
}
func NewBotWithDefault() Bot {
	return Bot{
		Commands:                []*disgord.CreateApplicationCommand{},
		CommandHandlers:         map[string]BotCommand{},
		ActiveComponentHandlers: map[string]BotComponent{},
		Config: BotConfig{
			DefaultPermissionErrorMessage: disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Flags: disgord.MessageFlagEphemeral,
					Embeds: []*disgord.Embed{
						{
							Title: "You do not have the required permission for this command.",
							Color: 0xff0000,
						},
					},
				},
			},
			DefaultOutsideInteractionErrorMessage: disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Flags: disgord.MessageFlagEphemeral,
					Embeds: []*disgord.Embed{
						{
							Title: "This isn't your button!",
							Color: 0xff0000,
						},
					},
				},
			},
			DefaultGeneralErrorMessage: disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Flags: disgord.MessageFlagEphemeral,
					Embeds: []*disgord.Embed{
						{
							Title: "An error occured performing this action.",
							Color: 0xff0000,
						},
					},
				},
			},
		},
	}
}

func (b *Bot) Run() error {
	if b.Client == nil {
		return errors.New("cannot run bot without a client")
	}
	client := b.Client
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
	client.Gateway().InteractionCreate(b.handleInteraction)

	err := client.Gateway().StayConnectedUntilInterrupted()
	if err != nil {
		return err
	}
	return nil
}
