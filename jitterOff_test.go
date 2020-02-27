package jitteroff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewJitter(t *testing.T) {
	/*
		Testing jitteroff with default timeouts
	*/
	// SETUP
	j := NewJitterOff()

	// VALIDATE
	assert.Equal(t, j.minTime, time.Duration(100*time.Millisecond), "Default timeout for wrapper set correctly")
	assert.Equal(t, j.capTime, time.Duration(400*time.Millisecond), "Default timeout for wrapper set correctly")
	assert.Equal(t, j.backOffTime, time.Duration(0*time.Millisecond), "Default timeout for wrapper set correctly")
	assert.Equal(t, j.maxAttempt, 2, "Max retries reached")
}

func TestExecute(t *testing.T) {
	/*
		Test-1
		with wait call so that call is successfull
	*/
	// SETUP
	j := NewJitterOff()

	// TEST
	_, err := j.Execute(doWaitCall(95 * time.Millisecond))

	// VALIDATE that request is successfull
	assert.Nil(t, err, "Function able to perform the task timely")
	assert.Equal(t, j.attempt, 0, "Max retries reached")

	/*
		Test-2
		with success call
	*/
	// TEST
	_, err = j.Execute(doSuccessCall())

	// VALIDATE that request is successfull
	assert.Nil(t, err, "Function able to perform the task timely")
	assert.Equal(t, j.attempt, 0, "Max retries reached")

	/*
		Test-3
		Request which fails for all max attempts.
	*/

	// TEST
	_, err = j.Execute(doFailCall())

	// VALIDATE
	assert.NotNil(t, err, "Function able to perform the task timely")
	assert.Equal(t, j.attempt, j.maxAttempt, "Max retries reached")
	assert.GreaterOrEqual(t, int(j.backOffTime), int(j.minTime), "backoff time reached cap time")
	assert.LessOrEqual(t, int(j.backOffTime), int(200*time.Millisecond), "backoff time reached cap time")
}

// do wait call
func doWaitCall(duration time.Duration) func() (interface{}, error) {
	return func() (interface{}, error) {
		// do some wait time task
		time.Sleep(duration)
		return nil, nil
	}
}

// do success call
func doSuccessCall() func() (interface{}, error) {
	return func() (interface{}, error) {
		return nil, nil
	}
}

// do fail call
func doFailCall() func() (interface{}, error) {
	return func() (interface{}, error) {
		return nil, errFailed
	}
}

func doPartialFailCall() func() (interface{}, error) {
	return func() (interface{}, error) {
		return nil, errFailed
	}
}
