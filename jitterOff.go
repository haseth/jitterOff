/*
jitteroff package implements backOff with jitter algorithm.

During high load, multiple requests would fight for the same
resources. To avoid high contention and bombarding a particular
service with all the request at the same time we can organize the
requests over period of time using back-off with jitter algorithm.

BackOff Algorithm:
- If request fails, we will retry the requests again after
back-off period of time rather than immediately.
- This would prevent bombarding all the request to a service again.

Jitter Algorithm:
- Making all the request again at same back-off time would again cause
high contention.
- We add jitter to back-off time to create a sense of randominess while
request is retried.

Package cap the number of retries by taking "maxAttempt" input from user.
It also caps the maximum back-off time period by "capTime" input provided
by user.

- BackOffTimeCalculation:
basetime* 2**attempt

if baseTime is 1s
attempt backOffTime (in sec.)
   1       0
   2       2
   3       4
   4       8
   .       .
   .       .
maxAttempt capTime

- Adding Random Jitter
baseTime/2 + [0,1]*baseTime/2

if baseTime is 1s
attempt backOffTime (in sec.)
   1       [0]
   2       [1,2]
   3       [2,4]
   4       [4,8]
   .       .
   .       .
maxAttempt [capTime/2+rand()*capTime/2]
*/

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

// JitterOff structures retry and time values required for
// backOff and Jitter algorithm.
type JitterOff struct {
	attempt    int
	maxAttempt int

	minTime     time.Duration
	capTime     time.Duration
	backOffTime time.Duration
}

// NewDefaultJitterOff creates a JitterOff structure with
// default values.
func NewDefaultJitterOff() *JitterOff {
	return &JitterOff{
		maxAttempt: 2,
		minTime:    100 * time.Millisecond,
		capTime:    400 * time.Millisecond,
	}
}

// NewCustomJitterOff creates a JitterOff structure with
// custom/user-defined values.
func NewCustomJitterOff(totalAttempt int, time1, time2 time.Duration) *JitterOff {
	return &JitterOff{
		maxAttempt: totalAttempt,
		minTime:    time1,
		capTime:    time2,
	}
}

// Do the requested function with BackOff and Jitter Algorithm
func (j *JitterOff) Do(request func() (interface{}, error)) (interface{}, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	for {
		out, err := request()
		if err == nil {
			return out, nil
		}

		j.attempt++
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
	// min of time.Duration
	if a < b {
		return a
	}
	return b
}
