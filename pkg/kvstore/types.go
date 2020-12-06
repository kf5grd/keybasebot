package kvstore

// KV holds a key/value pair for use with the kvstore. The Value will be marshaled to and
// from JSON.
type KV struct {
	Key      string
	Value    interface{}
	Revision *int
}

// New returns a KV struct. Passing a negative revision will cause the revision number
// to be ignored.
func New(key string, value interface{}, revision int) KV {
	kv := KV{
		Key:   key,
		Value: value,
	}

	if revision < 0 {
		return kv
	}
	kv.Revision = &revision
	return kv

}
