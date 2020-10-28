package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateParallel(t *testing.T) {
	parallel := CreateParallel(5)

	start := time.Now()
	// should take about 5 seconds
	for i:=0; i<5; i++ {
		c := i
		parallel.Run(func() (interface{}, error) {
			time.Sleep(time.Duration(i) * time.Second)
			return c, nil
		})
	}

	results := parallel.Sync()
	end := time.Now()

	assert.WithinDuration(t, end, start, time.Duration(6) * time.Second)
	for i, r := range results {
		assert.Equal(t, i, r.Result)
		assert.Nil(t, r.Error)
	}

	assert.False(t, parallel.HasErrors())
}