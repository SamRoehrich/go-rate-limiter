package limiter

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type limiter struct {
	mu          sync.Mutex
	maxCapacity int
	capacity    int
	rate        int
	lastFilled  time.Time
}

func (h *limiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	diff := now.Sub(h.lastFilled)

	if int(diff.Seconds()) > h.rate {
		h.capacity += h.maxCapacity - h.capacity
		h.lastFilled = now
	}

	if h.capacity < 1 {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	h.capacity--

	fmt.Fprintf(w, "capacity is %d\n", h.capacity)
	log.Printf("capacity is %d\n", h.capacity)
}

func New() *limiter {
	return &limiter{
		maxCapacity: 5,
		capacity:    5,
		rate:        10,
		lastFilled:  time.Now(),
	}
}
