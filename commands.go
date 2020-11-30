package keybasebot

import (
	"fmt"
	"strings"

	"github.com/kf5grd/keybasebot/pkg/util"
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

// Adapt loops through a set of Adapters and runs them on a given BotAction
func Adapt(b BotAction, adapters ...Adapter) BotAction {
	for i := len(adapters) - 1; i >= 0; i-- {
		b = adapters[i](b)
	}
	return b
}

// MessageType returns an Adapter that restricts a command to a specific message type
func MessageType(typeName string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			if m.Content.TypeName != typeName {
				return false, nil
			}
			return botAction(m, b)
		}
	}
}

// CommandPrefix returns an Adapter that specifies the specific prefix that this command responds to.
// Note that this will often require that MessageType is called _before_ this adapter
// to avoid a panic.
func CommandPrefix(prefix string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			if !strings.HasPrefix(m.Content.Text.Body, prefix) {
				return false, nil
			}
			return botAction(m, b)
		}
	}
}

// ReactionTrigger returns an Adapter that specifies the specific reaction that this command responds to.
// Note that you do not need to use the MessageType adapter when using this as we will
// already be checking to make sure the message type is a reaction.
func ReactionTrigger(trigger string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			if m.Content.TypeName != "reaction" {
				return false, nil
			}
			if m.Content.Reaction.Body != trigger {
				return false, nil
			}
			return botAction(m, b)
		}
	}
}

// MinRole returns an Adapter that restricts a command to users with _at least_ the specified role.
// Note that this _must_ be called _after_ CommandPrefix because this assumes that we
// already know we're executing the provided command.
func MinRole(kb *keybase.Keybase, role string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			if !util.HasMinRole(kb, role, m.Sender.Username, m.ConvID) {
				return true, fmt.Errorf("Your role must be at least %s to do that.", role)
			}
			return botAction(m, b)
		}
	}
}

// FromUser returns an Adapter that only runs a command when sent by a specific user
func FromUser(user string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			if m.Sender.Username != user {
				return false, nil
			}
			return botAction(m, b)
		}
	}
}

// FromUsers returns an Adapter that only runs a command when sent by one of a list of specific users
func FromUsers(users []string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			if !util.StringInSlice(m.Sender.Username, users) {
				return false, nil
			}
			return botAction(m, b)
		}
	}
}

// Contains returns an Adapter that only runs a command when the message contains a specific word.
func Contains(s string, ignoreCase bool, ignoreWhiteSpace bool) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			body := m.Content.Text.Body
			s := s
			if ignoreCase {
				body = strings.ToLower(body)
				s = strings.ToLower(s)
			}
			if ignoreWhiteSpace {
				body = strings.Join(strings.Fields(body), "")
				s = strings.Join(strings.Fields(s), "")
			}
			if !strings.Contains(body, s) {
				return false, nil
			}
			return botAction(m, b)
		}
	}
}

// AdvertiseCommands loops through all the bot's commands and sends their advertisements to the Keybase service
func (b *Bot) AdvertiseCommands() {
	var publicCommands = make([]chat1.UserBotCommandInput, 0)
	for _, ad := range b.Commands {
		if adRes := ad.Ad; adRes != nil {
			publicCommands = append(publicCommands, *adRes)
		}
	}

	public := chat1.AdvertiseCommandAPIParam{
		Typ:      "public",
		Commands: publicCommands,
	}

	publishAds := []chat1.AdvertiseCommandAPIParam{
		public,
	}

	ads := keybase.AdvertiseCommandsOptions{
		Advertisements: publishAds,
	}
	if b.Name != "" {
		ads.Alias = b.Name
	}

	err := b.KB.AdvertiseCommands(ads)
	if err != nil {
		b.Logger.Error("Error setting adverts: %v", err)
	}
}

// ClearCommands clears the advertised commands from the Keybase service
func (b *Bot) ClearCommands() {
	b.KB.ClearCommands()
}
