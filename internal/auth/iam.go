package auth

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type IAMTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type IAMAuthServiceInterface interface {
	GetToken(ctx context.Context) (string, error)
	IsAuthenticated() bool
	ClearToken()
}

type IAMAuthService struct {
	keyID    string
	secret   string
	endpoint string
	client   *http.Client

	token     string
	expiresAt time.Time
	mutex     sync.RWMutex
}

func NewIAMAuthService(keyID, secret, endpoint string) *IAMAuthService {
	keyID = strings.TrimSpace(keyID)
	secret = strings.TrimSpace(secret)
	endpoint = strings.TrimSpace(strings.TrimRight(endpoint, "/"))
	if endpoint == "" {
		endpoint = "https://iam.api.cloud.ru"
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Отключает проверку сертификата
		},
	}
	return &IAMAuthService{
		keyID:    keyID,
		secret:   secret,
		endpoint: endpoint,
		client:   &http.Client{Timeout: 30 * time.Second, Transport: tr},
	}
}

func (s *IAMAuthService) GetToken(ctx context.Context) (string, error) {
	s.mutex.RLock()
	if s.token != "" && time.Now().Before(s.expiresAt) {
		token := s.token
		s.mutex.RUnlock()
		return token, nil
	}
	s.mutex.RUnlock()
	return s.refreshToken(ctx)
}

func (s *IAMAuthService) refreshToken(ctx context.Context) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.token != "" && time.Now().Before(s.expiresAt) {
		return s.token, nil
	}

	payload, err := marshalIAMGetTokenBody(s.keyID, s.secret)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.endpoint+"/api/v1/auth/token",
		bytes.NewReader(payload),
	)
	if err != nil {
		return "", fmt.Errorf("create IAM token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request IAM token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read IAM token response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IAM API error (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResponse IAMTokenResponse
	if err = json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("unmarshal IAM token response: %w", err)
	}
	if tokenResponse.AccessToken == "" {
		return "", fmt.Errorf("IAM API returned empty access token")
	}

	expSec := tokenResponse.ExpiresIn - 300
	if expSec <= 0 {
		expSec = 300
	}

	s.token = tokenResponse.AccessToken
	s.expiresAt = time.Now().Add(time.Duration(expSec) * time.Second)
	return s.token, nil
}

// marshalIAMGetTokenBody — тело POST /api/v1/auth/token. Разные площадки шлют proto/REST JSON
// по-разному: clientId vs client_id, clientSecret vs client_secret, либо keyId/secret. Дублируем
// совместимые ключи, чтобы ненулевой ClientId всегда распознался.
func marshalIAMGetTokenBody(keyID, secret string) ([]byte, error) {
	if keyID == "" {
		return nil, fmt.Errorf("IAM Key ID (client id) is empty: set IAM_KEY_ID or run auth login")
	}
	if secret == "" {
		return nil, fmt.Errorf("IAM secret is empty: set IAM_SECRET or run auth login")
	}
	type kv struct{ k, v string }
	entries := []kv{
		{"clientId", keyID},
		{"client_id", keyID},
		{"keyId", keyID},
		{"key_id", keyID},
		{"clientSecret", secret},
		{"client_secret", secret},
		{"secret", secret},
	}
	var b strings.Builder
	b.WriteByte('{')
	for i, e := range entries {
		if i > 0 {
			b.WriteByte(',')
		}
		key, err := json.Marshal(e.k)
		if err != nil {
			return nil, err
		}
		val, err := json.Marshal(e.v)
		if err != nil {
			return nil, err
		}
		b.Write(key)
		b.WriteByte(':')
		b.Write(val)
	}
	b.WriteByte('}')
	return []byte(b.String()), nil
}

func (s *IAMAuthService) IsAuthenticated() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.token != "" && time.Now().Before(s.expiresAt)
}

func (s *IAMAuthService) ClearToken() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.token = ""
	s.expiresAt = time.Time{}
}
