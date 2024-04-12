package main

import (
    "bufio"
    "fmt"
    "net/http"
	"os"
	"time"
	"strconv"
	"errors"
	"io"
)

func main() {
    resp, err := http.Get("http://localhost:8080")

   if err != nil {
		if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Println("The server closed the connection unexpectedly. Please try again later.")
		} else {
			fmt.Println("Error:", err)
		}
		os.Exit(1)
	}
	}
    defer resp.Body.Close()

    fmt.Println("Response status:", resp.Status)

    scanner := bufio.NewScanner(resp.Body)
    for i := 0; scanner.Scan() && i < 10; i++ {
        fmt.Println(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }
	if resp.StatusCode == http.StatusTooManyRequests {
        handleTooManyRequests(resp)
    }
}

func handleTooManyRequests(resp *http.Response) {
    retryAfterHeader := resp.Header.Get("Retry-After")
    if retryAfterHeader != "" {
        retryAfterSeconds, err := strconv.Atoi(retryAfterHeader)
        if err == nil {
            fmt.Printf("Server is too busy. Waiting for %d seconds before retrying...\n", retryAfterSeconds)
            time.Sleep(time.Duration(retryAfterSeconds) * time.Second)
            return
        }

        retryAfterTime, err := time.Parse(http.TimeFormat, retryAfterHeader)
        if err == nil {
            durationUntilRetry := retryAfterTime.Sub(time.Now())
            fmt.Printf("Server is too busy. Waiting until %s before retrying...\n", retryAfterHeader)
            time.Sleep(durationUntilRetry)
            return
        }

		if retryAfterHeader == "a while" {
            defaultRetryDuration := 10 
            fmt.Printf("Server is too busy. Waiting for %d seconds before retrying...\n", defaultRetryDuration)
            time.Sleep(time.Duration(defaultRetryDuration) * time.Second)
            return
        }

        fmt.Println("Error parsing Retry-After header:", err)
    }
}