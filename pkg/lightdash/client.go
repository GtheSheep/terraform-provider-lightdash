package lightdash

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	URL        string
	HTTPClient *http.Client
	Username   string
	Password   string
	ApiURL     string
	Cookies    []*http.Cookie
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResults struct {
	UserUUID             string
	Email                string
	FirstName            string
	LastName             string
	OrganizationID       string
	OrganizationName     string
	IsTrackingAnonymized bool
	IsMarketingOptedIn   bool
	IsSetupComplete      bool
	Role                 string
	IsActive             bool
}

type LoginResponse struct {
	Status  string       `json:"status`
	Results LoginResults `json:"results`
}

// TODO: Convert to use a session
func NewClient(url *string, username *string, password *string) (*Client, error) {
	c := Client{
		URL:        *url,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Username:   *username,
		Password:   *password,
		ApiURL:     fmt.Sprintf("%s/api/v1", *url),
	}

	if (url != nil) && (username != nil) && (password != nil) {
		loginRequest := LoginRequest{
			Email:    *username,
			Password: *password,
		}
		loginRequestData, err := json.Marshal(loginRequest)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/login", c.ApiURL), strings.NewReader(string(loginRequestData)))
		if err != nil {
			return nil, err
		}

		body, err, cookies := c.doRequest(req)

		lr := LoginResponse{}
		err = json.Unmarshal(body, &lr)
		if err != nil {
			return nil, err
		}

		c.Cookies = cookies
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error, []*http.Cookie) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	for _, cookie := range c.Cookies {
		req.AddCookie(cookie)
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err, nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err, nil
	}

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != 201) {
		return nil, fmt.Errorf("%s url: %s, status: %d, body: %s", req.Method, req.URL, res.StatusCode, body), nil
	}

	return body, err, res.Cookies()
}
