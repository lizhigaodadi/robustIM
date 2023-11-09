package gateWay

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func TestGateWayRun(t *testing.T) {
	RunMain()
}

func TestClient(t *testing.T) {
	for i := 0; i < 100; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:6789")
		if err != nil {
			fmt.Printf("Failed to connect to server; %v\n", err)
			return
		} else {
			fmt.Printf("Success to connect to server;\n")
		}

		defer conn.Close()
		/*Send Message To Server*/
		message := []byte("hello world from client!")
		_, err = conn.Write(message)
		if err != nil {
			fmt.Printf("Failed to send data to server: %v\n", err)
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func TestTimeStampGenerator(t *testing.T) {
	generatorInit()
	for i := 0; i < 100; i++ {
		id, err := generator.NextId()
		if err != nil {
			panic(err)
		}
		log.Printf("id: %d\n", id)
		time.Sleep(200 * time.Millisecond)
	}
}
