package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	workerIDBits     = 5
	datacenterIDBits = 5
	sequenceBits     = 12

	maxWorkerID     = -1 ^ (-1 << workerIDBits)
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits)
	sequenceMask    = -1 ^ (-1 << sequenceBits)

	workerIDShift      = sequenceBits
	datacenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

type Generator struct {
	mu           sync.Mutex
	workerID     int64
	datacenterID int64
	sequence     int64
	lastStamp    int64
}

func New(workerID, datacenterID int64) (*Generator, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker id out of range")
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errors.New("datacenter id out of range")
	}
	return &Generator{
		workerID:     workerID,
		datacenterID: datacenterID,
	}, nil
}

func (g *Generator) NextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()
	if now < g.lastStamp {
		now = g.lastStamp
	}
	if now == g.lastStamp {
		g.sequence = (g.sequence + 1) & sequenceMask
		if g.sequence == 0 {
			for now <= g.lastStamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastStamp = now
	return ((now - epoch) << timestampLeftShift) | (g.datacenterID << datacenterIDShift) | (g.workerID << workerIDShift) | g.sequence
}
