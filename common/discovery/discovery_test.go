package discovery

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/magiconair/properties/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/common/config"
	"testing"
	"time"
)

func TestMod(t *testing.T) {
	e := &EndPointInfo{
		Ip:   "hello",
		Port: "world",
		Meta: make(map[string]interface{}),
	}
	e.Meta["1"] = "fuck"
	e.Meta["2"] = "you"

	buf := e.Marshal()
	message := string(buf)
	fmt.Println("messsage " + message)

	newE := UnMarshal(buf)
	fmt.Printf("%v", newE)
}

func TestEtcdConnect(t *testing.T) {
	DemoInit()
	ctx := context.Background()
	register := NewEtcdServiceRegister(&ctx, "test", &EndPointInfo{
		Ip:   "fuck",
		Port: "you",
		Meta: make(map[string]interface{}),
	}, 5)

	time.Sleep(20 * time.Second)

	register.Close()
}

func DemoInit() {
	err := config.Init("../")
	if err != nil {
		fmt.Printf("config init failed\n")
	}
}

func TestEtcdPutAndSet(t *testing.T) {
	DemoInit()
	/*Connect To Etcd*/
	fmt.Printf("config: %v,write: %v\n", config.GetEtcdEndPointsForDiscovery(), []string{"127.0.0.1:2379"})
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEtcdEndPointsForDiscovery(),
		DialTimeout: config.GetEtcdTimeOutDialTimeForDiscovery(),
	})
	assert.Equal(t, err, nil)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = cli.Put(ctx, "key2", "value1")
	assert.Equal(t, err, nil)
	resp, err := cli.Get(ctx, "test", clientv3.WithPrefix())
	assert.Equal(t, err, nil)
	for _, kv := range resp.Kvs {
		fmt.Printf("key :%s, value: %s\n", kv.Key, kv.Value)
	}
}

func TestEtcdWatch(t *testing.T) {
	DemoInit()

	time.Sleep(5 * time.Second)

	ctx := context.Background()
	/*watcher*/
	discovery := NewEtcdServiceDiscovery(&ctx)

	go func() {
		err := discovery.WatchService("test", setTest, delTest)
		assert.Equal(t, err, nil)
	}()

	go func() { /*register*/
		ctx := context.Background()
		register := NewEtcdServiceRegister(&ctx, "test4", &EndPointInfo{
			Ip:   "hello",
			Port: "world",
			Meta: make(map[string]interface{}),
		}, 0)

		time.Sleep(100 * time.Second)
		register.Close()
	}()

	time.Sleep(100 * time.Second)
}

func setTest(key, val string) {
	fmt.Printf("key: %s,val: %s has been set\n", key, val)
}

func delTest(key, val string) {
	fmt.Printf("key: %s,val: %s has been del\n", key, val)
}

func TestRedisConn(t *testing.T) {
	// 创建 Redis 客户端实例

	client := redis.NewClient(&redis.Options{
		Addr:     "120.76.241.187:6379", // Redis 服务器地址和端口
		Password: "",                    // Redis 服务器密码（如果有的话）
		DB:       0,                     // Redis 数据库索引
	})

	// 使用 Ping 方法测试连接是否成功
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)

	// 使用 SCAN 命令查找所有键
	var cursor uint64 = 0
	for {
		// 执行 SCAN 命令，每次返回一批匹配的键
		keys, nextCursor, err := client.Scan(cursor, "*", 10).Result()
		if err != nil {
			fmt.Printf("Failed to scan keys: %v\n", err)
			return
		}

		// 输出匹配的键
		for _, key := range keys {
			fmt.Println(key)
		}

		// 如果已经遍历完了所有键，则退出循环
		if nextCursor == 0 {
			break
		}

		// 更新游标，以便下一次 SCAN 命令从正确的位置开始
		cursor = nextCursor
	}

	client.Close()
}
