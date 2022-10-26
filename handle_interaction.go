package bot_template

import (
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func (b *Bot) handleInteraction(session disgord.Session, interaction *disgord.InteractionCreate) {

	empty := disgord.CreateInteractionResponse{}
	errResp := empty
	defer func() {
		if err := recover(); err != nil {
			b.Client.Logger().Error(fmt.Sprintf("A panic occurred: %v", err))
			session.SendInteractionResponse(context.Background(), interaction, &b.Config.DefaultGeneralErrorMessage)
			return
		}
	}()
	switch interaction.Type {
	case disgord.InteractionApplicationCommand:
		{
			guild, _ := b.Client.Guild(interaction.Member.GuildID).Get()
			cmd := b.CommandHandlers[interaction.Data.Name]
			perms, _ := interaction.Member.GetPermissions(context.Background(), session)
			if !perms.Contains(cmd.Options.RequiredPermission) && interaction.Member.UserID != guild.OwnerID {
				if msg := cmd.Options.PermissionErrorMessage; msg != empty {
					errResp = msg
				} else {
					errResp = b.Config.DefaultPermissionErrorMessage
				}
				break
			}
			err := cmd.Handler(session, interaction)
			if err != nil {
				b.Client.Logger().Error(fmt.Sprintf("Error occured running command %s: %v", interaction.Data.Name, err))
			}
			if msg := cmd.Options.GeneralErrorMessage; msg != empty {
				errResp = msg
			} else {
				errResp = b.Config.DefaultGeneralErrorMessage
			}

		}
	case disgord.InteractionMessageComponent:
		{
			var comp BotComponent
			if futureComp, ok := b.ActiveComponentHandlers[interaction.Data.CustomID]; !ok {
				return
			} else {
				comp = futureComp
			}
			if id := comp.Options.UserLockedTo; id != 0 {
				if id != interaction.Member.UserID {
					if msg := comp.Options.OutsideInteractionErrorMessage; msg != empty {
						errResp = msg
					} else {
						errResp = b.Config.DefaultOutsideInteractionErrorMessage
					}
					break
				}

			}
			if comp.Handler == nil {
				return
			}
			err := comp.Handler(session, interaction)
			if err != nil {
				if msg := comp.Options.GeneralErrorMessage; msg != empty {
					errResp = msg
				} else {
					errResp = b.Config.DefaultGeneralErrorMessage
				}
				b.Client.Logger().Error(fmt.Sprintf("There was an error on a component interaction. Component ID: %s. User:%s. Error:%v.", interaction.Data.CustomID, interaction.Member.UserID.String(), err))
				break
			}
			if comp.Options.OneClickOnly {
				delete(b.ActiveComponentHandlers, interaction.Data.CustomID)
			}
		}
	}
	if errResp != empty {
		session.SendInteractionResponse(context.Background(), interaction, &errResp)
	}
}
