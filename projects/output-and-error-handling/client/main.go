package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	// if an error occurs during the request, it is handled by the handleError function
	if err := makeRequestWithRetry("http://localhost:8080"); err != nil {
		handleError(err)
		// The program exits with a status code of 1.
		os.Exit(1)
	}
}

func handleError(err error) {
	//that prints the error message to the console
	fmt.Println("Error:", err)
}

// function makes an HTTP GET request to the specified URL
func makeRequestWithRetry(url string) error {
	// it attempts to make the request up to 3 times
	retryAttempts := 3
	for attempt := 1; attempt <= retryAttempts; attempt++ {
		//If an error occurs during the request, it immediately returns the error.
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		//if the response status code is 429 (indicating that the server is overloaded with requests.)
		if resp.StatusCode == http.StatusTooManyRequests {
			//extracts the value of the Retry-After header from the HTTP response, which indicates how long the client should wait before making another request.
			retryAfterHeader := resp.Header.Get("Retry-After")
			if retryAfterHeader != "" {
				//convert the Retry-After header value to an integer representing the number of seconds to wait before retrying the reques
				retryAfterSeconds, err := strconv.Atoi(retryAfterHeader)
				if err != nil {
					// Retry-After header failed or is not an integer, use default duration
					defaultRetryDuration := 30 // seconds
					fmt.Printf("Server is too busy. Waiting for %d seconds before retrying (attempt %d/%d)...\n", defaultRetryDuration, attempt, retryAttempts)
					//causes the program to sleep for the specified number of seconds
					time.Sleep(time.Duration(defaultRetryDuration) * time.Second)
					continue
				}

				fmt.Printf("Server is too busy. Waiting for %d seconds before retrying (attempt %d/%d)...\n", retryAfterSeconds, attempt, retryAttempts)
				time.Sleep(time.Duration(retryAfterSeconds) * time.Second)
				continue
			}

			return fmt.Errorf("server is too busy, but no Retry-After header provided")
		}
		// closing the response body until the end of the function and prints the response status code
		defer resp.Body.Close()

		fmt.Println("Response status:", resp.Status)

		//reads and prints the first 9 lines of the response body
		scanner := bufio.NewScanner(resp.Body)
		for i := 0; scanner.Scan() && i < 9; i++ {
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		//if the response is successfully processed without needing a retry, it returns nil to indicate success.
		return nil
	}
	//If the maximum number of retry attempts is reached without a successful response, it returns an error indicating that the maximum retry attempts have been reached.
	return fmt.Errorf("reached maximum retry attempts")
}
