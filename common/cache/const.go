package cache

import "time"

const (
	MaxClientIdKey      = "max_client_id_{%d}_%d"
	LastMsgKey          = "last_msg_{%d}_%d"
	LoginSlotSetKey     = "login_slot_set_{%d}"
	SessionStorageKey   = "session_key_{%d}"
	ClientIdToConnIdKey = "client_id_to_conn_id_{%d}"
	TTL7D               = 7 * 24 * time.Hour
)
