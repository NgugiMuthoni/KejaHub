package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type STKResponse struct {
	MerchantRequestID string `json:"MerchantRequestID"`
	CheckoutRequestID string `json:"CheckoutRequestID"`
	ResponseCode      string `json:"ResponseCode"`
	CustomerMessage   string `json:"CustomerMessage"`
}

// 🔐 GET ACCESS TOKEN
func getAccessToken() (string, error) {

	fmt.Println("REQUESTING ACCESS TOKEN...")

	key := strings.TrimSpace(os.Getenv("MPESA_CONSUMER_KEY"))
	secret := strings.TrimSpace(os.Getenv("MPESA_CONSUMER_SECRET"))

	url := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(key, secret)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("TOKEN REQUEST ERROR:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println("TOKEN RESPONSE RECEIVED:")
	fmt.Println(string(body))

	var result map[string]interface{}
	_ = json.Unmarshal(body, &result)

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get access token")
	}

	return token, nil
}

// 🚀 STK PUSH
func STKPush(phone string, amount int) (STKResponse, error) {

	token, err := getAccessToken()
	if err != nil {
		return STKResponse{}, err
	}

	fmt.Println("TOKEN ACQUIRED:", token)
	fmt.Println("PREPARING STK REQUEST...")

	shortcode := strings.TrimSpace(os.Getenv("MPESA_SHORTCODE"))
	passkey := strings.TrimSpace(os.Getenv("MPESA_PASSKEY"))
	callback := strings.TrimSpace(os.Getenv("MPESA_CALLBACK_URL"))

	timestamp := time.Now().Format("20060102150405")

	// 🔥 RAW STRING (THIS IS WHAT DARARAJA CARES ABOUT)
	raw := shortcode + passkey + timestamp

	password := base64.StdEncoding.EncodeToString([]byte(raw))

	// 🧠 DEBUG BLOCK (CRITICAL)
	fmt.Println("=== STK DEBUG ===")
	fmt.Println("SHORTCODE:", "["+shortcode+"]")
	fmt.Println("PASSKEY:", "["+passkey+"]")
	fmt.Println("TIMESTAMP:", "["+timestamp+"]")
	fmt.Println("RAW STRING:", "["+raw+"]")
	fmt.Println("PASSWORD:", password)

	payload := map[string]interface{}{
		"BusinessShortCode": shortcode,
		"Password":          password,
		"Timestamp":         timestamp,
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            amount,
		"PartyA":            phone,
		"PartyB":            shortcode,
		"PhoneNumber":       phone,
		"CallBackURL":       callback,
		"AccountReference":  "KejaHub Rent",
		"TransactionDesc":   "Rent Payment",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return STKResponse{}, err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var resp *http.Response

	// 🔁 RETRY (FIXED: recreate request each time)
	for i := 0; i < 2; i++ {

		fmt.Println("SENDING STK REQUEST... attempt:", i+1)

		req, err := http.NewRequest(
			"POST",
			"https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest",
			bytes.NewBuffer(body),
		)
		if err != nil {
			return STKResponse{}, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)

		if err == nil {
			fmt.Println("RESPONSE RECEIVED")
			break
		}

		fmt.Println("REQUEST ERROR:", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return STKResponse{}, err
	}

	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	fmt.Println("MPESA RAW RESPONSE:")
	fmt.Println(string(respBody))

	var result STKResponse
	_ = json.Unmarshal(respBody, &result)

	return result, nil
}
