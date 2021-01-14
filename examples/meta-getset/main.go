// This is a very simple bot that has 2 commands: set, and get. The set command sets a
// string variable named "message" in the Meta store, and the get command retrieves that
// variable and sends it to the user in a chat message.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	bot "github.com/kf5grd/keybasebot"
	"github.com/kf5grd/keybasebot/pkg/util"
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

// We'll use this to create a writer for the Logger which will be able to write logs to
// stdout, and optionally also to a Keybase chat conversation
type kbWriter struct {
	convID chat1.ConvIDStr
	bot    *bot.Bot
}

func (k kbWriter) Write(p []byte) (n int, err error) {
	opt := keybase.SendMessageOptions{
		ConversationID: k.convID,
		Message:        keybase.SendMessageBody{Body: string(p)},
	}
	go k.bot.KB.SendMessage("send", opt)
	fmt.Fprintf(os.Stdout, string(p))
	return len(p), nil
}

func main() {
	// setup flags
	var homePath = flag.String("home", "", "Set custom home directory for the Keybase client")
	var debug = flag.Bool("debug", false, "Enable debuging output")
	var json = flag.Bool("json", false, "Output logs in JSON format")
	var logConv = flag.String("log-conv", "", "Conversation ID to send logs to")
	flag.Parse()

	// setup bot
	b := bot.New("", keybase.SetHomePath(*homePath))
	b.Debug = *debug
	b.JSON = *json

	// setup the log writer
	w := kbWriter{
		// if convID is empty (which is the default) then logs will only be written to stdout,
		// but if a conversation id is passed here then logs will be written to stdout *and*
		// this conversation
		convID: chat1.ConvIDStr(*logConv),
		bot:    b,
	}
	b.LogWriter = w

	// register the bot's commands
	b.Commands = append(b.Commands,
		bot.BotCommand{
			Name: "SetMessage",
			Ad:   &setMessageAd,
			Run: bot.Adapt(setMessage,
				// this command can only be triggered by messages with
				// the "text" type...
				bot.MessageType("text"),

				// ...it will only be triggered if the message has this prefix
				bot.CommandPrefix("!set"),
			),
		},
		bot.BotCommand{
			Name: "GetMessage",
			Ad:   &getMessageAd,
			Run: bot.Adapt(getMessage,
				// this command can only be triggered by messages with
				// the "text" type...
				bot.MessageType("text"),

				// ...it will only be triggered if the message has this prefix
				bot.CommandPrefix("!get"),
			),
		},
	)

	// run bot
	b.Run()
}

// Advertisements for setMessage
var setMessageAd = chat1.UserBotCommandInput{
	Name:        "set",
	Usage:       "<message>",
	Description: "Set a message that can be displayed with the `!get` command",
}

// setMessage stores our message
func setMessage(m chat1.MsgSummary, b *bot.Bot) (bool, error) {
	message := strings.TrimSpace(strings.Replace(m.Content.Text.Body, "!set", "", 1))
	if message == "" {
		err := fmt.Errorf("Must provide a message.")
		b.Logger.Error("Error setting message value from '%s' in '%s': %v", m.Sender.Username, util.ChannelString(m.Channel), err)
		return true, err
	}

	// store the message
	b.Meta["message"] = message

	// send a reaction to the user letting them know we've processed the command
	b.KB.ReactByConvID(m.ConvID, m.Id, ":heavy_check_mark:")

	// setting this to true means the bot won't look for
	// any more commands to execute after this one runs
	return true, nil
}

// Advertisements for getMessage
var getMessageAd = chat1.UserBotCommandInput{
	Name:        "get",
	Description: "Get the message that was set with the `!set` command",
}

// getMessage retrieves our message
func getMessage(m chat1.MsgSummary, b *bot.Bot) (bool, error) {
	// fetch the message
	message, ok := b.Meta["message"]
	if !ok {
		// setting this to true and returning an error means the bot won't
		// look for any more commands to execute after this one runs, and
		// it will reply to the user with the error message. if we set the
		// boolean to false and return an error, the erro message gets sent
		// to the logs, but does not get sent to the user, and the bot
		// continues to loop through each of the commands trying to run them
		return true, fmt.Errorf("No message has been set yet. Send `!set <message>` to set one.")
	}

	// if we get this far it means there was a message set,
	// and we reply to the user with the message
	b.KB.ReplyByConvID(m.ConvID, m.Id, message.(string))
	return false, nil
}
