package util

import (
	"fmt"
	"strings"

	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
)

// StringInSlice returns true if the given string is present in the slice of strings
func StringInSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// ChannelString returns a string representation of a chat1.ChatChannel suitable for inclusion in chat and log messages
func ChannelString(channel chat1.ChatChannel) string {
	if channel.MembersType == keybase.USER {
		return channel.Name
	}

	return fmt.Sprintf("%s#%s", channel.Name, channel.TopicName)
}

// HasMinRole returns true if the given user has the given role or higher in the converation
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
