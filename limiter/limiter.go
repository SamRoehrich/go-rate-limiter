package limiter

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	mu          sync.Mutex
	maxCapacity int
	capacity    int
	rate        int
	lastFilled  time.Time
}

type Limiters struct {
	mu     sync.Mutex
	byUser map[string]*Limiter
}

func (ls *Limiters) get(user string) *Limiter {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, ok := ls.byUser[user]

	if !ok {
		l = New()
		ls.byUser[user] = l
	}

	return l
}

func (ls *Limiters) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	if len(user) == 0 {
		http.Error(w, "Missing user param", http.StatusBadRequest)
		return
	}

	l := ls.get(user)
	l.ServeHTTP(w, r)
}

func (l *Limiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	diff := now.Sub(l.lastFilled)

	if int(diff.Seconds()) > l.rate {
		l.capacity += l.maxCapacity - l.capacity
		l.lastFilled = now
	}

	if l.capacity < 1 {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	l.capacity--
	fmt.Fprintf(w, "capacity is %d\n", l.capacity)
	log.Printf("capacity is %d\n", l.capacity)
}

func New() *Limiter {
	return &Limiter{
		maxCapacity: 5,
		capacity:    5,
		rate:        10,
		lastFilled:  time.Now(),
	}
}

func Init() *Limiters {
	return &Limiters{
		byUser: make(map[string]*Limiter),
	}
}
