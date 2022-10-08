package bot_template

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

var Client *disgord.Client

type Bot struct {
	Client                *disgord.Client
	Commands              []*disgord.CreateApplicationCommand
	ActiveCommandHandlers map[string]interactionHandler
	ComponentHandlers     map[string]interactionHandler
}

func NewBot(config disgord.Config) Bot {
	return Bot{
		Client:                disgord.New(config),
		Commands:              []*disgord.CreateApplicationCommand{},
		ActiveCommandHandlers: map[string]interactionHandler{},
		ComponentHandlers:     map[string]interactionHandler{},
	}
}
func (b *Bot) Run() error {
	defer Client.Gateway().StayConnectedUntilInterrupted()
	Client.Gateway().BotReady(func() {

		user, _ := Client.CurrentUser().Get()
		for _, command := range b.Commands {
			err := Client.ApplicationCommand(0).Global().Create(command)
			if err != nil {
				Client.Logger().Error(err)
			}
		}
		Client.Logger().Info(fmt.Sprintf("Logged in as %s#%s ", user.Username, user.Discriminator))
	})
	Client.Gateway().InteractionCreate(func(s disgord.Session, h *disgord.InteractionCreate) {
		switch h.Type {
		case disgord.InteractionApplicationCommand:
			{
				b.ActiveCommandHandlers[h.Data.Name](s, h)
			}
		case disgord.InteractionMessageComponent:
			{
				b.ComponentHandlers[h.Data.CustomID](s, h)
			}
		}
	})
	return nil
}
