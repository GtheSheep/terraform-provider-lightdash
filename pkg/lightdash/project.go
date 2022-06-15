package lightdash

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Project struct {
	ProjectUUID string `json:"projectUuid,omitempty"`
	Name        string `json:"name,omitempty"`
}

type ProjectResponse struct {
	Results Project `json:"results"`
	Status  string  `json:"status"`
}

type ProjectsResponse struct {
	Results []Project `json:"results"`
	Status  string    `json:"status"`
}

func (c *Client) GetProject(projectUUID string) (*Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/org/projects", c.ApiURL), nil)
	if err != nil {
		return nil, err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	projectsResponse := ProjectsResponse{}
	err = json.Unmarshal(body, &projectsResponse)
	if err != nil {
		return nil, err
	}

	for i, project := range projectsResponse.Results {
		if project.ProjectUUID == projectUUID {
			return &projectsResponse.Results[i], nil
		}
	}

	return nil, fmt.Errorf("Project not found UUID %s", projectUUID)
}

// func (c *Client) CreateProject(name string) (*Project, error) {
//
// 	return &newProject, nil
// }
//
// func (c *Client) UpdateProject(projectUUID string) (*Project, error) {
//
// 	return &projectResponse.Results, nil
// }

func (c *Client) DeleteProject(projectUUID string) (string, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/org/projects/%s", c.ApiURL, projectUUID), nil)
	if err != nil {
		return "", err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return "", err
	}

	projectResponse := ProjectResponse{}
	err = json.Unmarshal(body, &projectResponse)
	if err != nil {
		return "", err
	}

	return projectResponse.Status, nil
}
