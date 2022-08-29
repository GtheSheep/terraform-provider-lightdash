package lightdash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Project struct {
	ProjectUUID         string              `json:"projectUuid,omitempty"`
	Name                string              `json:"name,omitempty"`
	OrganisationUUID    string              `json:"organizationUuid"`
	Type                string              `json:"type"`
	DbtConnection       DbtConnection       `json:"dbtConnection"`
	WarehouseConnection WarehouseConnection `json:"warehouseConnection"`
}

type CreateProjectRequest struct {
	OrganisationUUID    string              `json:"organizationUuid"`
	Name                string              `json:"name"`
	Type                string              `json:"type"`
	DbtConnection       DbtConnection       `json:"dbtConnection"`
	WarehouseConnection WarehouseConnection `json:"warehouseConnection"`
}

type CreateProjectResponse struct {
	Results Project `json:"results"`
	Status  string  `json:"status"`
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

func (c *Client) CreateProject(organisationUUID, name, projectType string, dbtConnection DbtConnection, warehouseConnection WarehouseConnection) (*Project, error) {
	createProjectRequest := CreateProjectRequest{
		OrganisationUUID:    organisationUUID,
		Name:                name,
		Type:                projectType,
		DbtConnection:       dbtConnection,
		WarehouseConnection: warehouseConnection,
	}
	newProjectData, err := json.Marshal(createProjectRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/org/projects", c.ApiURL), strings.NewReader(string(newProjectData)))
	if err != nil {
		return nil, err
	}

	body, err, _ := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	createProjectResponse := CreateProjectResponse{}
	err = json.Unmarshal(body, &createProjectResponse)
	if err != nil {
		return nil, err
	}
	return &createProjectResponse.Results, nil
}

func (c *Client) UpdateProject(projectUUID string) (*Project, error) {

	// TODO: Implement Updates
	project, err := c.GetProject(projectUUID)
	if err != nil {
		return nil, err
	}

	return project, nil
}

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
