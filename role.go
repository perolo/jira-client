package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qri-io/jsonschema"
	"strings"
)

// RoleService handles roles for the JIRA instance / API.
//
// JIRA API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-group-Role
type RoleService struct {
	client *Client
}

// Role represents a JIRA product role
type Role struct {
	Self        string   `json:"self" structs:"self"`
	Name        string   `json:"name" structs:"name"`
	ID          int      `json:"id" structs:"id"`
	Description string   `json:"description" structs:"description"`
	Actors      []*Actor `json:"actors" structs:"actors"`
}

// Actor represents a JIRA actor
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
// JIRA API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-api-3-role-get
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
// JIRA API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-api-3-role-id-get
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
	Name string
	Rollnk string
	ID string
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
	var schemaData = []byte(`{
	"type": "object",
		"patternProperties": {
		".+": {
			"type": "string",
				"format": "uri"
		}
	},
	"additionalProperties": false
}`)
	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		panic("unmarshal schema: " + err.Error())
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
	for k,v := range doc.(map[string]interface{}) { // .(map[string]string)
		var r RoleType
		r.Name = k
		r.Rollnk = v.(string)
		pos := strings.LastIndex(r.Rollnk, "/role/")
		adjustedPos := pos + len("/role/")
		r.ID =  r.Rollnk[adjustedPos:len(r.Rollnk)]
		//r.ID = strings.TrimLeft(v.(string), "/role/")

		rl = append(rl, r)
	}

	if errors, _ := rs.ValidateBytes(resp); len(errors) > 0 {
		panic(errors)
	}

/*
	roles := new([]Role)
	resp, err := s.client.Do(req, roles)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

 */
	return &rl, nil, err
}
/*
type ActorStruct struct {
	Self        string `json:"self"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
	Description string `json:"description"`
	Actors      []struct {
		ID          int    `json:"id"`
		DisplayName string `json:"displayName"`
		Type        string `json:"type"`
		Name        string `json:"name"`
		AvatarURL   string `json:"avatarUrl"`
	} `json:"actors"`
}
*/
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
