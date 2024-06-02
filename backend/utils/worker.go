package utils

import (
	"fmt"
	"context"
	"encoding/gob"
	"bytes"
	"sync"
	"time"
	"errors"
	//"sync/atomic"
	"strings"
	"runtime/debug"

	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"
)


var (
	WorkerAddressFlag = cli.StringFlag {
		Name: "address",
		Usage: "address of Redis Stream to publish to",
		Aliases: []string{"a", "addr"},
		Value: "127.0.0.1",
	}
	WorkerPortFlag = cli.IntFlag {
		Name: "port",
		Usage: "port of Redis Stream to publish to",
		Aliases: []string{"p"},
		Value: 6379,
	}
	WorkerTopicFlag = cli.StringFlag {
		Name: "topic",
		Usage: "topic to send message to",
		Aliases: []string{"t", "c"},
		Value: "task-queue2",
	}
	WorkerTypFlag = cli.StringFlag {
		Name: "type",
		Usage: "type of task",
		Aliases: []string{"tt"},
	}
	WorkerMasterFlag = cli.StringFlag {
		Name: "master",
		Usage: "master id of task",
		Aliases: []string{"m", "mid"},
	}
	WorkerRunFlag = cli.StringFlag {
		Name: "run",
		Usage: "run id of task",
		Aliases: []string{"r", "rid"},
	}


)

func WorkerPeek(ctx *cli.Context) error {
	return workerPeek(
		ctx.String("address"),
		ctx.Int("port"),
		ctx.String("topic"),
	)
}

func workerPeek(addr string, port int, topic string) error {
	worker := NewWorker(addr, port, topic)
	worker.Peek(context.TODO())
	return nil 
}

func WorkerPush(ctx *cli.Context) error {
	return workerPush( 
		ctx.String("address"),
		ctx.Int("port"),
		ctx.String("topic"),
		ctx.String("type"),
		ctx.String("master"),
		ctx.String("run"),
	)
}

func workerPush(addr string, port int, topic string, typ string, master string, run string) error {
	worker := NewWorker(addr, port, topic)

	worker.Push(context.TODO(), &Task{
		Type: TaskType(typ),
		MasterId: master,
		RunId: run,
	})

	return nil 
}

func WorkerPop(ctx *cli.Context) error {
	return workerPop(
		ctx.String("address"),
		ctx.Int("port"),
		ctx.String("topic"),
	)
}

func workerPop(addr string, port int, topic string) error {
	worker := NewWorker(addr, port, topic)
	worker.PopLoop(context.TODO())
	return nil 
}

type Worker struct {
	r *redis.Client
	topic string

	loop sync.WaitGroup
}

func NewWorker(address string, port int, topic string) *Worker {
	r := redis.NewClient(&redis.Options {
		Addr: fmt.Sprintf("%s:%d", address, port),
	})

	return &Worker {
		r: r,
		topic: topic,
	}
}

func (w *Worker) Push(ctx context.Context, task *Task) error {
	var err error
	if task == nil {
		return fmt.Errorf("task is nil")
	}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err = enc.Encode(&task)
	if err != nil {
		fmt.Println("no task", err)
		return fmt.Errorf("unable to encode task: %w", err)
	}

	cmd := w.r.LPush(ctx, w.topic, buffer.Bytes())
	if cmd.Err() != nil {
		return fmt.Errorf("unable to process results: %w", err)
	}

	return nil
}

func (w *Worker) Peek(ctx context.Context) error {
	vals, err := w.r.LRange(ctx, w.topic, 0, -1).Result()

	if err != nil {
		fmt.Println("Peeking! - error")
		return err
	}

	fmt.Println("There is", len(vals), "messages in topic", w.topic)

	for ix, val := range vals {
		fmt.Println(ix, val)
	}

	return nil
}

func (w *Worker) PopLoop(ctx context.Context) {
	w.popLoop(ctx, 1, 1)
}

func (w *Worker) popLoop(ctx context.Context, i int, total int) {
	debug.SetPanicOnFault(true)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.pop(ctx, 5 * time.Second)
		}
	}
}

func (w *Worker) pop(ctx context.Context, timeout time.Duration) error {
	res, err := w.r.BRPop(ctx, timeout, w.topic).Result()

	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("unable to BRPOP: %w", err)
	}

	var task Task
	dec := gob.NewDecoder(strings.NewReader(res[1]))
	err = dec.Decode(&task)
	if err != nil {
		return fmt.Errorf("unable to decode task: %w", err)
	}

	w.handle(ctx, &task)

	return nil
}

func (w *Worker) handle(ctx context.Context, task *Task) {
	if task == nil {
		return
	}

	var handler TaskHandler
	if h, exists := TaskHandlers[task.Type]; exists {
		handler = h
	} else {
		handler = TaskHandlers["default"]
	}

	taskCtx, _ := context.WithCancel(ctx)
	handler(taskCtx, task)	
}

