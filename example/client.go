package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	// jitteroff "github.com/haseth/jitterOff"
	"../../jitteroff"
)

// Get has a middleware
func Get(url string) ([]byte, error) {
	joff := jitteroff.NewDefaultJitterOff()
	body, err := joff.Do(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})
	if err != nil {
		return nil, err
	}

	return body.([]byte), nil
}

func main() {
	b, err := Get("http://google.com")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
