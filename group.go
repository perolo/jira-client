package jira

import (
	"fmt"
//	"net/url"
)

// GroupService handles Groups for the JIRA instance / API.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/server/#api/2/group
type GroupService struct {
	client *Client
}
// GroupOptions specifies the optional parameters to various List methods that
// support pagination.
// Pagination is used for the JIRA REST APIs to conserve server resources and limit
// response size for resources that return potentially large collection of items.
// A request to a pages API will result in a values array wrapped in a JSON object with some paging metadata
// Default Pagination options
type GroupOptions struct {
	// StartAt: The starting index of the returned projects. Base index: 0.
	StartAt int `url:"startAt,omitempty"`
	// MaxResults: The maximum number of projects to return per page. Default: 50.
	MaxResults int `url:"maxResults,omitempty"`
	// Expand: Expand specific sections in the returned issues
	Expand string `url:expand,omitempty"`
}

// groupMembersResult is only a small wrapper around the Group* methods
// to be able to parse the results
type groupMembersResult struct {
	StartAt    int           `json:"startAt"`
	MaxResults int           `json:"maxResults"`
	Total      int           `json:"total"`
	Members    []GroupMember `json:"values"`
}

// GroupMember reflects a single member of a group
type GroupMember struct {
	Self         string `json:"self,omitempty"`
	Name         string `json:"name,omitempty"`
	Key          string `json:"key,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	Active       bool   `json:"active,omitempty"`
	TimeZone     string `json:"timeZone,omitempty"`
}

// Get returns a paginated list of users who are members of the specified group and its subgroups.
// Users in the page are ordered by user names.
// User of this resource is required to have sysadmin or admin permissions.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/server/#api/2/group-getUsersFromGroup
func (s *GroupService) Get(name string, options *GroupOptions) ([]GroupMember, *Response, error, int) {
	var u string
	if options == nil {
		u = fmt.Sprintf("rest/api/2/group/member?groupname=%s", name)
	} else {
		u = fmt.Sprintf("rest/api/2/group/member?groupname=%s&startAt=%d&maxResults=%d", name,
			options.StartAt, options.MaxResults)
	}
	fmt.Println("u: " + u)
	apiEndpoint := u
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err, 0
	}

	group := new(groupMembersResult)
	resp, err := s.client.Do(req, group)
	if err != nil {
		return nil, resp, err,0
	}

	return group.Members, resp, nil, group.Total
}
/*
func (s *GroupService) Get(name string) ([]GroupMember, *Response, error) {
	return GetOpt(name, nil)
}*/

//Perolo
func (s *GroupService) PermissionSearch(projid string, options *GroupOptions) ([]GroupMember, *Response, error, int) {
	var u string
	if options == nil {
		u = fmt.Sprintf("rest/api/2/user/permission/search?permissions=BROWSE&projectKey=%s", projid)
	} else {
		u = fmt.Sprintf("rest/api/2/group/member?groupname=%s&startAt=%d&maxResults=%d", projid,
			options.StartAt, options.MaxResults)
	}
	fmt.Println("u: " + u)
	apiEndpoint := u
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err, 0
	}

	group := new([]GroupMember)
	resp, err := s.client.Do(req, group)
	if err != nil {
		return nil, resp, err,0
	}

	return *group, resp, nil, 42
}
