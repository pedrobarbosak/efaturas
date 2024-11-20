package main

import (
	"log"
	"time"

	cache "github.com/jellydator/ttlcache/v3"
)

func main() {
	s := *cache.New[string, *int](cache.WithTTL[string, *int](time.Second))
	go s.Start()

	i := 12
	s.Set("key", &i, cache.DefaultTTL)

	time.Sleep(time.Second * 5)

	v := s.Get("key")
	log.Println(v)
}
