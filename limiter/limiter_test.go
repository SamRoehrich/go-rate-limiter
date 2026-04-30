package limiter

import (
	"testing"
)

func TestLimitersInit(t *testing.T) {
	ls := Init()

	if ls == nil {
		t.Errorf("Expected limiters struct to be created but it is nil")
	}

	if ls.byUser == nil {
		t.Errorf("Expected limiters.byUser to be created but it is nil")
	}

	if len(ls.byUser) != 0 {
		t.Errorf("Expected byUser length to be 0, got %d\n", len(ls.byUser))
	}
}

func TestLimitersCreateLimiter(t *testing.T) {
	ls := Init()
	l := ls.get("Sam", 10, 5, 1)

	if l == nil {
		t.Errorf("Expected limiter to be created but it was not")
	}

	if l.capacity != 5 {
		t.Errorf("Expected capacity to equal 5, got %d\n", l.capacity)
	}
	if l.maxCapacity != 10 {
		t.Errorf("Expected capacity to equal 10, got %d\n", l.maxCapacity)
	}
	if l.rate != 1 {
		t.Errorf("Expected capacity to equal 5, got %d\n", l.rate)
	}
}
