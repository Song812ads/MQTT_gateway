package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

func getRequest(c http.Client, url string) bool {

	resp, err := c.Get(url)
	var send = false

	if err != nil {
		// Handle the error (e.g., network issue, service down)
		fmt.Println("Waiting for service to begin...")
	} else {
		// Check if the response status code is 200
		if resp.StatusCode == 200 {
			send = true
			// fmt.Println("Send success")
			resp.Body.Close() // Close the response body when done
		} else {
			// If status code is not 200, handle the error or retry
			fmt.Printf("Received status %d. Retrying...\n", resp.StatusCode)
			resp.Body.Close() // Close the response body even on failure
		}
	}
	return send
}

func postRequest(c http.Client, url string, data []byte) error {
	reader := bytes.NewReader(data)

	resp, err := c.Post(url, "application/json", reader)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 { // 200 OK
		// WriteError(w, resp.StatusCode, fmt.Errorf("Received non-OK status code: %d", resp.StatusCode))

		// fmt.Println(resp.StatusCode)
		return fmt.Errorf("Error adding device")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error reading response body"))
		// fmt.Println(err)
		return err
	}

	// Parse the JSON response body into a slice of maps (assuming the response is an array of objects)
	var responseData []map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		// WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error parsing response body"))
		// fmt.Println(err)
		return err
	}

	for _, item := range responseData {
		// Type assertion to extract the statusCode
		if statusCode, ok := item["statusCode"].(float64); ok { // Use float64 because JSON numbers are usually parsed as float64
			if math.Round(statusCode/100) != 2 {
				// WriteError(w, http.StatusBadRequest, fmt.Errorf("Error adding profile device"))
				// fmt.Println(statusCode)
				return fmt.Errorf("Error adding device")
			}
			fmt.Printf("StatusCode: %.0f\n", statusCode)
		} else {
			// WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error parsing response body"))
			// fmt.Println("StatusCode not found in map")
			return fmt.Errorf("Error parsing response body")
		}
	}

	return nil

}
