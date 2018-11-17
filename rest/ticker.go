package rest

import (
	"fmt"
	"time"
)

func ticker(){


	fmt.Println("hello")
	for t := range time.NewTicker(10 * time.Second).C {
		heartBeat(t)
	}


}

func heartBeat(tick time.Time){
	//for range time.Tick(time.Second *1){
		fmt.Println("Foo")
	//}
}