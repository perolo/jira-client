package jira

import "fmt"

// ProfieldService handles projects for the JIRA instance / API.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/latest/#api/2/project
type ProfieldService struct {
	client *Client
}

type ProFieldsList []struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	IsSystem    bool   `json:"isSystem"`
	CustomIds   bool   `json:"customIds,omitempty"`
}

type ProFieldsValue struct {
	Field struct {
		ID          int    `json:"id"`
		ParentID    int    `json:"parentId"`
		Description string `json:"description"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		CustomIds   bool   `json:"customIds"`
		IsSystem    bool   `json:"isSystem"`
	} `json:"field"`
	Value struct {
		Value struct {
			ID          int    `json:"id"`
			ParentID    int    `json:"parentId"`
			ParentValue string `json:"parentValue"`
			Text        string `json:"text"`
			Value       string `json:"value"`
		} `json:"value"`
		Formatted string `json:"formatted"`
	} `json:"value"`
}

func (s *ProfieldService) GetFields() (*ProFieldsList, *Response, error) {
	apiEndpoint := "rest/profields/api/2.0/fields"
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	profieldslist := new(ProFieldsList)
	resp, err := s.client.Do(req, profieldslist)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return profieldslist, resp, nil
}

func (s *ProfieldService) GetProjectField(projkey string, fieldid int) (*ProFieldsValue, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/profields/api/2.0/values/projects/%s/fields/%v", projkey, fieldid)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	profieldsvalue := new(ProFieldsValue)
	resp, err := s.client.Do(req, profieldsvalue)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return profieldsvalue, resp, nil
}
