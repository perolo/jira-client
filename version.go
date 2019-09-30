package jira

import (
	"fmt"
	"github.com/google/go-querystring/query"
)

/*
const (
	// AssigneeAutomatic represents the value of the "Assignee: Automatic" of JIRA
	AssigneeAutomatic = "-1"
)
*/

// IssueService handles Issues for the JIRA instance / API.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/issue
type VersionService struct {
	client *Client
}

type Version2 struct {
	Self            string `json:"self"`
	ID              string `json:"id"`
	Description     string `json:"description"`
	Name            string `json:"name"`
	Archived        bool   `json:"archived"`
	Released        bool   `json:"released"`
	StartDate       string `json:"startDate"`
	ReleaseDate     string `json:"releaseDate"`
	Overdue         bool   `json:"overdue"`
	UserStartDate   string `json:"userStartDate"`
	UserReleaseDate string `json:"userReleaseDate"`
	ProjectID       int    `json:"projectId"`
}

func (s *VersionService) Get(versionID string, options *GetQueryOptions) (*Version2, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/latest/version/%s", versionID)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, nil, err
		}
		req.URL.RawQuery = q.Encode()
	}

	version := new(Version2)
	resp, err := s.client.Do(req, version)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return version, resp, nil
}

type RelatedIssueCounts struct {
	Self                                     string `json:"self"`
	IssuesFixedCount                         int    `json:"issuesFixedCount"`
	IssuesAffectedCount                      int    `json:"issuesAffectedCount"`
	IssueCountWithCustomFieldsShowingVersion int    `json:"issueCountWithCustomFieldsShowingVersion"`
}


func (s *VersionService) GetRelatedIssueCounts(versionID string, options *GetQueryOptions) (*RelatedIssueCounts, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/latest/version/%s/relatedIssueCounts", versionID)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, nil, err
		}
		req.URL.RawQuery = q.Encode()
	}

	relatedissuecounts := new(RelatedIssueCounts)
	resp, err := s.client.Do(req, relatedissuecounts)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return relatedissuecounts, resp, nil
}

type IssuesUnresolvedCount struct {
	Self                  string `json:"self"`
	IssuesUnresolvedCount int    `json:"issuesUnresolvedCount"`
}

func (s *VersionService) GetIssuesUnresolvedCount(versionID string, options *GetQueryOptions) (*IssuesUnresolvedCount, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/latest/version/%s/unresolvedIssueCount", versionID)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, nil, err
		}
		req.URL.RawQuery = q.Encode()
	}

	issuesunresolvedcount := new(IssuesUnresolvedCount)
	resp, err := s.client.Do(req, issuesunresolvedcount)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return issuesunresolvedcount, resp, nil
}
