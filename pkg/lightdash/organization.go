package lightdash

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Organization struct {
	UUID string `json:"organizationUuid,omitempty"`
	Name string `json:"name,omitempty"`
}

type OrganizationResponse struct {
	Results Organization `json:"results"`
	Status  string       `json:"status"`
}

func (c *Client) GetOrganization() (*Organization, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/org", c.ApiURL), nil)
	if err != nil {
		return nil, err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	organizationResponse := OrganizationResponse{}
	err = json.Unmarshal(body, &organizationResponse)
	if err != nil {
		return nil, err
	}

	return &organizationResponse.Results, nil
}
