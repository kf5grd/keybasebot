package keybasebot

import (
	"io"

	"github.com/kf5grd/keybasebot/pkg/logr"
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

// BotAction is a function that's run when a command is received by the bot. If the boolean
// return is true, the bot will not attempt to execute any other commands after this one.
// If an error is returned, it will be sent to the logger. If an error is returned and the
// boolean is also set to true, the returned error will be sent back to the chat as a reply
// to the message that triggered the command.
type BotAction func(chat1.MsgSummary, *Bot) (bool, error)

// BotCommand holds information regarding a command and its advertisements
type BotCommand struct {
	Name string                     // Name of the command for use in the logs
	Ad   *chat1.UserBotCommandInput // This will show in your chat channels when someone starts typing a command in the message box
	Run  BotAction                  // The function to run when the command is triggered
}

// Adapter can modify the behavior of a BotAction
type Adapter func(BotAction) BotAction

// Bot is where we'll hold the necessary information for the bot to run
type Bot struct {
	Name          string             // This will show up next to the bot's username in chat messages
	CommandPrefix string             // If this is set, any messages that are received with a Text type that do not have this string prefix will be discarded. This can be useful if all of your commands start with the same prefix.
	KB            *keybase.Keybase   // The Keybase instance
	Logger        *logr.Logger       // The logr instance
	LogWriter     io.Writer          // Where you want log messages sent
	JSON          bool               // Whether log messages should be in JSON format
	Debug         bool               // Whether to show debug messages in log output
	Handlers      keybase.Handlers   // Message handlers. You probably should leave the Chat handler alone
	Opts          keybase.RunOptions // Custom run options for the message listener
	Commands      []BotCommand       // A slice holding all of you BotCommands. Be sure to populate this prior to calling Run()
}

// New returns a new Bot instance. name will set the Bot.Name and will show up next to the bot's username in chat messages. You can set name to an empty string.
// You can also pass in any keybase.KeybaseOpt options and they will be passed to keybase.New()
func New(name string, opts ...keybase.KeybaseOpt) *Bot {
	var b Bot
	b.Name = name
	b.KB = keybase.New(opts...)
	b.Handlers = keybase.Handlers{}
	b.Opts = keybase.RunOptions{}
	b.Commands = make([]BotCommand, 0)

	return &b
}
