package cache

import (
	"context"
	"fmt"
)

const (
	LuaCompareAndIncrClientId = "LuaCompareAndIncrClientId"
)

type luaScript struct {
	LuaScript string
	Sha       string
}

var luaScriptTable = map[string]*luaScript{
	LuaCompareAndIncrClientId: &luaScript{
		LuaScript: "if redis.call('exists',KEYS[1]) end;if redis.call('get', KEYS[1]) == ARGV[1] then redis.call('incr', KEYS[1]); redis.call('expire', KEYS[1],ARGV[2]); return 1 else return -1 end",
	},
}

func initLuaScript(ctx context.Context) {
	for name, script := range luaScriptTable {
		/*Pre Load LuaScript To Generate Sha*/
		cmd := rdb.ScriptLoad(ctx, script.LuaScript)
		if cmd == nil {
			panic(fmt.Sprintf("lua init failed lua:%s", name))
		}
		if cmd.Err() != nil {
			panic(cmd.Err())
		}
		script.Sha = cmd.Val()
	}
}
