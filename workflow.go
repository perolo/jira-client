package jira

import (
	"fmt"
)

// Workflow represents.
type Workflow struct {
	Name            string     `json:"name,omitempty"`
	Desc            string     `json:"description,omitempty"`
	Steps           int        `json:"steps,omitempty"`
	Default         bool       `json:"default,omitempty"`
}

func (s *UserService) GetWorkflow() (*[]Workflow, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/workflow")
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	workflow := new([]Workflow)
//fmt.Println("apiEndpoint: " + apiEndpoint)
	resp, err := s.client.Do(req, workflow)
	if err != nil {
		return nil, resp, NewJiraError(resp, err)
	}
	return workflow, resp, nil
}
// Not a REST api  - requires websudo disabled
func (s *UserService) SaveWorkflow(workflow string) (error) {
	apiEndpoint := fmt.Sprintf("/secure/admin/workflows/ViewWorkflowXml.jspa?workflowMode=live&workflowName=%s", workflow)
//fmt.Println("apiEndpoint: " + apiEndpoint)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return err
	}
	err = s.client.Save(req, workflow + ".xml")
	if err != nil {
		return err
	}
	return err
}

