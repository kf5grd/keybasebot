package util

import (
	"fmt"
	"strings"

	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

func ChannelString(channel chat1.ChatChannel) string {
	if channel.MembersType == keybase.USER {
		return channel.Name
	}

	return fmt.Sprintf("%s#%s", channel.Name, channel.TopicName)
}

func HasMinRole(kb *keybase.Keybase, role string, user string, conv chat1.ConvIDStr) bool {
	conversation, err := kb.ListMembersOfConversation(conv)
	if err != nil {
		return false
	}

	memberTypes := make(map[string]struct{})
	memberTypes["owner"] = struct{}{}
	memberTypes["admin"] = struct{}{}
	memberTypes["writer"] = struct{}{}
	memberTypes["reader"] = struct{}{}

	if _, ok := memberTypes[strings.ToLower(role)]; !ok {
		// invalid role
		return false
	}

	for _, member := range conversation.Owners {
		if strings.ToLower(member.Username) == strings.ToLower(user) {
			return true
		}
	}
	if strings.ToLower(role) == "owner" {
		return false
	}

	for _, member := range conversation.Admins {
		if strings.ToLower(member.Username) == strings.ToLower(user) {
			return true
		}
	}
	if strings.ToLower(role) == "admin" {
		return false
	}

	for _, member := range conversation.Writers {
		if strings.ToLower(member.Username) == strings.ToLower(user) {
			return true
		}
	}
	if strings.ToLower(role) == "writer" {
		return false
	}

	for _, member := range conversation.Readers {
		if strings.ToLower(member.Username) == strings.ToLower(user) {
			return true
		}
	}
	if strings.ToLower(role) == "reader" {
		return false
	}

	return false
}
