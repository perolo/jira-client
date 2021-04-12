package jira

import (
	"context"
	"fmt"
	"net/url"
)

// GroupService handles Groups for the JIRA instance / API.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/server/#api/2/group
type GroupService struct {
	client *Client
}

// groupMembersResult is only a small wrapper around the Group* methods
// to be able to parse the results
type groupMembersResult struct {
	StartAt    int           `json:"startAt"`
	MaxResults int           `json:"maxResults"`
	Total      int           `json:"total"`
	Members    []GroupMember `json:"values"`
}

// Group represents a JIRA group
type Group struct {
	ID                   string          `json:"id"`
	Title                string          `json:"title"`
	Type                 string          `json:"type"`
	Properties           groupProperties `json:"properties"`
	AdditionalProperties bool            `json:"additionalProperties"`
}

type groupProperties struct {
	Name groupPropertiesName `json:"name"`
}

type groupPropertiesName struct {
	Type string `json:"type"`
}

// GroupMember reflects a single member of a group
type GroupMember struct {
	Self         string `json:"self,omitempty"`
	Name         string `json:"name,omitempty"`
	Key          string `json:"key,omitempty"`
	AccountID    string `json:"accountId,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	Active       bool   `json:"active,omitempty"`
	TimeZone     string `json:"timeZone,omitempty"`
	AccountType  string `json:"accountType,omitempty"`
}

// GroupSearchOptions specifies the optional parameters for the Get Group methods
type GroupSearchOptions struct {
	StartAt              int
	MaxResults           int
	IncludeInactiveUsers bool
}

type PermissionSearchOptions struct {
	StartAt  int    `url:"startAt,omitempty"`
	MaxResults  int    `url:"maxResults,omitempty"`
	UserName  string `url:"username,omitempty"`
	Permissions   string `url:"permissions,omitempty"`
	IssueKey string `url:"issueKey,omitempty"`
	ProjectKey string `url:"projectKey,omitempty"`
}

type PermissionSearchResultType []struct {
	Self       string `json:"self"`
	Name       string `json:"name"`
	AvatarUrls struct {
		Two4X24   string `json:"24x24"`
		One6X16   string `json:"16x16"`
		Three2X32 string `json:"32x32"`
		Four8X48  string `json:"48x48"`
	} `json:"avatarUrls"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
}
// GetWithContext returns a paginated list of users who are members of the specified group and its subgroups.
// Users in the page are ordered by user names.
// User of this resource is required to have sysadmin or admin permissions.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/server/#api/2/group-getUsersFromGroup
//
// WARNING: This API only returns the first page of group members
func (s *GroupService) GetWithContext(ctx context.Context, name string) ([]GroupMember, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/group/member?groupname=%s", url.QueryEscape(name))
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	group := new(groupMembersResult)
	resp, err := s.client.Do(req, group)
	if err != nil {
		return nil, resp, err
	}

	return group.Members, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *GroupService) Get(name string) ([]GroupMember, *Response, error) {
	return s.GetWithContext(context.Background(), name)
}

// GetWithOptionsWithContext returns a paginated list of members of the specified group and its subgroups.
// Users in the page are ordered by user names.
// User of this resource is required to have sysadmin or admin permissions.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/server/#api/2/group-getUsersFromGroup
func (s *GroupService) GetWithOptionsWithContext(ctx context.Context, name string, options *GroupSearchOptions) ([]GroupMember, *Response, error) {
	var apiEndpoint string
	if options == nil {
		apiEndpoint = fmt.Sprintf("/rest/api/2/group/member?groupname=%s", url.QueryEscape(name))
	} else {
		apiEndpoint = fmt.Sprintf(
			"/rest/api/2/group/member?groupname=%s&startAt=%d&maxResults=%d&includeInactiveUsers=%t",
			url.QueryEscape(name),
			options.StartAt,
			options.MaxResults,
			options.IncludeInactiveUsers,
		)
	}
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	group := new(groupMembersResult)
	resp, err := s.client.Do(req, group)
	if err != nil {
		return nil, resp, err
	}
	return group.Members, resp, nil
}

// GetWithOptions wraps GetWithOptionsWithContext using the background context.
func (s *GroupService) GetWithOptions(name string, options *GroupSearchOptions) ([]GroupMember, *Response, error) {
	return s.GetWithOptionsWithContext(context.Background(), name, options)
}

//	/rest/api/2/user/permission/search

func (s *GroupService) SearchPermissionsWithOptionsWithContext(ctx context.Context, permissiosn string, options *PermissionSearchOptions) (*PermissionSearchResultType, *Response, error) {
	var apiEndpoint string
	apiEndpoint, _ = addOptions("/rest/api/2/user/permission/search", options)

	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	group := new(PermissionSearchResultType)
	resp, err := s.client.Do(req, group)
	if err != nil {
		return nil, resp, err
	}
	return group, resp, nil
}

// AddWithContext adds user to group
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/group-addUserToGroup
func (s *GroupService) AddWithContext(ctx context.Context, groupname string, username string) (*Group, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/group/user?groupname=%s", groupname)
	var user struct {
		Name string `json:"name"`
	}
	user.Name = username
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, &user)
	if err != nil {
		return nil, nil, err
	}

	responseGroup := new(Group)
	resp, err := s.client.Do(req, responseGroup)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return responseGroup, resp, nil
}

// Add wraps AddWithContext using the background context.
func (s *GroupService) Add(groupname string, username string) (*Group, *Response, error) {
	return s.AddWithContext(context.Background(), groupname, username)
}

// RemoveWithContext removes user from group
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/group-removeUserFromGroup
func (s *GroupService) RemoveWithContext(ctx context.Context, groupname string, username string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/group/user?groupname=%s&username=%s", groupname, username)
	req, err := s.client.NewRequestWithContext(ctx, "DELETE", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return resp, jerr
	}

	return resp, nil
}

// Remove wraps RemoveWithContext using the background context.
func (s *GroupService) Remove(groupname string, username string) (*Response, error) {
	return s.RemoveWithContext(context.Background(), groupname, username)
}

type GroupsResult struct {
	Header string `json:"header"`
	Total  int    `json:"total"`
	Groups []struct {
		Name   string        `json:"name"`
		HTML   string        `json:"html"`
		Labels []interface{} `json:"labels"`
	} `json:"groups"`
}

// Get wraps GetWithContext using the background context.
func (s *GroupService) GetGroups() (*GroupsResult, *Response, error) {
	return s.GetGroupsWithContext(context.Background())
}
func (s *GroupService) GetGroupsWithContext(ctx context.Context) (*GroupsResult, *Response, error) {
	///rest/api/2/groups/picker?query&exclude&maxResults
	apiEndpoint := fmt.Sprintf("/rest/api/2/groups/picker?maxResults=10000")
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	groups := new(GroupsResult)
	resp, err := s.client.Do(req, groups)
	if err != nil {
		return nil, resp, err
	}

	return groups, resp, nil
}
