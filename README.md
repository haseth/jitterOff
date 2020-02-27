# JitterOff 

During high load, multiple requests would fight for the same
resources. To avoid high contention and bombarding a particular service with all the request at the same time we can organize the requests over period of time using back-off with jitter algorithm.

JitterOff package implements backOff with jitter algorithm.

BackOff Algorithm:
- If request fails, we will retry the requests again after
back-off period of time rather than immediately.
- This would prevent bombarding all the request to a service again.

Jitter Algorithm:
- Making all the request again at same back-off time would again cause high contention.
- We add jitter to back-off time to create a sense of randominess while request is retried.

# BackOff with Jitter Calculation 

```BackOffTimeCalculation```:
basetime* 2**attempt

**if baseTime is 1s**
| attempt | backOffTime (in sec.) | 
|--------|------|
|1|0|
|2|2|
|3|4|
|4|8|
|.|.|
|.|.|
|maxAttempt |capTime|

```Adding Random Jitter```
baseTime/2 + [0,1]*baseTime/2

if baseTime is 1s
|attempt |backOffTime (range in sec.)|
|-------|-------|
|1| [0]|
|2| [1,2]|
|3| [2,4]|
|4| [4,8]|
|.| .|
|.| .|
|maxAttempt |[capTime/2+rand()*capTime/2]|

# Example 

``` go 
    joff := jitteroff.NewDefaultJitterOff()
	body, err := joff.Do(func() (interface{}, error) {
        // do something 
        // request if failed will be retried with back-off and jitter algorithm
		return nil, nil
	})
```

License
-------
The MIT License (MIT)



# TODO 
1. Add more features 
