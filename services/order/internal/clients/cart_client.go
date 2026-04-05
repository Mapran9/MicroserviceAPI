package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"order/internal/models"
)

func GetCart(cartID string) (*models.CartDTO, error) {
	baseURL := strings.TrimRight(os.Getenv("CART_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = "http://localhost:8003"
	}

	url := fmt.Sprintf("%s/api/Carts/%s", baseURL, cartID)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
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

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
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
