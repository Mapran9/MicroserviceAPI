package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type ProductResponse struct {
	ProductID string   `json:"product_id"`
	Price     *float64 `json:"price"`
}

var productHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	},
}

func GetProductPrice(productID string) (float64, error) {
	baseURL := strings.TrimRight(os.Getenv("PRODUCT_BASE_URL"), "/")
	if baseURL == "" {
		// สำหรับรันนอก docker เท่านั้น
		baseURL = "http://localhost:8002"
	}

	url := fmt.Sprintf("%s/api/Products/%s", baseURL, productID)

	resp, err := productHTTPClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("call product-service failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return 0, fmt.Errorf("product not found (404): %s", productID)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("product-service returned %d: %s", resp.StatusCode, string(body))
	}

	var p ProductResponse
	if err := json.Unmarshal(body, &p); err != nil {
		return 0, fmt.Errorf("decode product response failed: %w (body=%s)", err, string(body))
	}

	if p.Price == nil {
		return 0, fmt.Errorf("product price is null: %s", productID)
	}
	return *p.Price, nil
}
