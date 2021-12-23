package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// RoleService handles roles for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-group-Role
type RoleService struct {
	client *Client
}

// Role represents a Jira product role
type Role struct {
	Self        string   `json:"self" structs:"self"`
	Name        string   `json:"name" structs:"name"`
	ID          int      `json:"id" structs:"id"`
	Description string   `json:"description" structs:"description"`
	Actors      []*Actor `json:"actors" structs:"actors"`
}

// Actor represents a Jira actor
type Actor struct {
	ID          int        `json:"id" structs:"id"`
	DisplayName string     `json:"displayName" structs:"displayName"`
	Type        string     `json:"type" structs:"type"`
	Name        string     `json:"name" structs:"name"`
	AvatarURL   string     `json:"avatarUrl" structs:"avatarUrl"`
	ActorUser   *ActorUser `json:"actorUser" structs:"actoruser"`
}

// ActorUser contains the account id of the actor/user
type ActorUser struct {
	AccountID string `json:"accountId" structs:"accountId"`
}

// GetListWithContext returns a list of all available project roles
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-api-3-role-get
func (s *RoleService) GetListWithContext(ctx context.Context) (*[]Role, *Response, error) {
	apiEndpoint := "rest/api/3/role"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	roles := new([]Role)
	resp, err := s.client.Do(req, roles)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}
	return roles, resp, err
}

// GetList wraps GetListWithContext using the background context.
func (s *RoleService) GetList() (*[]Role, *Response, error) {
	return s.GetListWithContext(context.Background())
}

// GetWithContext retreives a single Role from Jira
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-api-3-role-id-get
func (s *RoleService) GetWithContext(ctx context.Context, roleID int) (*Role, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/3/role/%d", roleID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	role := new(Role)
	resp, err := s.client.Do(req, role)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}
	if role.Self == "" {
		return nil, resp, fmt.Errorf("no role with ID %d found", roleID)
	}

	return role, resp, err
}

// Get wraps GetWithContext using the background context.
func (s *RoleService) Get(roleID int) (*Role, *Response, error) {
	return s.GetWithContext(context.Background(), roleID)
}

type RoleType struct {
	Name   string
	Rollnk string
	ID     string
}

// GetListWithContext returns a list of all available roles for a project
// /rest/api/2/project/{projectIdOrKey}/role
// JIRA API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-api-3-role-get
func (s *RoleService) GetRolesForProjectWithContext(ctx context.Context, proj string) (*[]RoleType, *Response, error) {
	var rl []RoleType
	apiEndpoint := fmt.Sprintf("/rest/api/latest/project/%s/role", proj)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do2(req)
	if err != nil {
		//jerr := NewJiraError(resp, err)
		return nil, nil, nil
	}
	var doc interface{}
	if err := json.Unmarshal(resp, &doc); err != nil {
		return nil, nil, nil
	}
	// Should be a better way of doing this:
	for k, v := range doc.(map[string]interface{}) {
		var r RoleType
		r.Name = k
		r.Rollnk = v.(string)
		pos := strings.LastIndex(r.Rollnk, "/role/")
		adjustedPos := pos + len("/role/")
		r.ID = r.Rollnk[adjustedPos:len(r.Rollnk)]
		rl = append(rl, r)
	}
	return &rl, nil, err
}

// /rest/api/2/project/{projectIdOrKey}/role/{id}
func (s *RoleService) GetActorsForProjectRoleWithContext(ctx context.Context, proj string, roleid string) (*Role, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/project/%s/role/%s", proj, roleid)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	roles := new(Role)
	resp, err := s.client.Do(req, roles)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}
	return roles, resp, err
}

type GroupAddType struct {
	Group []string `json:"group"`
}

// /rest/api/2/project/{projectIdOrKey}/role/{id}
func (s *RoleService) AddActorsForProjectRoleWithContext(ctx context.Context, proj string, roleid string, actor string) (*Role, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/project/%s/role/%s", proj, roleid)
	var payload = new(GroupAddType)
	payload.Group = append(payload.Group, actor)
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, payload)
	if err != nil {
		return nil, nil, err
	}
	roles := new(Role)
	resp, err := s.client.Do(req, roles)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}
	return roles, resp, err
	/*
		var user struct {
			Name string `json:"name"`
		}
		user.Name = username
		req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, &user)
		if err != nil {
			return nil, nil, err
		}
	*/
	//groups := new(AddGroupsResponseType)
	//c.doRequest("POST", u, payload, &groups)
}

type GroupRemoveType struct {
	User []string `json:"user"`
}

// DELETE /rest/project/{projectIdOrKey}/role/{id}
func (s *RoleService) RemoveUserActorsForProjectRole(proj string, roleid int, user string) (*Role, *Response, error) {
	ctx := context.Background()
	apiEndpoint := fmt.Sprintf("/rest/api/2/project/%s/role/%v?user=%s", proj, roleid, user)
	req, err := s.client.NewRequestWithContext(ctx, "DELETE", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}
	return nil, resp, err
}
