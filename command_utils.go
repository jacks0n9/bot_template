package bot_template

import (
	"fmt"
	"math/rand"

	"github.com/andersfylling/disgord"
)

type InteractionHandler func(session disgord.Session, interaction *disgord.InteractionCreate) error

func (b *Bot) AddCommand(command *disgord.CreateApplicationCommand, options BotCommandOptions, handler InteractionHandler) {
	b.Commands = append(b.Commands, command)
	b.CommandHandlers[command.Name] = BotCommand{
		Handler: handler,
		Options: options,
	}
}

type ComponentHandler struct {
	Expiry  int64
	Handler InteractionHandler
}

func (b *Bot) NewComponentHandler(handler InteractionHandler) string {
	randID := b.NewComponentHandlerFromOptions(BotComponentOptions{}, handler)
	return randID
}

func (b *Bot) NewComponentHandlerFromOptions(options BotComponentOptions, handler InteractionHandler) string {
	randNum := rand.Intn(999999999999)
	// Append prefix so you don't accidentally add a custom id link thing
	randID := fmt.Sprint(randNum) + "-g"
	newComponent := BotComponent{
		Handler: handler,
		Options: options,
	}
	b.LinkIDToHandler(randID, newComponent)
	return randID
}
func (b *Bot) LinkIDToHandler(ID string, comp BotComponent) {
	b.ActiveComponentHandlers[ID] = comp
}
func (b *Bot) LinkIDsToHandlers(compMap map[string]BotComponent) {
	for name, comp := range compMap {
		b.ActiveComponentHandlers[name] = comp
	}
}
