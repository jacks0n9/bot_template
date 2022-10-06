package bot_template

import (
	"fmt"
	"math/rand"

	"github.com/andersfylling/disgord"
)

type interactionHandler func(s disgord.Session, h *disgord.InteractionCreate)

func (b *Bot) AddCommand(command *disgord.CreateApplicationCommand, handler interactionHandler) {
	b.Commands = append(b.Commands, command)
	b.ActiveCommandHandlers[command.Name] = handler
}
func (b *Bot) NewComponentHandler(handler interactionHandler) string {
	randNum := rand.Intn(99999999999)
	randID := fmt.Sprint(randNum)
	b.ComponentHandlers[randID] = handler
	return randID
}
