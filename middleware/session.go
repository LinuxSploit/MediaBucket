// example validation util func
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ValidateData struct {
	SessionToken string `json:"sessionKey"`
}

type ValidateResponse struct {
	IsValid  bool   `json:"isValid"`
	UserType string `json:"userType"`
}

// ValidateSessionAndPerm checks if the session is valid and returns the user's permission level
// Return values: 0 = Invalid session, 1 = Regular user, 2 = Influencer
func ValidateSessionAndPerm(sessionToken string) int {

	const influencerType = "Influencer"
	const validSessionCode = 1
	const influencerSessionCode = 2
	const invalidSessionCode = 0

	body := ValidateData{
		SessionToken: sessionToken,
	}

	bodyWriter := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(bodyWriter).Encode(body); err != nil {
		log.Printf("Error encoding request body: %v", err)
		return invalidSessionCode
	}

	req, err := http.NewRequest("POST", os.Getenv("ValidateSessionURL"), bodyWriter)
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return invalidSessionCode
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return invalidSessionCode
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-200 response: %v", resp.Status)
		return invalidSessionCode
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return invalidSessionCode
	}

	var respData ValidateResponse
	if err := json.Unmarshal(data, &respData); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return invalidSessionCode
	}

	if !respData.IsValid {
		log.Println("Invalid session")
		return invalidSessionCode
	}

	if respData.UserType == influencerType {
		return influencerSessionCode
	}

	return validSessionCode
}
