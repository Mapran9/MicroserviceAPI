package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type CreatePaymentRequest struct {
	CustomerID    string  `json:"customer_id"`
	OrderID       string  `json:"order_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
}

type CreatePaymentResponse struct {
	Message   string  `json:"message"`
	PaymentID string  `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
}

func CreatePayment(req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	baseURL := strings.TrimRight(os.Getenv("PAYMENT_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = "http://localhost:8005"
	}

	url := fmt.Sprintf("%s/api/payments/internal", baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call payment-service failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("payment-service returned %d: %s", resp.StatusCode, string(respBody))
	}

	var out CreatePaymentResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("decode payment response failed: %w", err)
	}

	return &out, nil
}
