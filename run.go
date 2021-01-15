package keybasebot

import (
	"github.com/kf5grd/keybasebot/pkg/logr"
)

// sendLogMessages pulls log messages from the channel and writes them to a Keybase conversation
func (b *Bot) sendLogMessages(ch chan string) {
	for {
		b.KB.SendMessageByConvID(b.LogConv, <-ch)
	}
}

// Run starts the bot listening for new messages
func (b *Bot) Run() error {
	// set up logger
	logWriter := convWriter{
		// if convID is empty (which is the default) then logs will only be written to stdout,
		// but if a conversation id is passed here then logs will be written to stdout *and*
		// this conversation
		ConvID: b.LogConv,
		ch:     make(chan string, 100),
		Writer: b.LogWriter,
	}
	go b.sendLogMessages(logWriter.ch)
	b.Logger = logr.New(logWriter, b.Debug, b.JSON)

	b.registerHandlers()
	b.AdvertiseCommands()
	defer b.ClearCommands()

	b.Logger.Info("Running as user %s", b.KB.Username)
	b.KB.Run(b.Handlers, &b.Opts)
	return nil
}
