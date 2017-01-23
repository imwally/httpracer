package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Response struct {
	Error       error
	Location    string
	StatusCode  int
	TimeElapsed time.Duration
}

func fetcher(location string) *Response {
	u, err := url.ParseRequestURI(location)
	if err != nil {
		return &Response{err, location, 0, time.Duration(0)}
	}

	startTime := time.Now()
	resp, err := http.Get(u.String())
	if err != nil {
		return &Response{err, location, 0, time.Duration(0)}
	}
	defer resp.Body.Close()

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	response := &Response{
		nil,
		u.String(),
		resp.StatusCode,
		elapsedTime,
	}

	return response
}

func fetch(responses chan<- *Response, locations []string) {
	var wg sync.WaitGroup
	for _, location := range locations {
		wg.Add(1)
		go func(location string) {
			defer wg.Done()
			responses <- fetcher(location)
		}(location)
	}

	wg.Wait()
	close(responses)
}

func printer(r <-chan *Response) {
	for response := range r {
		if response.Error != nil {
			fmt.Printf("%s: %s\n",
				os.Args[0],
				response.Error,
			)
		} else {
			fmt.Printf("%s:\t%d\t%s\n",
				response.TimeElapsed.String(),
				response.StatusCode,
				response.Location,
			)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s: need at least one argument.\n", os.Args[0])
	}

	responses := make(chan *Response)

	go fetch(responses, os.Args[1:])
	printer(responses)
}
