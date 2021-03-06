package api

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/deepfabric/busybee/pkg/core"
	"github.com/deepfabric/busybee/pkg/pb/metapb"
	"github.com/deepfabric/busybee/pkg/storage"
	"github.com/fagongzi/log"
	"github.com/fagongzi/util/protoc"
	"github.com/fagongzi/util/task"
)

var (
	closeFlag = &struct{}{}
)

type tenantQueue struct {
	id      uint64
	running uint32
	eng     core.Engine
	cb      func(arg interface{}, value []byte, err error)

	runner         *task.Runner
	partitons      uint32
	partitonQueues []*task.Queue
	created        uint32
}

func newTenantQueue(id uint64, eng core.Engine, cb func(arg interface{}, value []byte, err error)) *tenantQueue {
	return &tenantQueue{
		id:     id,
		eng:    eng,
		cb:     cb,
		runner: task.NewRunner(),
	}
}

func (q *tenantQueue) waitCreated() {
	for {
		if atomic.LoadUint32(&q.created) == 1 {
			return
		}

		time.Sleep(time.Millisecond * 10)
	}
}

func (q *tenantQueue) add(ctx ctx) {
	q.waitCreated()
	q.partitonQueues[ctx.req.AddEvent.Event.UserID%q.partitons].Put(ctx)
}

func (q *tenantQueue) start() error {
	if atomic.CompareAndSwapUint32(&q.running, 0, 1) {
		value, err := q.eng.Storage().Get(storage.TenantMetadataKey(q.id))
		if err != nil {
			return err
		}

		if len(value) == 0 {
			return fmt.Errorf("tenant %d not init", q.id)
		}

		meta := &metapb.Tenant{}
		protoc.MustUnmarshal(meta, value)

		q.partitons = meta.Input.Partitions
		for i := uint32(0); i < q.partitons; i++ {
			queue := task.New(1024)
			q.partitonQueues = append(q.partitonQueues, queue)
			q.startPartition(i, queue)
		}

		atomic.StoreUint32(&q.created, 1)
	}

	return nil
}

func (q *tenantQueue) stop() {
	if atomic.CompareAndSwapUint32(&q.running, 1, 0) {
		for _, tq := range q.partitonQueues {
			tq.Put(closeFlag)
		}
	}
}

func (q *tenantQueue) startPartition(partition uint32, pq *task.Queue) {
	q.runner.RunCancelableTask(func(c context.Context) {
		items := make([]interface{}, 64, 64)
		var events [][]byte

		for {
			select {
			case <-c.Done():
				return
			default:
				n, err := pq.Get(64, items)
				if err != nil {
					log.Fatalf("BUG: queue must closed by self goroutine")
				}

				events = events[:0]
				for i := int64(0); i < n; i++ {
					item := items[i]
					if item == closeFlag {
						pq.Dispose()
						return
					}

					events = append(events, protoc.MustMarshal(&metapb.Event{
						Type: metapb.UserType,
						User: &item.(ctx).req.AddEvent.Event,
					}))
				}

				err = q.eng.Storage().PutToQueue(q.id, partition,
					metapb.TenantInputGroup, events...)
				for i := int64(0); i < n; i++ {
					q.cb(items[i], nil, err)
				}
			}
		}
	})
}
