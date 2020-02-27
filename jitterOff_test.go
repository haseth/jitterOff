package jitteroff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewJitter(t *testing.T) {
	/*
		Testing jitteroff with default settings
	*/
	// SETUP
	j := NewDefaultJitterOff()

	// VALIDATE
	assert.Equal(t, j.minTime, time.Duration(100*time.Millisecond), "Default min timeout correctly set.")
	assert.Equal(t, j.capTime, time.Duration(400*time.Millisecond), "Default cap timeout correctly set.")
	assert.Equal(t, j.backOffTime, time.Duration(0*time.Millisecond), "Default backoff timeout correctly set.")
	assert.Equal(t, j.maxAttempt, 2, "Default max attempt correctly set.")
}

func TestCustomerJitter(t *testing.T) {
	/*
		Testing jitteroff with custom settings
	*/
	totalAttempt := 2
	time1 := 100 * time.Millisecond
	time2 := 400 * time.Millisecond

	// SETUP
	j := NewCustomJitterOff(totalAttempt, time1, time2)

	// VALIDATE
	assert.Equal(t, j.minTime, time1, "Custom min timeout correctly set.")
	assert.Equal(t, j.capTime, time2, "Custom max timeout correctly set.")
	assert.Equal(t, j.backOffTime, time.Duration(0*time.Second), "Custom backoff timeout correctly set.")
	assert.Equal(t, j.maxAttempt, totalAttempt, "Custom max attempt correctly set.")
}

func TestExecute_DefaultSetting(t *testing.T) {
	/*
		Test-1
		with wait call so that call is successfull
	*/
	// SETUP
	j := NewDefaultJitterOff()

	// TEST
	_, err := j.Do(doWaitCall(95 * time.Millisecond))

	// VALIDATE that request is successfull
	assert.Nil(t, err, "Function able to perform the task timely")
	assert.Equal(t, j.attempt, 0, "No retry required for successful call")

	/*
		Test-2
		with success call
	*/
	// TEST
	_, err = j.Do(doSuccessCall())

	// VALIDATE that request is successfull
	assert.Nil(t, err, "Function able to perform the task timely")
	assert.Equal(t, j.attempt, 0, "No retry required for successful call")

	/*
		Test-3
		Request which fails for all max attempts.
	*/

	// TEST
	_, err = j.Do(doFailCall())

	// VALIDATE
	assert.NotNil(t, err, "Function able to perform the task timely")
	assert.Equal(t, j.attempt, j.maxAttempt, "Max retries reached")
	assert.Greater(t, int(j.backOffTime), int(j.minTime), "No retry required for successful call")
	assert.LessOrEqual(t, int(j.backOffTime), int(200*time.Millisecond), "No retry required for successful call")
}

func TestExecute_CustomSetting(t *testing.T) {
	/*
		Testing jitteroff execute method with custom settings
	*/
	totalAttempt := 5
	time1 := 200 * time.Millisecond
	time2 := 800 * time.Millisecond

	// SETUP
	j := NewCustomJitterOff(totalAttempt, time1, time2)

	// TEST
	_, err := j.Do(doWaitCall(95 * time.Millisecond))

	// VALIDATE that request is successfull
	assert.Nil(t, err, "Function able to perform the task timely")
	assert.Equal(t, j.attempt, 0, "No retry required for successful call")

	/*
		Test-2
		with success call
	*/
	// TEST
	_, err = j.Do(doSuccessCall())

	// VALIDATE that request is successfull
	assert.Nil(t, err, "Function able to perform the task timely")
	assert.Equal(t, j.attempt, 0, "No retry required for successful call")

	/*
		Test-3
		Request which fails for all max attempts.
	*/

	// TEST
	_, err = j.Do(doFailCall())

	// VALIDATE
	assert.NotNil(t, err, "Function able to perform the task timely")
	assert.Equal(t, j.attempt, totalAttempt, "Max retries reached")
	assert.GreaterOrEqual(t, int(j.backOffTime), int(time1), "No retry required for successful call")
	assert.LessOrEqual(t, int(j.backOffTime), int(time2), "No retry required for successful call")
}

// HELPER

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
