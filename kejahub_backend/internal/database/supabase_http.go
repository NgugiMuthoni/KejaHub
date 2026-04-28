package database

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func getSupabaseConfig() (string, string) {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")
	return url + "/rest/v1", key
}

func Get(table string, query string) ([]map[string]interface{}, error) {

	baseURL, key := getSupabaseConfig()

	fullURL := baseURL + "/" + table + query

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	return result, err
}

func Insert(table string, data map[string]interface{}) error {

	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", baseURL+"/"+table, bytes.NewBuffer(jsonData))

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}