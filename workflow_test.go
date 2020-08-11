package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)


func TestUserService_GetWorkflow(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/2/workflow"
	raw, err := ioutil.ReadFile("./mocks/workflow.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, "GET")
		testRequestURL(t, request, testAPIEndpoint)
		_, err = fmt.Fprint(writer, string(raw))
		if err != nil {
			t.Error(err.Error())
		}
	})

	workflows, _, err := testClient.User.GetWorkflow()
	if workflows == nil {
		t.Errorf("Expected GetWorkflow, got nil")
	} else {
		if len(*workflows) != 6 {
			t.Errorf("Expected GetWorkflow = 6, got " + strconv.Itoa(len(*workflows)))
		}
	}
	if err != nil {
		t.Errorf("Error given: %s", err.Error())
	}
}



