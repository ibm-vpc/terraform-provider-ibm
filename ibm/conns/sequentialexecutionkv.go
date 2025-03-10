// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package conns

import (
	"log"
	"sync"
)

var IbmSequentialExecutionKV = NewSequentialExecutionKV()

// SequentialExecutionKV ensures that executions happen sequentially but do not wait for completion.
type SequentialExecutionKV struct {
	mu    sync.Mutex
	store map[string]*sync.Cond
}

// NewSequentialExecutionKV initializes the store.
func NewSequentialExecutionKV() *SequentialExecutionKV {
	return &SequentialExecutionKV{
		store: make(map[string]*sync.Cond),
	}
}

// LockSequential ensures only one execution enters at a time but immediately releases for the next.
func (m *SequentialExecutionKV) LockSequential(key string) {
	m.mu.Lock()

	// If no condition variable exists for the key, create one
	if _, exists := m.store[key]; !exists {
		m.store[key] = sync.NewCond(&sync.Mutex{})
	}
	cond := m.store[key]
	m.mu.Unlock()

	// Ensure sequential access
	cond.L.Lock()
	log.Printf("[DEBUG] Entering execution for %q", key)
	cond.Broadcast() // Notify the next one to proceed
	cond.L.Unlock()
}

// UnlockSequential does nothing in this case, but it's here for consistency.
func (m *SequentialExecutionKV) UnlockSequential(key string) {
	log.Printf("[DEBUG] Execution completed for %q", key)
}
