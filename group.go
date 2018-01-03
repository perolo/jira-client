package jira

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	//"os/user"
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

// Group represents a JIRA group
type Group struct {
	ID                   string          `json:"id,omitempty"`
	Title                string          `json:"title,omitempty"`
	Type                 string          `json:"type,omitempty"`
	Properties           groupProperties `json:"properties,omitempty"`
	AdditionalProperties bool            `json:"additionalProperties,omitempty"`
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
	EmailAddress string `json:"emailAddress,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	Active       bool   `json:"active,omitempty"`
	TimeZone     string `json:"timeZone,omitempty"`
}
type Groups struct {
	Name         string `json:"name,omitempty" structs:"name,omitempty"`
	Html         string `json:"html,omitempty" structs:"html,omitempty"`
	labels       []string `json:"labels,omitempty"  structs:"labels,omitempty`
}

type GroupResp struct {
	Self         string `json:"self,omitempty"`
	Name         string `json:"name,omitempty"`
	Expand       string `json:"expand,omitempty"`
	Users        string `json:"users,omitempty" structs:"users,omitempty"`
}

type GroupsType2 struct {
	Header 	   string        `json:"header,omitempty" structs:"header,omitempty"`
	Total      int           `json:"total,omitempty" structs:"total,omitempty"`
	Groups     []Groups      `json:"groups,omitempty"  structs:"groups,omitempty`
}

func (s *GroupService) GetGroups(opt *GroupOptions) (*GroupsType2, *Response, error) {
	var u string
	if opt == nil {
		u = fmt.Sprintf("rest/api/2/groups/picker")
	} else {
		u = fmt.Sprintf("rest/api/2/groups/picker?startAt=%d&maxResults=%d", opt.StartAt, opt.MaxResults)
	}
	apiEndpoint := u
	//url, err := addOptions(apiEndpoint, opt)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	groups := new(GroupsType2)
	resp, err := s.client.Do(req, groups)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return groups, resp, err
}

func (s *GroupService) GetGroups2(options *GroupOptions) (*GroupsType2, *Response, error, int) {
	var u string
	if options == nil {
		u = fmt.Sprintf("/rest/api/2/groups/picker")
	} else {
		u = fmt.Sprintf("/rest/api/2/groups/picker&maxResults=%d", options.MaxResults)
	}
//	fmt.Println("u: " + u)
	apiEndpoint := u
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err, 0
	}

	groups := new(GroupsType2)
	resp, err := s.client.Do(req, groups)
	if err != nil {
		return nil, resp, err,0
	}

	return groups, resp, nil, groups.Total
}


// Get returns a paginated list of users who are members of the specified group and its subgroups.
// Users in the page are ordered by user names.
// User of this resource is required to have sysadmin or admin permissions.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/server/#api/2/group-getUsersFromGroup
func (s *GroupService) Get1(name string, options *GroupOptions) ([]GroupMember, *Response, error, int) {
	var u string
	if options == nil {
		u = fmt.Sprintf("rest/api/2/group/member?groupname=%s", name)
	} else {
		u = fmt.Sprintf("rest/api/2/group/member?groupname=%s&startAt=%d&maxResults=%d", name,
			options.StartAt, options.MaxResults)
	}
//	fmt.Println("u: " + u)
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



// Get returns a paginated list of users who are members of the specified group and its subgroups.
// Users in the page are ordered by user names.
// User of this resource is required to have sysadmin or admin permissions.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/server/#api/2/group-getUsersFromGroup
func (s *GroupService) Get2(name string) ([]GroupMember, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/group/member?groupname=%s", name)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
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

// Add adds user to group
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/group-addUserToGroup
func (s *GroupService) Add(groupname string, username string) (*GroupResp, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/group/user?groupname=%s", groupname)
	var user struct {
		Name string `json:"name"`
	}
	user.Name = username

	//fmt.Println("apiEndpoint: " + apiEndpoint)
	//fmt.Println("name: " + user.Name)

	req, err := s.client.NewRequest("POST", apiEndpoint, &user)
	if err != nil {
		fmt.Println("Not OK: " + user.Name)
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		// in case of error return the resp for further inspection
		return nil, resp, err
	}

	responseGroup := new(GroupResp)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not read the returned data")
	}
	//fmt.Println("resp.Response.Status: " + resp.Response.Status)
	//fmt.Println("Data: " + string(data))

	err = json.Unmarshal(data, responseGroup)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not unmarshall the data into struct")
	}
	return responseGroup, resp, nil
}

// Remove removes user from group
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/group-removeUserFromGroup
func (s *GroupService) Remove(groupname string, username string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/group/user?groupname=%s&username=%s", groupname, username)
	req, err := s.client.NewRequest("DELETE", apiEndpoint, nil)
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
