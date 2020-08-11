package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestVersionService_Get_Success(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/version/10002", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/rest/api/2/version/10002")

		_, err := fmt.Fprint(w, `{
			"self": "http://www.example.com/jira/rest/api/2/version/10002",
			"id": "10002",
			"description": "An excellent version",
			"name": "New Version 1",
			"archived": false,
			"released": true,
			"releaseDate": "2010-07-06",
			"overdue": true,
			"userReleaseDate": "6/Jul/2010",
			"startDate" : "2010-07-01",
			"projectId": 10000
		}`)
		if err != nil {
			t.Errorf("Error given: %s", err)
		}
	})

	version, _, err := testClient.Version.Get(10002)
	if version == nil {
		t.Error("Expected version. Issue is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestVersionService_Create(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/version", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/rest/api/2/version")

		w.WriteHeader(http.StatusCreated)
		_, err := fmt.Fprint(w, `{
			"description": "An excellent version",
			"name": "New Version 1",
			"archived": false,
			"released": true,
			"releaseDate": "2010-07-06",
			"userReleaseDate": "6/Jul/2010",
			"project": "PXA",
			"projectId": 10000
		  }`)
		if err != nil {
			t.Errorf("Error given: %s", err)
		}

	})

	v := &Version{
		Name:            "New Version 1",
		Description:     "An excellent version",
		ProjectID:       10000,
		Released:        true,
		ReleaseDate:     "2010-07-06",
		UserReleaseDate: "6/Jul/2010",
		StartDate:       "2018-07-01",
	}

	version, _, err := testClient.Version.Create(v)
	if version == nil {
		t.Error("Expected version. Version is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestServiceService_Update(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/rest/api/2/version/10002", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testRequestURL(t, r, "/rest/api/2/version/10002")
		_, err := fmt.Fprint(w, `{
			"description": "An excellent updated version",
			"name": "New Updated Version 1",
			"archived": false,
			"released": true,
			"releaseDate": "2010-07-06",
			"userReleaseDate": "6/Jul/2010",
			"startDate" : "2010-07-01",
			"project": "PXA",
			"projectId": 10000
		  }`)
		if err != nil {
			t.Errorf("Error given: %s", err)
		}

	})

	v := &Version{
		ID:          "10002",
		Name:        "New Updated Version 1",
		Description: "An excellent updated version",
	}

	version, _, err := testClient.Version.Update(v)
	if version == nil {
		t.Error("Expected version. Version is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestServiceService_GetRelatedIssueCounts(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/latest/version/12201/relatedIssueCounts"
	raw, err := ioutil.ReadFile("./mocks/version.json")
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

	issuecount, _, err := testClient.Version.GetRelatedIssueCounts("12201", nil)
	if issuecount == nil {
		t.Errorf("Expected IssuesFixedCount, got nil")
	} else {
		if issuecount.IssuesFixedCount != 21 {
			t.Errorf("Expected IssuesFixedCount = 21, got " + strconv.Itoa(issuecount.IssuesFixedCount))
		}
	}
	if err != nil {
		t.Errorf("Error given: %s", err.Error())
	}
}

func TestServiceService_GetIssuesUnresolvedCount(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/rest/api/latest/version/12201/unresolvedIssueCount"
	raw, err := ioutil.ReadFile("./mocks/version_unresolved.json")
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

	issuecount, _, err := testClient.Version.GetIssuesUnresolvedCount("12201", nil)
	if issuecount == nil {
		t.Errorf("Expected IssuesFixedCount, got nil")
	} else {
		if issuecount.IssuesUnresolvedCount != 15 {
			t.Errorf("Expected IssuesUnresolvedCount = 15, got " + strconv.Itoa(issuecount.IssuesUnresolvedCount))
		}
	}
	if err != nil {
		t.Errorf("Error given: %s", err.Error())
	}
}



