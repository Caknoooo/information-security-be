package utils

import (
	"log"
	"time"
)

// func timer(name string) {
// 	start := time.Now()
// 	log.Printf("Start %s", name)
// 	defer func() {
// 		log.Printf("End %s, took %v", name, time.Since(start))
// 	}()
// }

func timerWithReturn(name string) func() time.Duration {
	start := time.Now()
	log.Printf("Start %s", name)
	return func() time.Duration {
		log.Printf("End %s, took %v", name, time.Since(start))
		return time.Since(start)
	}
}