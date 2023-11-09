package relation

import (
	"context"
	"fmt"
	"im/common/cache"
	"im/common/config"
	"im/common/utils"
	"log"
	"strconv"
)

/*clientId is The UserId in GateWay*/
func UpdateLogin(ctx context.Context, clientId uint64, connId uint64) error {

	slot := config.GetStateServerLoginSlot()
	/*Generate Login Slot Key*/
	loginKey := utils.HashWithSlot(uint32(clientId), slot)
	loginSlotKey := fmt.Sprintf(cache.LoginSlotSetKey, loginKey)
	/*TODO: Generate LoginKey(connId or Did ?) MayBe We Need To optimize*/

	err := cache.SAdd(ctx, loginSlotKey, connId)
	if err != nil {
		log.Printf("Update Login Status For clientId: %d, connId:%d Failed\n", clientId, connId)
		return err
	}

	return nil
}

func UpdateUnLogin(ctx context.Context, clientId uint64, connId uint64) error {
	slot := config.GetStateServerLoginSlot()
	/*Generate Login Slot Key*/
	loginKey := utils.HashWithSlot(uint32(clientId), slot)
	loginSlotKey := fmt.Sprintf(cache.LoginSlotSetKey, loginKey)
	err := cache.SRem(ctx, loginSlotKey, connId)
	if err != nil {
		log.Printf("Update unLogin Status For clientId: %d, connId:%d Failed\n", clientId, connId)
		return err
	}

	return nil
}

func GetLoginStatus(ctx context.Context, clientId uint64) bool {
	/*TODO:Check User is Login?*/
	connIds := GetConnIdsByClientId(ctx, clientId)
	if connIds == nil || len(connIds) == 0 {
		return false /*UnLogin Status*/
	} else {
		return true
	}
}

func AddConnIdToClient(ctx context.Context, clientId, connId uint64) error {
	clientKey := fmt.Sprintf(cache.ClientIdToConnIdKey, clientId)
	err := cache.SAdd(ctx, clientKey, connId)
	if err != nil {
		log.Printf("Add ConnId To ClientId Failed\n")
		return nil
	}

	return nil
}

func GetConnIdsByClientId(ctx context.Context, clientId uint64) []uint64 {
	clientKey := fmt.Sprintf(cache.ClientIdToConnIdKey, clientId)

	member, err := cache.SMemberStringSlice(ctx, clientKey)
	if err != nil {
		return nil
	}
	connIds := make([]uint64, len(member))
	for _, connId := range member {
		id, err := strconv.ParseUint(connId, 10, 64)
		if err != nil {
			log.Printf("Parse Int Failed\n")
			return nil
		}
		connIds = append(connIds, id)
	}

	return connIds
}
