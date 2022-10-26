# Discord Bot Template
This package allows you to get up and running with disgord, with application command and handler support built in. It also is really easy to use components!

## To start out with:
```go get github.com/jacks0n9/bot_template```

## To make a bot
```go
package main

import (
    "log"

	"github.com/andersfylling/disgord"
	"github.com/jacks0n9/bot_template"

)

var bot = bot_template.NewBotWithDefault()


func init() {
    // We set the client after the initial creation so the commands can be created properly
    // You may want to use the bot with a config file, in which case, you load the token in a init function or main function
    // these functions run after the commands are registered, using the new bot.
    // it will throw an error if we create a bot with an invalid or empty token, so to be as flexible as possible, please just do this.
	bot.Client = disgord.New(disgord.Config{
		BotToken: "your-token-here (don't put it here though, put it in a config file it's a best practice)",
		Logger:   logger,
		Intents:  disgord.IntentGuilds,
	})
}
func main() {
	err := bot.Run()
	if err != nil {
        log.Errorln("error running bot: "+err.Error())
	}
}

```

## Commands
#### An interaction (command/button) handler in this package has the following function signature:
```go
func(session disgord.Session, interaction *disgord.InteractionCreate) error
```
See how the function returns an error? Returning an error will give the user a simple error message you can specify in the client options (but those messages are already filled in by `NewBotWithDefault`), but also give you the actual error message in your logger.

Commands are normally created in a seperate file in an init function like so:
```go
package main

import (
	"context"

	"github.com/andersfylling/disgord"
)

func init() {
	bot.AddCommand(&disgord.CreateApplicationCommand{
		Name:        "setup",
		Description: "Configure roles and permissions.",
	}, nil, func(session disgord.Session, interaction *disgord.InteractionCreate) error {
		session.SendInteractionResponse(context.Background(),interaction,&disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "This is a response, that's all",
			},
		})
		return nil
	})
}

```
**REMEMBER: Bot is a variable and has to be defined so if it isn't working that might be why**

## Message Components
### How does it work?
A message component is a thing on a message that you can interact with, like a button or select menu. Messages components must be wrapped in one action row to work. Every message component, besides an action row, has a custom id, which is sent to the bot server when it's interacted with. Bot template maps each of these components to a function that is called upon interaction.
### Example:
```go
package main

import (
	"context"

	"github.com/andersfylling/disgord"
)

func init() {
	bot.AddCommand(&disgord.CreateApplicationCommand{
		Name:        "setup",
		Description: "Configure roles and permissions.",
	}, nil, func(session disgord.Session, interaction *disgord.InteractionCreate) error {
		session.SendInteractionResponse(context.Background(), interaction, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "This is a response with components",
				Components: []*disgord.MessageComponent{
					{
						Type: disgord.MessageComponentActionRow,
						Components: []*disgord.MessageComponent{
							{
								Label: "I am a button",
								Style: disgord.Primary,
								CustomID: bot.NewComponentHandler(),
							},
						},
					},
				},
			},
		})
		return nil
	})
}

```
You may be wondering: What does `bot.NewComponentHandler` do?

This function generates a random number and uses it as a custom ID after linking it to your handler.

If you don't want to define the functions inline, you can create another init function:
```go
package main

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/jacks0n9/bot_template"
)

func init() {
	bot.LinkIDsToHandlers(map[string]bot_template.BotComponent{
		"coolbutton": {
			Handler: func(session disgord.Session, interaction *disgord.InteractionCreate) error {
				session.SendInteractionResponse(context.Background(), interaction, &disgord.CreateInteractionResponse{
					Type: disgord.InteractionCallbackChannelMessageWithSource,
					Data: &disgord.CreateInteractionResponseData{
						Content: "You clicked this button!",
					},
				})
				return nil
			}},
	})
}

```
And then in the customID field, you can just put "coolbutton".

## Error handling in commands/components:
As I talked about earlier, error handling is as simple as returning an error and it will be handled for you.

In this package, the default config looks like this:
```go
BotConfig{
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
		}
```
To use your own config refer to this example:
```go
package main

import (
	"github.com/jacks0n9/bot_template"
)

func init() {
	bot_template.NewBotWithConfig(bot_template.BotConfig{...})
}
```
That's pretty much it!