package main

import (
	"log"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/go-redis/redis"
)

var (
	PublishAddressFlag = cli.StringFlag {
		Name: "address",
		Usage: "address of Redis Stream to publish to",
		Aliases: []string{"a", "addr"},
		Value: "127.0.0.1",
	}
	PublishPortFlag = cli.IntFlag {
		Name: "port",
		Usage: "port of Redis Stream to publish to",
		Aliases: []string{"p"},
		Value: 6379,
	}
	PublishTopicFlag = cli.StringFlag {
		Name: "topic",
		Usage: "topic to send message to",
		Aliases: []string{"t", "c"},
		Value: "aida",
	}
	PublishMasterIdFlag = cli.StringFlag {
		Name: "master-id",
		Usage: "master id",
		Aliases: []string{"m"},
	}
)

func PublishRegistrar(ctx *cli.Context) error {
	return publishRegistrar(
		ctx.String("address"),
		ctx.Int("port"),
		ctx.String("topic"),
	)
}

func publishRegistrar(addr string, port int, topic string) error {
	pub, err := MakeRedisPublisher(addr, port, topic)
	if err != nil {
		log.Fatal("Unable to connect to Redis: ", err)
		return err
	}

	err = pub.Publish(map[string]any{
		"start": 1,
		"end": 1,
		"memory": 1,
		"txCount": 1,
		"gas": 1,
		"totalTxCount": 2,
		"totalGas": 2,
		"lDisk": 3,
		"aDisk": 3,
	})
	if err != nil {
		log.Fatal("Unable to publish to Redis")
		return err
	}

	return nil
}

type Message map[string]any
type msg Message

type Publisher interface {
	Publish(msg) error
}

type NilPublisher struct {}
func (_ *NilPublisher) Publish(m msg) error {
	return nil
}

type RedisPublisher struct {
	NilPublisher
	r *redis.Client
	topic string
}

func MakeRedisPublisher(addr string, port int, topic string) (*RedisPublisher, error) {
	r := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", addr, port),
	})

	_, err := r.Ping().Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis", err)
		return nil, err
	}

	log.Println("Connected to Redis server", addr, port)

	return &RedisPublisher{r: r, topic: topic}, nil
}

func (pub *RedisPublisher) Publish(m msg) error {
	err := pub.r.XAdd(&redis.XAddArgs{
		Stream: pub.topic,
		Values: m,
	}).Err()

	log.Println("msg sent!")

	return err
}
