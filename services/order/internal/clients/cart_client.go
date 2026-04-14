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

	"order/internal/models"
)

var cartHTTPClient = &http.Client{
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

func GetCart(cartID string) (*models.CartDTO, error) {
	baseURL := strings.TrimRight(os.Getenv("CART_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = "http://localhost:8003"
	}

	url := fmt.Sprintf("%s/api/Carts/%s", baseURL, cartID)

	resp, err := cartHTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("call cart-service failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("cart not found: %s", cartID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cart-service returned %d: %s", resp.StatusCode, string(body))
	}

	var cart models.CartDTO
	if err := json.Unmarshal(body, &cart); err != nil {
		return nil, fmt.Errorf("decode cart response failed: %w", err)
	}

	return &cart, nil
}

func UpdateCartStatus(cartID, status string) error {
	baseURL := strings.TrimRight(os.Getenv("CART_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = "http://localhost:8003"
	}

	url := fmt.Sprintf("%s/api/Carts/%s/status", baseURL, cartID)

	reqBody := fmt.Sprintf(`{"status":"%s"}`, status)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := cartHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("update cart status failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cart-service returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
