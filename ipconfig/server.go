package ipconfig

import (
	"context"
	"im/common/config"
	"im/ipconfig/pkg/domain"
	"im/ipconfig/web/router"
)

func Init() {
	/*TODO: Start all functions of the server*/
	config.Init("../")
	/*dispatcher Init*/
	ctx := context.Background()
	domain.Init(&ctx)

	go router.Init()

}
