package kvstore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"samhofi.us/x/keybase/v2"
)

// Namespaces returns a slice of strings containing all the namespaces for a team
func Namespaces(kb *keybase.Keybase, team string) ([]string, error) {
	var teamName *string

	teamName = &team
	if team == "" {
		teamName = nil
	}

	namespaces, err := kb.KVListNamespaces(teamName)
	if err != nil {
		return []string{}, err
	}
	return namespaces.Namespaces, nil
}

// Keys returns a slice of strings containing all the keys for a namespace
func Keys(kb *keybase.Keybase, team, namespace string) ([]string, error) {
	var teamName *string

	teamName = &team
	if team == "" {
		teamName = nil
	}

	keys, err := kb.KVListKeys(teamName, namespace)
	if err != nil {
		return []string{}, err
	}

	// base64 decode all the keys before returning them
	var ret = make([]string, 0)
	for _, key := range keys.EntryKeys {
		s, _ := base64.StdEncoding.DecodeString(key.EntryKey)
		ret = append(ret, string(s))
	}
	return ret, nil
}

// Get fetches a key from the store
func Get(kb *keybase.Keybase, team, namespace string, kv *KV) error {
	var teamName *string

	teamName = &team
	if team == "" {
		teamName = nil
	}

	key := base64.StdEncoding.EncodeToString([]byte(kv.Key))
	val, err := kb.KVGet(teamName, namespace, key)
	if err != nil {
		return fmt.Errorf("unable to fetch key from store: %v", err)
	}

	value, err := base64.StdEncoding.DecodeString(val.EntryValue)
	if err != nil {
		return fmt.Errorf("unable to base64 decode value data from store: %v", err)
	}
	err = json.Unmarshal(value, kv.Value)
	if err != nil {
		return fmt.Errorf("unable to unmarshal value data from store: %v", err)
	}
	return nil
}

// Put writes a key to the store
func Put(kb *keybase.Keybase, team, namespace string, kv KV) error {
	var teamName *string

	teamName = &team
	if team == "" {
		teamName = nil
	}

	// base64 encode key
	key := base64.StdEncoding.EncodeToString([]byte(kv.Key))

	// base64 encode value
	jsonBytes, err := json.Marshal(kv.Value)
	if err != nil {
		return fmt.Errorf("unable to marshal value data: %v", err)
	}
	value := base64.StdEncoding.EncodeToString(jsonBytes)

	// If revision is not nil, we need to pass the revision when deleting from the store
	if kv.Revision != nil {
		_, err = kb.KVPutWithRevision(teamName, namespace, key, value, *kv.Revision)
		return err
	}
	_, err = kb.KVPut(teamName, namespace, key, value)
	return err
}

// Delete deletes a key from the store
func Delete(kb *keybase.Keybase, team, namespace string, kv KV) error {
	var (
		teamName *string
		key      string
	)

	// We base64 encode keys before putting them into the store, so we need to encode it here
	// before deleting
	key = base64.StdEncoding.EncodeToString([]byte(kv.Key))

	// If the team string is empty we'll pass nil to the kvstore in order to use the implicit
	// self-team
	teamName = &team
	if team == "" {
		teamName = nil
	}

	// If revision is not nil, we need to pass the revision when deleting from the store
	if kv.Revision != nil {
		_, err := kb.KVDeleteWithRevision(teamName, namespace, key, *kv.Revision)
		return err
	}
	_, err := kb.KVDelete(teamName, namespace, key)
	return err
}
