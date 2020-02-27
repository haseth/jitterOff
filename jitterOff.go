package jitteroff

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

var (
	errFailed error = errors.New(("Failed call"))
)

// JitterOff ...
type JitterOff struct {
	// Max retries
	attempt    int
	maxAttempt int

	// BackOff and jitter
	minTime     time.Duration
	capTime     time.Duration
	backOffTime time.Duration
}

// NewJitterOff creates a default
func NewJitterOff() *JitterOff {
	// Currently only supporting default values
	return &JitterOff{
		maxAttempt: 2,
		minTime:    100 * time.Millisecond,
		capTime:    400 * time.Millisecond,
	}
}

// Execute the requested function
func (j *JitterOff) Execute(request func() (interface{}, error)) (interface{}, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	for {
		out, err := request()
		if err == nil {
			return out, nil
		}

		j.attempt += 1
		// Or you could ask user to handle this
		if j.attempt >= j.maxAttempt {
			return out, err
		}

		// Back-off with jitter
		timeDuration := minimum(j.capTime, time.Duration((float64(j.minTime) * math.Pow(2.0, float64(j.attempt)))))
		j.backOffTime = timeDuration/2 + time.Duration(rand.Float64()*float64(timeDuration/2))

		time.Sleep(j.backOffTime)
	}

	return nil, nil
}

func minimum(a, b time.Duration) time.Duration {
	if a < b {
		return a
	} else {
		return b
	}
}
