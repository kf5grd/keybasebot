package keybasebot

import (
	"fmt"
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
	// Name of the command for use in the logs
	Name string

	// This will show in your chat channels when someone starts typing a command in the
	// message box. Set this to nil if you don't want the command advertised
	Ad *chat1.UserBotCommandInput

	// AdType can be one of "public", "teamconvs", and "teammembers". If AdType is empty or
	// unknown, it will default to "public". If AdType is one of "teamconvs" or "teammembers",
	// be sure to specify the corresponding team name in AdTeamName. Note: These settings only
	// restrict where the advertisements will show. You will still need to either use
	// appropriate Adapters, or write your BotAction in a way that limits where the commands
	// can be called from if that is your intention
	AdType string

	// If AdType is one of "teamconvs" or "teammembers", be sure to enter a team name in
	// AdTeamName, which will restrict your commands to either be advertised in this team, or
	// only be advertised to members of this team
	AdTeamName string

	// The function to run when the command is triggered
	Run BotAction
}

// Adapter can modify the behavior of a BotAction
type Adapter func(BotAction) BotAction

// Bot is where we'll hold the necessary information for the bot to run
type Bot struct {
	// The Name string will show up next to the bot's username in chat messages. Setting this
	// to empty will cause it to be ignored
	Name string

	// If CommandPrefix is set, any messages that are received with a TypeName of "text," and
	// do not have this string prefix will be discarded. This can be useful if all of your
	// commands start with the same prefix.
	CommandPrefix string

	// The Keybase instance
	KB *keybase.Keybase

	// The logr instance
	Logger *logr.Logger

	// Where you want log messages written to
	LogWriter io.Writer

	// If LogConv is not empty, log messages will be sent to this conversation, in addition to
	// the LogWriter
	LogConv chat1.ConvIDStr

	// Whether log messages should be in JSON format
	JSON bool

	// Whether to show debug messages in log output
	Debug bool

	// Message handlers. You probably should leave the Chat handler alone
	Handlers keybase.Handlers

	// Custom run options for the message listener
	Opts keybase.RunOptions

	// A slice holding all of you BotCommands. Be sure to populate this prior to calling Run()
	Commands []BotCommand

	// You can use this to store custom info in order to pass it around to your bot commands
	Meta map[string]interface{}
}

// New returns a new Bot instance. name will set the Bot.Name and will show up next to the
// bot's username in chat messages. You can set name to an empty string. You can also pass
// in any keybase.KeybaseOpt options and they will be passed to keybase.New()
func New(name string, opts ...keybase.KeybaseOpt) *Bot {
	var b Bot
	b.Name = name
	b.KB = keybase.New(opts...)
	b.Handlers = keybase.Handlers{}
	b.Opts = keybase.RunOptions{}
	b.Commands = make([]BotCommand, 0)
	b.Meta = make(map[string]interface{})

	return &b
}

// We'll use this to create a writer for the Logger which will be able to write logs to
// stdout, and optionally also to a Keybase chat conversation
type convWriter struct {
	ConvID chat1.ConvIDStr
	Writer io.Writer
	KB     *keybase.Keybase
}

// Write sends log message strings to a channel
func (cw convWriter) Write(p []byte) (n int, err error) {
	if cw.ConvID != "" {
		opts := keybase.SendMessageOptions{
			ConversationID: cw.ConvID,
			NonBlock:       true,
			Message:        keybase.SendMessageBody{Body: string(p)},
		}
		cw.KB.SendMessage("send", opts)
	}

	fmt.Fprintf(cw.Writer, string(p))
	return len(p), nil
}
