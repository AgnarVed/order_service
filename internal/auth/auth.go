package auth

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tests2/internal/config"
	"tests2/internal/models"
	"time"
)

type Auth struct {
	gatewayURL string
	httpClient *http.Client
}

type validateTokenResponse struct {
	UID       string                 `json:"uid"`
	IsActive  bool                   `json:"isActive"`
	User      models.UserFromRequest `json:"user"`
	IssuedAt  int64                  `json:"issuedAt"`
	ExpiresAt int64                  `json:"expiresAt"`
	TimeLeft  int64                  `json:"timeLeft"`
}

func NewAuth(cfg *config.Config) *Auth {
	return &Auth{
		gatewayURL: cfg.GatewayURL,
		httpClient: &http.Client{
			Timeout: time.Minute,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

// GatewayValidateToken отправляет запрос к Gateway-сервису для валидации токена
func (m *Auth) GatewayValidateToken(ctx context.Context, accessToken string) (*validateTokenResponse, error) {
	values := map[string]string{"accessToken": accessToken}
	json_data, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/validate", m.gatewayURL),
		bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	req = req.WithContext(ctx)
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth return status code %d", resp.StatusCode)
	}

	b, _ := io.ReadAll(resp.Body)

	vvv := &validateTokenResponse{}
	var got = models.Response{
		Resp: vvv,
	}
	_ = json.Unmarshal(b, &got)

	vvv = got.Resp.(*validateTokenResponse)
	return vvv, nil
}
