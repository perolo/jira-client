package jira

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

// ProjectService handles projects for the JIRA instance / API.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project
type ProjectService struct {
	client *Client
}

// ProjectList represent a list of Projects
type ProjectList []struct {
	Expand          string          `json:"expand" structs:"expand"`
	Self            string          `json:"self" structs:"self"`
	ID              string          `json:"id" structs:"id"`
	Key             string          `json:"key" structs:"key"`
	Name            string          `json:"name" structs:"name"`
	AvatarUrls      AvatarUrls      `json:"avatarUrls" structs:"avatarUrls"`
	ProjectTypeKey  string          `json:"projectTypeKey" structs:"projectTypeKey"`
	ProjectCategory ProjectCategory `json:"projectCategory,omitempty" structs:"projectsCategory,omitempty"`
}

// ProjectCategory represents a single project category
type ProjectCategory struct {
	Self        string `json:"self" structs:"self,omitempty"`
	ID          string `json:"id" structs:"id,omitempty"`
	Name        string `json:"name" structs:"name,omitempty"`
	Description string `json:"description" structs:"description,omitempty"`
}

// Project represents a JIRA Project.
type Project struct {
	Expand          string             `json:"expand,omitempty" structs:"expand,omitempty"`
	Self            string             `json:"self,omitempty" structs:"self,omitempty"`
	ID              string             `json:"id,omitempty" structs:"id,omitempty"`
	Key             string             `json:"key,omitempty" structs:"key,omitempty"`
	Description     string             `json:"description,omitempty" structs:"description,omitempty"`
	Lead            User               `json:"lead,omitempty" structs:"lead,omitempty"`
	Components      []ProjectComponent `json:"components,omitempty" structs:"components,omitempty"`
	IssueTypes      []IssueType        `json:"issueTypes,omitempty" structs:"issueTypes,omitempty"`
	URL             string             `json:"url,omitempty" structs:"url,omitempty"`
	Email           string             `json:"email,omitempty" structs:"email,omitempty"`
	AssigneeType    string             `json:"assigneeType,omitempty" structs:"assigneeType,omitempty"`
	Versions        []Version          `json:"versions,omitempty" structs:"versions,omitempty"`
	Name            string             `json:"name,omitempty" structs:"name,omitempty"`
	Roles           map[string]string  `json:"roles,omitempty" structs:"roles,omitempty"`
	AvatarUrls      AvatarUrls         `json:"avatarUrls,omitempty" structs:"avatarUrls,omitempty"`
	ProjectCategory ProjectCategory    `json:"projectCategory,omitempty" structs:"projectCategory,omitempty"`
}

//Roles        struct {
//Developers string `json:"Developers,omitempty" structs:"Developers,omitempty"`
//} `json:"roles,omitempty" structs:"roles,omitempty"`

// Version represents a single release version of a project
type Version struct {
	Self            string `json:"self" structs:"self,omitempty"`
	ID              string `json:"id" structs:"id,omitempty"`
	Name            string `json:"name" structs:"name,omitempty"`
	Archived        bool   `json:"archived" structs:"archived,omitempty"`
	Released        bool   `json:"released" structs:"released,omitempty"`
	ReleaseDate     string `json:"releaseDate" structs:"releaseDate,omitempty"`
	UserReleaseDate string `json:"userReleaseDate" structs:"userReleaseDate,omitempty"`
	ProjectID       int    `json:"projectId" structs:"projectId,omitempty"` // Unlike other IDs, this is returned as a number
}

// ProjectComponent represents a single component of a project
type ProjectComponent struct {
	Self                string `json:"self" structs:"self,omitempty"`
	ID                  string `json:"id" structs:"id,omitempty"`
	Name                string `json:"name" structs:"name,omitempty"`
	Description         string `json:"description" structs:"description,omitempty"`
	Lead                User   `json:"lead,omitempty" structs:"lead,omitempty"`
	AssigneeType        string `json:"assigneeType" structs:"assigneeType,omitempty"`
	Assignee            User   `json:"assignee" structs:"assignee,omitempty"`
	RealAssigneeType    string `json:"realAssigneeType" structs:"realAssigneeType,omitempty"`
	RealAssignee        User   `json:"realAssignee" structs:"realAssignee,omitempty"`
	IsAssigneeTypeValid bool   `json:"isAssigneeTypeValid" structs:"isAssigneeTypeValid,omitempty"`
	Project             string `json:"project" structs:"project,omitempty"`
	ProjectID           int    `json:"projectId" structs:"projectId,omitempty"`
}

type PropType struct {
	Self string `json:"self,omitempty" structs:"self,omitempty"`
	Key  string `json:"key,omitempty" structs:"key,omitempty"`
}
type ProjectProperties struct {
	Keys []PropType `json:"keys,omitempty" structs:"keys,omitempty"`
}

// GetList gets all projects form JIRA
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project-getAllProjects
func (s *ProjectService) GetList() (*ProjectList, *Response, error) {
	apiEndpoint := "rest/api/2/project"
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	projectList := new(ProjectList)
	resp, err := s.client.Do(req, projectList)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return projectList, resp, nil
}

// Get returns a full representation of the project for the given issue key.
// JIRA will attempt to identify the project by the projectIdOrKey path parameter.
// This can be an project id, or an project key.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project-getProject
func (s *ProjectService) Get(projectID string) (*Project, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/project/%s", projectID)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(req, project)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return project, resp, nil
}

// Get gets Roles for project from JIRA
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-getUser
func (s *UserService) GetRoles(project string) (*map[string]string, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/project/%s/role", project)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	role := make(map[string]string)
	//fmt.Println("apiEndpoint: " + apiEndpoint)
	resp, err := s.client.Do(req, &role)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return &role, resp, nil
}

// Get gets Role info from JIRA
//
func (s *UserService) GetProjectRole(link string) (*ProjectRole, *Response, error) {
	req, err := s.client.NewRequest("GET", link, nil)
	if err != nil {
		return nil, nil, err
	}

	projrole := new(ProjectRole)
	resp, err := s.client.Do(req, projrole)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return projrole, resp, nil
}

func (s *UserService) SetProjectRole(user *User) (*User, *Response, error) {
	apiEndpoint := "api/2/project/{projectIdOrKey}/role/{id}"
	req, err := s.client.NewRequest("POST", apiEndpoint, user)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	responseUser := new(User)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := fmt.Errorf("Could not read the returned data")
		return nil, resp, NewJiraError(resp, e)
	}
	err = json.Unmarshal(data, responseUser)
	if err != nil {
		e := fmt.Errorf("Could not unmarshall the data into struct")
		return nil, resp, NewJiraError(resp, e)
	}
	return responseUser, resp, nil
}

// Get gets Role info from JIRA
//
func (s *ProjectService) GetProjectProperties(project string) (*ProjectProperties, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/project/%s/properties", project)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	projprop := new(ProjectProperties)
	resp, err := s.client.Do(req, projprop)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return projprop, resp, nil
}
