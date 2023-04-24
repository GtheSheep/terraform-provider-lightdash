package lightdash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type User struct {
	UserUUID         string  `json:"userUuid,omitempty"`
	FirstName        string  `json:"firstName,omitempty"`
	LastName         string  `json:"lastName,omitempty"`
	Email            string  `json:"email,omitempty"`
	Password         string  `json:"password,omitempty"`
	OrganizationUUID string  `json:"organizationUuid,omitempty"`
	Role             string  `json:"role,omitempty"`
	IsActive         bool    `json:"isActive,omitempty"`
	IsInviteExpired  bool    `json:"isInviteExpired,omitempty"`
	InviteCode       *string `json:"inviteCode,omitempty"`
}

type UserResponse struct {
	Results User   `json:"results"`
	Status  string `json:"status"`
}

type UsersResponse struct {
	Results []User `json:"results"`
	Status  string `json:"status"`
}

func (c *Client) GetUser(userUUID string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/org/users", c.ApiURL), nil)
	if err != nil {
		return nil, err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	usersResponse := UsersResponse{}
	err = json.Unmarshal(body, &usersResponse)
	if err != nil {
		return nil, err
	}

	for i, user := range usersResponse.Results {
		if user.UserUUID == userUUID {
			return &usersResponse.Results[i], nil
		}
	}

	return nil, fmt.Errorf("User not found UUID %s", userUUID)
}

func (c *Client) UpdateUser(userID string, role string) (*User, error) {
	updatedUser := User{
		Role: role,
	}
	updatedUserData, err := json.Marshal(updatedUser)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/org/users/%s", c.ApiURL, userID), strings.NewReader(string(updatedUserData)))
	if err != nil {
		return nil, err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	userResponse := UserResponse{}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, err
	}

	return &userResponse.Results, nil
}

func (c *Client) DeleteUser(userID string) (string, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/org/user/%s", c.ApiURL, userID), nil)
	if err != nil {
		return "", err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return "", err
	}

	userResponse := UserResponse{}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return "", err
	}

	return userResponse.Status, nil
}
