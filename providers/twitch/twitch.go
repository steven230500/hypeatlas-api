package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	oauthURL = "https://id.twitch.tv/oauth2/token"
	baseAPI  = "https://api.twitch.tv/helix"
)

type Client struct {
	ClientID string
	Secret   string
	http     *http.Client
	token    string
	exp      time.Time
}

func New(clientID, secret string) *Client {
	return &Client{
		ClientID: clientID,
		Secret:   secret,
		http:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) ensureToken(ctx context.Context) error {
	if c.token != "" && time.Now().Before(c.exp.Add(-1*time.Minute)) {
		return nil
	}
	form := url.Values{}
	form.Set("client_id", c.ClientID)
	form.Set("client_secret", c.Secret)
	form.Set("grant_type", "client_credentials")

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, oauthURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return fmt.Errorf("twitch oauth: %s", res.Status)
	}

	var body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return err
	}

	c.token = body.AccessToken
	c.exp = time.Now().Add(time.Duration(body.ExpiresIn) * time.Second)
	return nil
}

func (c *Client) auth(req *http.Request) {
	req.Header.Set("Client-Id", c.ClientID)
	req.Header.Set("Authorization", "Bearer "+c.token)
}

type Stream struct {
	UserID       string `json:"user_id"`
	UserLogin    string `json:"user_login"`
	Language     string `json:"language"`
	Title        string `json:"title"`
	ViewerCount  int    `json:"viewer_count"`
	Type         string `json:"type"` // "live" o ""
	StartedAt    string `json:"started_at"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// GetStreamsByLogin acepta hasta 100 logins por llamada.
func (c *Client) GetStreamsByLogin(ctx context.Context, logins []string) (map[string]Stream, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}

	u, _ := url.Parse(baseAPI + "/streams")
	q := u.Query()
	for _, l := range logins {
		l = strings.TrimSpace(strings.ToLower(l))
		if l != "" {
			q.Add("user_login", l)
		}
	}
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	c.auth(req)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("twitch streams: %s", res.Status)
	}

	var body struct {
		Data []Stream `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	out := map[string]Stream{}
	for _, s := range body.Data {
		out[strings.ToLower(s.UserLogin)] = s
	}
	return out, nil
}

// Chunk parte un slice en trozos de tamaño n.
func Chunk[T any](xs []T, n int) [][]T {
	if n <= 0 {
		n = 1
	}
	var out [][]T
	for len(xs) > 0 {
		k := n
		if len(xs) < k {
			k = len(xs)
		}
		out = append(out, xs[:k])
		xs = xs[k:]
	}
	return out
}
