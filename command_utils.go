package bot_template

import (
	"fmt"
	"math/rand"

	"github.com/andersfylling/disgord"
)

type interactionHandler func(s disgord.Session, h *disgord.InteractionCreate) error

func (b *Bot) AddCommand(command *disgord.CreateApplicationCommand, options BotCommandOptions, handler interactionHandler) {
	b.Commands = append(b.Commands, command)
	b.CommandHandlers[command.Name] = BotCommand{
		Handler: handler,
		Options: options,
	}
}

type ComponentHandler struct {
	Expiry  int64
	Handler interactionHandler
}

func (b *Bot) NewComponentHandler(handler interactionHandler) string {
	randID := b.NewComponentHandlerFromOptions(BotComponentOptions{}, handler)
	return randID
}

func (b *Bot) NewComponentHandlerFromOptions(options BotComponentOptions, handler interactionHandler) string {
	randNum := rand.Intn(999999999999)
	randID := fmt.Sprint(randNum)
	newComponent := BotComponent{
		Handler: handler,
		Options: options,
	}
	b.ActiveComponentHandlers[randID] = newComponent
	return randID
}
