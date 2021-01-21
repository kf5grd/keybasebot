package keybasebot

import (
	"strings"

	"github.com/kf5grd/keybasebot/pkg/util"
	"samhofi.us/x/keybase/v2/types/chat1"
)

func (b *Bot) registerHandlers() {
	chat := b.chatHandler
	b.Handlers.ChatHandler = &chat
}

func (b *Bot) chatHandler(m chat1.MsgSummary) {
	var (
		sender  = m.Sender.Username
		channel = util.ChannelString(m.Channel)
	)

	// If message comes from the bot, and b.AllowSelfMessages is false, ignore the message
	if sender == b.KB.Username && !b.AllowSelfMessages {
		return
	}

	// If CommandPrefix is set and message is a text message, make sure it has the
	// correct prefix
	if b.CommandPrefix != "" && m.Content.TypeName == "text" {
		if !strings.HasPrefix(m.Content.Text.Body, b.CommandPrefix) {
			return
		}
	}

	// Cycle through each action and run them until we reach the end, or until a command
	// requests to stop execution of subsequent commands
	b.Logger.Debug("Incoming message from %s", sender)
	for _, action := range b.Commands {
		actionName := action.Name
		b.Logger.Debug("Trying %s", actionName)
		ok, err := action.Run(m, b)
		if err != nil {
			b.Logger.Error("[%v][%s in %s] %s returned error: %v", m.ConvID, sender, channel, actionName, err)
			if ok {
				b.KB.ReplyByConvID(m.ConvID, m.Id, err.Error())
			}
		}
		if ok {
			b.Logger.Debug("%s ok = true, cancelling execution of subsequent commands", actionName)
			return
		}
	}
}
