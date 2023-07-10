package spotify

import (
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/karim-w/stdlib/httpclient"
)

type Client struct {
	token        TokenResponse
	expiresAt    time.Time
	clientId     string
	clientSecret string
}

const baseURL = "https://api.spotify.com"

func InitOrDieClient(
	code string,
	clientId string,
	clientSecret string,
) *Client {
	formData := url.Values{
		"grant_type":   {"authorization_code"},
		"redirect_uri": {"http://localhost:8080/callback"},
		"code":         {code},
	}
	res := httpclient.Req("https://accounts.spotify.com/api/token").
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		AddBodyRaw([]byte(formData.Encode())).
		AddBasicAuth(clientId, clientSecret).Post()
	if !res.IsSuccess() {
		body := res.GetBody()
		log.Println(string(body))
		panic("failed to get request")
	}
	var token TokenResponse
	res.SetResult(&token)
	expiresAt := time.UnixMicro(time.Now().UnixMicro() + int64(token.ExpiresIn-10))
	return &Client{
		token,
		expiresAt,
		clientId,
		clientSecret,
	}
}

func (c *Client) GetCurrentlyPlayingSong() (*GetCurrentPlayingSongQueryResponse, error) {
	res := httpclient.Req(baseURL + "/v1/me/player/currently-playing").
		AddBearerAuth(c.token.AccessToken).
		Get()
	body := res.GetBody()
	if !res.IsSuccess() {
		log.Println(res.GetStatusCode())
		return nil, errors.New("failed to get current song with " + string(body))
	}
	var response GetCurrentPlayingSongQueryResponse
	res.SetResult(&response)
	return &response, nil
}

func (c *Client) getAccessToken() string {
	if time.Now().Before(c.expiresAt) {
		return c.token.AccessToken
	}
	formData := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.token.RefreshToken},
	}
	res := httpclient.Req("https://accounts.spotify.com/api/token").
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		AddBodyRaw([]byte(formData.Encode())).
		AddBasicAuth(c.clientId, c.clientSecret).Post()
	if !res.IsSuccess() {
		body := res.GetBody()
		log.Println(string(body))
		panic("failed to get request")
	}
	var token TokenResponse
	res.SetResult(&token)
	expiresAt := time.UnixMicro(time.Now().UnixMicro() + int64(token.ExpiresIn-10))
	c.expiresAt = expiresAt
	c.token = token
	return token.AccessToken
}
