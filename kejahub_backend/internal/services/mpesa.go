package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func STKPush(phone string, amount int) error {

	url := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"

	payload := map[string]interface{}{
		"BusinessShortCode": "YOUR_SHORTCODE",
		"Password":          "GENERATED_PASSWORD",
		"Timestamp":         "TIMESTAMP",
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            amount,
		"PartyA":            phone,
		"PartyB":            "YOUR_SHORTCODE",
		"PhoneNumber":       phone,
		"CallBackURL":       "https://yourdomain.com/mpesa/callback",
		"AccountReference":  "Rent Payment",
		"TransactionDesc":   "Rent",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer ACCESS_TOKEN")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)

	return err
}
