package keybasebot

import (
	"keybasebot/pkg/logr"
)

// Run starts the bot listening for new messages
func (b *Bot) Run() error {
	b.Logger = logr.New(b.LogWriter, b.Debug, b.JSON)
	b.registerHandlers()
	b.AdvertiseCommands()
	defer b.ClearCommands()

	b.Logger.Info("Running as user %s", b.KB.Username)
	b.KB.Run(b.Handlers, &b.Opts)
	return nil
}
