package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	// HTTP Basic Authentication
	authTypeBasic = 1
	// HTTP Session Authentication
	authTypeSession = 2
	// HTTP Token Authentication
	authTypeToken = 3
	// No auth
	authTypeNo = 4
)

// AuthenticationService handles authentication for the Jira instance / API.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#authentication
type AuthenticationService struct {
	client *Client

	// Authentication type
	authType int

	// Basic auth username
	username string

	// Basic auth password
	password string

	Usetoken bool
}

// Session represents a Session JSON response by the Jira API.
type Session struct {
	Self    string `json:"self,omitempty"`
	Name    string `json:"name,omitempty"`
	Session struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"session,omitempty"`
	LoginInfo struct {
		FailedLoginCount    int    `json:"failedLoginCount"`
		LoginCount          int    `json:"loginCount"`
		LastFailedLoginTime string `json:"lastFailedLoginTime"`
		PreviousLoginTime   string `json:"previousLoginTime"`
	} `json:"loginInfo"`
	Cookies []*http.Cookie
}

// SetBasicAuth sets username and password for the basic auth against the Jira instance.
//
// Reactivating to manage token login
func (s *AuthenticationService) SetBasicAuth(username, password string, usetoken bool) {
	s.username = username
	s.password = password
	s.Usetoken = usetoken
	s.authType = authTypeBasic
}
func (s *AuthenticationService) SetTokenAuth(password string, usetoken bool) {
	s.password = password
	s.Usetoken = usetoken
	s.authType = authTypeToken
}
func (s *AuthenticationService) SetAuthNo() {
	s.authType = authTypeNo
}

// Authenticated reports if the current Client has authentication details for Jira
func (s *AuthenticationService) Authenticated() bool {
	if s != nil {
		if s.authType == authTypeSession {
			return s.client.session != nil
		} else if s.authType == authTypeBasic {
			return s.username != ""
		} else if s.authType == authTypeToken {
			return s.password != ""
		}

	}
	return false
}

// LogoutWithContext logs out the current user that has been authenticated and the session in the client is destroyed.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#auth/1/session
//
// Deprecated: Use CookieAuthTransport to create base client.  Logging out is as simple as not using the
// client anymore
func (s *AuthenticationService) LogoutWithContext(ctx context.Context) error {
	if s.authType != authTypeSession || s.client.session == nil {
		return fmt.Errorf("no user is authenticated")
	}

	apiEndpoint := "rest/auth/1/session"
	req, err := s.client.NewRequestWithContext(ctx, "DELETE", apiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("creating the request to log the user out failed : %s", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("error sending the logout request: %s", err)
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf("the logout was unsuccessful with status %d", resp.StatusCode)
	}

	// If logout successful, delete session
	s.client.session = nil

	return nil

}

// Logout wraps LogoutWithContext using the background context.
//
// Deprecated: Use CookieAuthTransport to create base client.  Logging out is as simple as not using the
// client anymore
func (s *AuthenticationService) Logout() error {
	return s.LogoutWithContext(context.Background())
}

// GetCurrentUserWithContext gets the details of the current user.
//
// Jira API docs: https://docs.atlassian.com/jira/REST/latest/#auth/1/session
func (s *AuthenticationService) GetCurrentUserWithContext(ctx context.Context) (*Session, error) {
	if s == nil {
		return nil, fmt.Errorf("authentication Service is not instantiated")
	}
	if s.authType != authTypeSession || s.client.session == nil {
		return nil, fmt.Errorf("no user is authenticated yet")
	}

	apiEndpoint := "rest/auth/1/session"
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request for getting user info : %s", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("error sending request to get user info : %s", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("getting user info failed with status : %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	ret := new(Session)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read body from the response : %s", err)
	}

	err = json.Unmarshal(data, &ret)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshall received user info : %s", err)
	}

	return ret, nil
}

// GetCurrentUser wraps GetCurrentUserWithContext using the background context.
func (s *AuthenticationService) GetCurrentUser() (*Session, error) {
	return s.GetCurrentUserWithContext(context.Background())
}
