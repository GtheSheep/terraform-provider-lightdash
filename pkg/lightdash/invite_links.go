package lightdash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type InviteLink struct {
	InviteCode       string `json:"inviteCode"`
	ExpiresAt        string `json:"expiresAt"`
	InviteUrl        string `json:"inviteUrl"`
	OrganizationUUID string `json:"organizationUuid"`
	UserUUID         string `json:"userUuid"`
	Email            string `json:"email"`
}

type InviteLinkRequest struct {
	ExpiresAt string `json:"expiresAt"`
	Email     string `json:"email"`
}

type InviteLinkResponse struct {
	Results InviteLink `json:"results"`
	Status  string     `json:"status"`
}

func (c *Client) GetInviteLink(inviteCode string) (*InviteLink, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/invite-links/%s", c.ApiURL, inviteCode), nil)
	if err != nil {
		return nil, err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	inviteLinkResponse := InviteLinkResponse{}
	err = json.Unmarshal(body, &inviteLinkResponse)
	if err != nil {
		return nil, err
	}
	return &inviteLinkResponse.Results, nil
}

func (c *Client) CreateInviteLink(email string) (*InviteLink, error) {
	newInviteLinkRequest := InviteLinkRequest{
		ExpiresAt: "2099-01-01T23:59:59Z",
		Email:     email,
	}
	newInviteLinkRequestData, err := json.Marshal(newInviteLinkRequest)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/invite-links", c.ApiURL), strings.NewReader(string(newInviteLinkRequestData)))
	if err != nil {
		return nil, err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	inviteLinkResponse := InviteLinkResponse{}
	err = json.Unmarshal(body, &inviteLinkResponse)
	if err != nil {
		return nil, err
	}

	return &inviteLinkResponse.Results, nil
}

func (c *Client) DeleteInviteLink(inviteCode string) (string, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/invite-links/%s", c.ApiURL, inviteCode), nil)
	if err != nil {
		return "", err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return "", err
	}

	inviteLinkResponse := InviteLinkResponse{}
	err = json.Unmarshal(body, &inviteLinkResponse)
	if err != nil {
		return "", err
	}

	return inviteLinkResponse.Status, nil
}

func (c *Client) DeleteAllInviteLinks() (string, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/invite-links", c.ApiURL), nil)
	if err != nil {
		return "", err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return "", err
	}

	inviteLinkResponse := InviteLinkResponse{}
	err = json.Unmarshal(body, &inviteLinkResponse)
	if err != nil {
		return "", err
	}

	return inviteLinkResponse.Status, nil
}
