package main


import (
"fmt"
"log"
"time"
)

func main() {
	interval := float64(10000)

	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	go func() {
		counter := 10.0
		for {
			select {
			case <-ticker.C:
				log.Println("ticker accelerating to " + fmt.Sprint(interval/counter) + " ms")
				ticker = time.NewTicker(time.Duration(interval/counter) * time.Millisecond)
				counter++
			}
		}
		log.Println("stopped")
	}()
	time.Sleep(50 * time.Second)
	log.Println("stopping ticker")
	ticker.Stop()
}