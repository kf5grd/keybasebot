package keybasebot

import (
	"fmt"
	"strings"

	"github.com/kf5grd/keybasebot/pkg/util"
	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

// Adapt loops through a set of Adapters and runs them on a given BotAction in the order
// that they're provided. It's important to make sure you're passing adapters in the
// correct order as some things will need to be checked before others. As an example, the
// Contains adapter assumes that the incoming message has a MessageType of "text." If we
// pass that adapter prior to passing the MessageType adapter, then we will end up with a
// panic any time a message comes through with a MessageType other than "text."
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
			b.Logger.Debug("Verifying message type is '%s'", typeName)
			if m.Content.TypeName != typeName {
				b.Logger.Debug("Message type is '%s', exiting command", m.Content.TypeName)
				return false, nil
			}
			b.Logger.Debug("Message type is '%s', continuing", typeName)
			return botAction(m, b)
		}
	}
}

// CommandPrefix returns an Adapter that specifies the specific prefix that this command
// responds to. Note that this will often require that MessageType is called _before_ this
// adapter to avoid a panic.
func CommandPrefix(prefix string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			b.Logger.Debug("Verifying message contains prefix '%s'", prefix)
			if !strings.HasPrefix(m.Content.Text.Body, prefix) {
				b.Logger.Debug("Message does not contain prefix '%s', exiting command", prefix)
				return false, nil
			}
			b.Logger.Debug("Message does contain prefix '%s', continuing", prefix)
			return botAction(m, b)
		}
	}
}

// ReactionTrigger returns an Adapter that specifies the specific reaction that this
// command responds to. Note that you do not need to use the MessageType adapter when using
// this as we will already be checking to make sure the message type is a reaction.
func ReactionTrigger(trigger string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			b.Logger.Debug("Verifying message type is 'reaction'")
			if m.Content.TypeName != "reaction" {
				b.Logger.Debug("Message type is '%s', exiting command", m.Content.TypeName)
				return false, nil
			}
			b.Logger.Debug("Verifying reaction body is '%s'", trigger)
			if m.Content.Reaction.Body != trigger {
				b.Logger.Debug("Reaction body is '%s', exiting command", m.Content.Reaction.Body)
				return false, nil
			}
			b.Logger.Debug("Reaction body is '%s', continuing", m.Content.Reaction.Body)
			return botAction(m, b)
		}
	}
}

// MinRole returns an Adapter that restricts a command to users with _at least_ the
// specified role. Note that this _must_ be called _after_ CommandPrefix because this
// assumes that we already know we're executing the provided command.
func MinRole(kb *keybase.Keybase, role string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			b.Logger.Debug("Verifying user '%s' has minimum role '%s' in '%s'", m.Sender.Username, role, util.ChannelString(m.Channel))
			if !util.HasMinRole(kb, role, m.Sender.Username, m.ConvID) {
				b.Logger.Debug("User '%s' does not have minimum role '%s' in '%s', exiting command and replying with error", m.Sender.Username, role, util.ChannelString(m.Channel))
				return true, fmt.Errorf("Your role must be at least %s to do that.", role)
			}
			b.Logger.Debug("User '%s' has minimum role '%s' in '%s', continuing", m.Sender.Username, role, util.ChannelString(m.Channel))
			return botAction(m, b)
		}
	}
}

// FromUser returns an Adapter that only runs a command when sent by a specific user
func FromUser(user string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			b.Logger.Debug("Verifying received message was sent by '%s'", user)
			if m.Sender.Username != user {
				b.Logger.Debug("Received message was sent by '%s', exiting command", m.Sender.Username)
				return false, nil
			}
			b.Logger.Debug("Received message was sent by '%s', continuing", user)
			return botAction(m, b)
		}
	}
}

// FromUsers returns an Adapter that only runs a command when sent by one of a list of
// specific users
func FromUsers(users []string) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			b.Logger.Debug("Verifying received message was sent by one of '%s'", strings.Join(users, ","))
			if !util.StringInSlice(m.Sender.Username, users) {
				b.Logger.Debug("Received message was sent by '%s', exiting command", m.Sender.Username)
				return false, nil
			}
			b.Logger.Debug("Received message was sent by '%s', continuing", m.Sender.Username)
			return botAction(m, b)
		}
	}
}

// Contains returns an Adapter that only runs a command when the message contains a
// specific string. This will also make sure the received message has a type of 'text' or
// 'edit'
func Contains(s string, ignoreCase bool, ignoreWhiteSpace bool) Adapter {
	return func(botAction BotAction) BotAction {
		return func(m chat1.MsgSummary, b *Bot) (bool, error) {
			b.Logger.Debug("Verifying message contains '%s'", s)
			var body string

			switch m.Content.TypeName {
			case "text":
				body = m.Content.Text.Body
			case "edit":
				body = m.Content.Edit.Body
			default:
				b.Logger.Debug("Received message does not have type 'text' or 'edit', exiting command")
				return false, nil
			}

			var s = s
			if ignoreCase {
				body = strings.ToLower(body)
				s = strings.ToLower(s)
			}
			if ignoreWhiteSpace {
				body = strings.Join(strings.Fields(body), "")
				s = strings.Join(strings.Fields(s), "")
			}
			if !strings.Contains(body, s) {
				b.Logger.Debug("Message does not contain '%s', exiting command", s)
				return false, nil
			}
			b.Logger.Debug("Message does contain '%s', continuing", s)
			return botAction(m, b)
		}
	}
}

// AdvertiseCommands loops through all the bot's commands and sends their advertisements
// to the Keybase service
func (b *Bot) AdvertiseCommands() {
	var publicCommands = make([]chat1.UserBotCommandInput, 0)
	for _, ad := range b.Commands {
		if adRes := ad.Ad; adRes != nil {
			publicCommands = append(publicCommands, *adRes)
		}
	}

	if len(publicCommands) == 0 {
		b.Logger.Debug("Bot has no command advertisements")
		return
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
