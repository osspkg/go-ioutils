package cache

import (
	"context"
	"testing"
	"time"

	"go.osspkg.com/casecheck"
)

func TestUnit_WithTTL_Pointer(t *testing.T) {
	type A struct {
		Data uint64
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Second)
	c := NewWithTTL[string, *A](ctx, time.Minute)
	v, ok := c.Get("a")
	casecheck.False(t, ok)
	casecheck.Nil(t, v)
}
