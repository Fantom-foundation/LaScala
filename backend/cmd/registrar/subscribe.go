package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/urfave/cli/v2"
	"github.com/go-redis/redis"
	"github.com/rs/xid"
)

const (
	CreateTableIfNotExist = `
		CREATE TABLE IF NOT EXISTS %s (
			start INTEGER NOT NULL,
			end INTEGER NOT NULL,
			memory int,
			live_disk int,
			archive_disk int,
			tx_rate float,
			gas_rate float,
			overall_tx_rate float,
			overall_gas_rate float
		)
	`

	InsertOrReplace = `
		INSERT or REPLACE INTO %s (
			start, end,
			memory, live_disk, archive_disk,
			tx_rate, gas_rate, overall_tx_rate, overall_gas_rate
		) VALUES (
			?, ?,
			?, ?, ?,
			?, ?, ?, ?
		)
	`
)

func SubscribeRegistrar(ctx *cli.Context) error {
	return subscribeRegistrar(
		ctx.String("address"),
		ctx.Int("port"),
		ctx.String("topic"),
		ctx.String("master-id"),
	)
}

func subscribeRegistrar(addr string, port int, topic string, masterId string) error {
	sub, err := MakeRedisSubscriber(addr, port, masterId)
	if err != nil {
		log.Fatal("unable to start subscriber at ", addr, port, masterId)
	}
	log.Println("subscriber started: ", addr, port, masterId)

	sub.Subscribe(topic)
	return nil
}

type Subscriber interface {
	Subscribe(string) error
}

type NilSubscriber struct {}
func (_ *NilSubscriber) Subscribe(topic string) error {
	return nil
}

type RedisSubscriber struct {
	NilSubscriber
	r *redis.Client
	ps *Printers

	// who am i?
	masterId string
	runId string

	// load
	p progress
}

func MakeRedisSubscriber(addr string, port int, masterId string) (*RedisSubscriber, error) {
	r := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", addr, port),
	})

	_, err := r.Ping().Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis ", err)
		return nil, err
	}
	log.Println("Connected to Redis server ", addr, port)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get current directory")
	}
	log.Println("Using ", wd, " as deposit directory for dbs")

	if masterId == "" {
		masterId = xid.New().String()
	}
	runId := xid.New().String()

	sub := &RedisSubscriber {
		r: r, 
		masterId: masterId, 
		runId: runId, 
		ps: NewPrinters(),
	} 

	db := filepath.Join(wd, fmt.Sprintf("%s.db", masterId))
	p2db, err := NewPrinterToSqlite3(sub.sqlite3(db, runId))
	if err != nil {
		log.Fatal("Unable to connect to ", db, " and access table", runId)
	}
	sub.ps.AddPrinter(p2db)

	return sub, nil
}

func (sub *RedisSubscriber) Subscribe(topic string) error {
	cg := fmt.Sprintf("%s-listener", topic)

	err := sub.r.XGroupCreate(topic, cg, "0").Err()
	if err != nil {
		log.Println(err)
		return err
	}

	uid := xid.New().String()

	for {
		entries, err := sub.r.XReadGroup(&redis.XReadGroupArgs{
			Group: cg,
			Consumer: uid,
			Streams: []string{topic, ">"},
			Count: 2,
			Block: 0,
			NoAck: false,
		}).Result()

		if err != nil {
			log.Fatal("Could not read entries", err)
			return err
		}

		for i := 0; i < len(entries[0].Messages); i++ {
			msgId := entries[0].Messages[i].ID
			values := entries[0].Messages[i].Values

			err := sub.handleProgress(values)
			if err != nil {
				log.Fatal("Could not handle event", err)
				return err
			}

			sub.r.XAck(topic, cg, msgId)
		}
	}

	return nil
}

func (sub *RedisSubscriber) Print() error {
	sub.ps.Print()
	return nil
}

func (sub *RedisSubscriber) handleProgress(values map[string]any) error {
	start, _ := strconv.ParseUint(values["start"].(string), 10, 64)
	end, _ := strconv.ParseUint(values["end"].(string), 10, 64)
	memory, _ := strconv.ParseUint(values["memory"].(string), 10, 64)
	txCount, _ := strconv.ParseUint(values["txCount"].(string), 10, 64)
	gas, _ := strconv.ParseUint(values["gas"].(string), 10, 64)
	totalTxCount, _ := strconv.ParseUint(values["totalTxCount"].(string), 10, 64)
	totalGas, _ := strconv.ParseUint(values["totalGas"].(string), 10, 64)
	lDisk, _ := strconv.ParseUint(values["lDisk"].(string), 10, 64)
	aDisk, _ := strconv.ParseUint(values["aDisk"].(string), 10, 64)

	sub.p = progress{
		start: start,
		end: end,
		memory: memory,
		txCount: txCount,
		gas: gas,
		totalTxCount: totalTxCount,
		totalGas: totalGas,
		lDisk: lDisk,
		aDisk: aDisk,
	}
	
	log.Println(sub.p)
	return sub.Print()
}

func (sub *RedisSubscriber) sqlite3(conn string, tableName string) (string, string, string, func() [][]any) {
	return 	conn, 
		fmt.Sprintf(CreateTableIfNotExist, tableName),
		fmt.Sprintf(InsertOrReplace, tableName),
		func() [][]any {
			return [][]any{{
					sub.p.start,
					sub.p.end,
					sub.p.memory,
					sub.p.txCount,
					sub.p.gas,
					sub.p.totalTxCount,
					sub.p.totalGas,
					sub.p.lDisk,
					sub.p.aDisk,
			}}
		}}

type progress struct {
	start		uint64
	end		uint64
	memory 		uint64
	txCount 	uint64
	gas		uint64
	totalTxCount 	uint64
	totalGas	uint64
	lDisk		uint64
	aDisk		uint64
}


