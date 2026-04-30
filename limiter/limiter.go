package limiter

import (
	"net/http"
	"strconv"
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

func (ls *Limiters) get(user string, m int, c int, r int) *Limiter {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, ok := ls.byUser[user]

	if !ok {
		l = New(m, c, r)
		ls.byUser[user] = l
	}

	return l
}

func (ls *Limiters) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	m, err := strconv.Atoi(r.URL.Query().Get("maxCapacity"))
	c, err := strconv.Atoi(r.URL.Query().Get("capacity"))
	rate, err := strconv.Atoi(r.URL.Query().Get("rate"))

	if len(user) == 0 {
		http.Error(w, "Missing user param", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Missing param", http.StatusBadRequest)
		return
	}

	l := ls.get(user, m, c, rate)
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
		diff = 0
	}

	if l.capacity < 1 {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	l.capacity--
	w.Header().Add("X-RateLimit-Remaining", strconv.Itoa(l.capacity))
	w.Header().Add("X-RateLimit-Limit", strconv.Itoa(l.maxCapacity))
	w.Header().Add("X-RateLimit-Reset", diff.String())
}

func New(m int, c int, r int) *Limiter {
	return &Limiter{
		maxCapacity: m,
		capacity:    c,
		rate:        r,
		lastFilled:  time.Now(),
	}
}

func Init() *Limiters {
	return &Limiters{
		byUser: make(map[string]*Limiter),
	}
}
