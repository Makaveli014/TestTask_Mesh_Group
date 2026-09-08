package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIConfig struct {
	URI              string
	AuthLoginPwd     string
	UserAgent        string
	TimeoutSeconds   int
	ImportBatchSize  int
}

// FetchSAPData fetches a batch of segmentation data from the SAP API.
func FetchSAPData(client *http.Client, cfg APIConfig, offset int) ([]Segmentation, error) {
	url := fmt.Sprintf("%s?p_limit=%d&p_offset=%d", cfg.URI, cfg.ImportBatchSize, offset)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(cfg.AuthLoginPwd))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("User-Agent", cfg.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	var result []Segmentation
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return result, nil
}
