// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0
package conns

import (
	"log"
	"sync"
	"time"
)

var BatchIbmMutexKV = NewBatchExecutionKV(3, 120*time.Second)

// BatchExecutionKV controls execution in batches with a delay between them.
type BatchExecutionKV struct {
	mu          sync.Mutex
	cond        *sync.Cond
	activeCount int
	batchSize   int
	delay       time.Duration
	waiting     int
}

// NewBatchExecutionKV initializes a new batch executor.
func NewBatchExecutionKV(batchSize int, delay time.Duration) *BatchExecutionKV {
	b := &BatchExecutionKV{
		batchSize: batchSize,
		delay:     delay,
	}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// LockBatch ensures only a certain number of executions run at the same time.
func (b *BatchExecutionKV) LockBatch(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// If the batch is full, wait for an available slot
	for b.activeCount >= b.batchSize {
		b.waiting++
		log.Printf("[BatchTest] Waiting for execution slot for %s. Active: %d/%d, Waiting: %d", key, b.activeCount, b.batchSize, b.waiting)
		b.cond.Wait()
		b.waiting--
	}

	// Start execution
	b.activeCount++
	log.Printf("[BatchTest] Started execution for %s (Active: %d/%d, Waiting: %d)", key, b.activeCount, b.batchSize, b.waiting)
}

// UnlockBatch releases the execution slot.
func (b *BatchExecutionKV) UnlockBatch(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.activeCount--
	log.Printf("[BatchTest] Execution completed for %s. Remaining active: %d", key, b.activeCount)

	// If this was the last in the batch, delay before allowing the next batch
	if b.activeCount == 0 && b.waiting > 0 {
		log.Printf("[BatchTest] Batch execution complete. Waiting %v before next batch starts.", b.delay)
		time.Sleep(b.delay)
	}

	// Notify waiting executions
	b.cond.Broadcast()
}
