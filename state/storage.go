package state

import (
	"context"
	"fmt"
	"im/common/cache"
	"strconv"
)

func StorageSession(ctx context.Context, session Session, userIds []uint64) error {
	key := fmt.Sprintf(cache.SessionStorageKey, session.sessionId)

	for i := 1; i < len(userIds); i++ {
		err := cache.SAdd(ctx, key, userIds[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func GetSessionUserIds(ctx context.Context, sessionId uint64) ([]uint64, error) {
	key := fmt.Sprintf(cache.SessionStorageKey, sessionId)
	slice, err := cache.SMemberStringSlice(ctx, key)
	if err != nil {
		return nil, err
	}
	res := make([]uint64, len(slice))

	for i := 0; i < len(slice); i++ {
		userId, err := strconv.ParseInt(slice[i], 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, uint64(userId))
	}

	return res, nil
}
