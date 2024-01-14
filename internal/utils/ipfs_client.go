package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetMetadataFromIPFS retrieves JSON metadata from IPFS and stores it in the database.
func GetMetadataFromIPFS(ctx context.Context, cid string) ([]byte, error) {
	// Fetch data from IPFS
	resp, err := NewHTTPClient().Get("https://ipfs.io/ipfs/" + cid) // do not use default client

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch data from IPFS. Status: %s", resp.Status)
	}

	// Parse JSON data
	var jsonData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		return nil, err
	}

	// Convert the JSON data to a string
	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	return jsonString, nil
}
