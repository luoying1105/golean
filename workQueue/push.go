package workQueue

import (
	"fmt"
	"log"
	"time"
)

func MakePush() {
	go func() {
		body := "hello world.luoying"
		for i := 0; i < 10; i++ {
			PushData(&body)
			time.Sleep(1 * time.Second)
		}
	}()
	log.Println("get receieve")
	Receive()
	fmt.Println("end")

}
