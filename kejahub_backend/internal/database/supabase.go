package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var baseURL = os.Getenv("SUPABASE_URL") + "/rest/v1"
var apiKey = os.Getenv("SUPABASE_API_KEY")

// INSERT


// GET ONE
func GetOne(table string, query string) (map[string]interface{}, error) {

	url := baseURL + "/" + table + query

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result []map[string]interface{}
	json.Unmarshal(body, &result)

	if len(result) == 0 {
		return nil, fmt.Errorf("not found")
	}

	return result[0], nil
}

// GET ALL
func GetAll(table string, query string) ([]map[string]interface{}, error) {

	url := baseURL + "/" + table + query

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result []map[string]interface{}
	json.Unmarshal(body, &result)

	return result, nil
}

// UPDATE
func Update(table string, id string, data map[string]interface{}) error {

	jsonData, _ := json.Marshal(data)

	url := baseURL + "/" + table + "?id=eq." + id

	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}