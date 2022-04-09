package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestProfieldsService_GetFields(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/profields/api/2.0/fields"

	raw, err := ioutil.ReadFile("./mocks/profields_fields.json")
	if err != nil {
		t.Error(err.Error())
	}

	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEndpoint)
		_, err = fmt.Fprint(w, string(raw))
		if err != nil {
			t.Error(err.Error())
		}
	})

	fields, _, err := testClient.ProField.GetFields()
	if fields == nil {
		t.Errorf("Expected fields list is nil")
	} else {
		if len(*fields) != 8 {
			t.Errorf("Expected 8 fields but received %v ", len(*fields))
		}
	}
	if err != nil {
		t.Errorf("Error given")
	}
}

func TestProfieldsService_GetProjectField(t *testing.T) {
	setup()
	defer teardown()
	//	testAPIEndpoint := "/rest/profields/api/2.0/fields"
	projkey := "STP"
	fieldid := -2
	testAPIEndpoint := fmt.Sprintf("/rest/profields/api/2.0/values/projects/%s/fields/%v", projkey, fieldid)

	raw, err := ioutil.ReadFile("./mocks/profields_value.json")
	if err != nil {
		t.Error(err.Error())
	}

	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEndpoint)
		_, err = fmt.Fprint(w, string(raw))
		if err != nil {
			t.Error(err.Error())
		}
	})

	value, _, err := testClient.ProField.GetProjectField(projkey, fieldid)
	if value == nil {
		t.Errorf("Expected fields list is nil")
	} else {
		if value.Value.Formatted != "Scrum Test Project" {
			t.Errorf("Expected Scrum Test Project but received :%s ", value.Value.Formatted)
		}
	}
	if err != nil {
		t.Errorf("Error given: " + err.Error())

	}
}
