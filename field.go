package jira

import (
	"context"
	"github.com/google/go-querystring/query"
)

// FieldService handles fields for the Jira instance / API.
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-Field
type FieldService struct {
	client *Client
}

// Field represents a field of a Jira issue.
type Field struct {
	ID          string      `json:"id,omitempty" structs:"id,omitempty"`
	Key         string      `json:"key,omitempty" structs:"key,omitempty"`
	Name        string      `json:"name,omitempty" structs:"name,omitempty"`
	Custom      bool        `json:"custom,omitempty" structs:"custom,omitempty"`
	Navigable   bool        `json:"navigable,omitempty" structs:"navigable,omitempty"`
	Searchable  bool        `json:"searchable,omitempty" structs:"searchable,omitempty"`
	ClauseNames []string    `json:"clauseNames,omitempty" structs:"clauseNames,omitempty"`
	Schema      FieldSchema `json:"schema,omitempty" structs:"schema,omitempty"`
}

// FieldSchema represents a schema of a Jira field.
// Documentation: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-issue-fields/#api-rest-api-2-field-get
type FieldSchema struct {
	Type     string `json:"type,omitempty" structs:"type,omitempty"`
	Items    string `json:"items,omitempty" structs:"items,omitempty"`
	Custom   string `json:"custom,omitempty" structs:"custom,omitempty"`
	System   string `json:"system,omitempty" structs:"system,omitempty"`
	CustomID int64  `json:"customId,omitempty" structs:"customId,omitempty"`
}

// GetListWithContext gets all fields from Jira
//
// Jira API docs: https://developer.atlassian.com/cloud/jira/platform/rest/#api-api-2-field-get
func (s *FieldService) GetListWithContext(ctx context.Context) ([]Field, *Response, error) {
	apiEndpoint := "rest/api/2/field"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	fieldList := []Field{}
	resp, err := s.client.Do(req, &fieldList)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return fieldList, resp, nil
}

// GetList wraps GetListWithContext using the background context.
func (s *FieldService) GetList() ([]Field, *Response, error) {
	return s.GetListWithContext(context.Background())
}

type FieldOptions struct {
	// StartAt: The starting index of the returned projects. Base index: 0.
	StartAt int `url:"startAt,omitempty"`
	// MaxResults: The maximum number of projects to return per page. Default: 50.
	MaxResults int `url:"maxResults,omitempty"`
	// Expand: Expand specific sections in the returned issues
	ProjectIds      string `url:"projectIds,omitempty"`
	ScreenIds       string `url:"screenIds,omitempty"`
	Types           string `url:"types,omitempty"`
	LastValueUpdate int    `url:"lastValueUpdate,omitempty"`
}

type CustomFieldsType struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Type          string `json:"type"`
	SearcherKey   string `json:"searcherKey"`
	Self          string `json:"self"`
	NumericID     int    `json:"numericId"`
	IsLocked      bool   `json:"isLocked"`
	IsManaged     bool   `json:"isManaged"`
	IsAllProjects bool   `json:"isAllProjects"`
	ProjectsCount int    `json:"projectsCount"`
	ScreensCount  int    `json:"screensCount"`
}

type CustomFieldsResponseType struct {
	MaxResults int                `json:"maxResults"`
	StartAt    int                `json:"startAt"`
	Total      int                `json:"total"`
	IsLast     bool               `json:"isLast"`
	Values     []CustomFieldsType `json:"values"`
}

// Should be a separate file + ...
//rest/api/2/customFields

func (s *FieldService) GetAllCustomFieldsWithContext(ctx context.Context, options *FieldOptions) (*CustomFieldsResponseType, *Response, error) {
	apiEndpoint := "rest/api/2/customFields"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
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

	issue := new(CustomFieldsResponseType)
	resp, err := s.client.Do(req, issue)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return issue, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *FieldService) GetAllCustomFields(options *FieldOptions) (*CustomFieldsResponseType, *Response, error) {
	return s.GetAllCustomFieldsWithContext(context.Background(), options)
}
