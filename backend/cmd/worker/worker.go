package main

import (
	"fmt"
	
	"github.com/go-redis/redis"
)

type Worker struct {
	r *redis.Client
	topic string
}

func NewWorker(address string, port int, db int, topic string) *Worker  {
	r := redis.NewClient(&redis.Options{
		Addr: address,
		DB: db,
	})

	return &Worker {

	}
}


