package keybasebot

import (
	"github.com/kf5grd/keybasebot/pkg/logr"
)

// Run starts the bot listening for new messages
func (b *Bot) Run() error {
	// set up logger
	logWriter := convWriter{
		// if convID is empty (which is the default) then logs will only be written to stdout,
		// but if a conversation id is passed here then logs will be written to stdout *and*
		// this conversation
		ConvID: b.LogConv,
		Writer: b.LogWriter,
		KB:     b.KB,
	}
	b.Logger = logr.New(logWriter, b.Debug, b.JSON)

	b.registerHandlers()
	b.AdvertiseCommands()
	defer b.ClearCommands()

	b.Logger.Info("Running as user %s", b.KB.Username)
	b.running = true
	b.KB.Run(b.Handlers, &b.Opts)
	b.running = false

	return nil
}

// Running indicates whether the bot is currently running
func (b *Bot) Running() bool {
	return b.running
}
