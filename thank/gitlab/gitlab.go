package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/adamliesko/go-thanks/discover"
)

const apiEndpoint = "https://gitlab.com/api/v4/"

// Thanker is a light Gitlab API HTTP client capable of starring a repository.
type Thanker struct {
	cl       *http.Client
	apiToken string
}

// New creates a new thanker for Gitlab with the apiToken set.
func New(token string) Thanker {
	return Thanker{
		cl:       &http.Client{Timeout: 10 * time.Second},
		apiToken: token,
	}
}

// Auth checks whether thanker's apiToken is a valid one for Gitlab.
func (t Thanker) Auth() error {
	uri := fmt.Sprintf("%s/user/%s", apiEndpoint, t.authTokenParams())

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}

	resp, err := t.cl.Do(req)
	if err != nil {
		return fmt.Errorf("gitlab authentication failed with an errror: %v", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return fmt.Errorf("gitlab authentication failed with HTTP status code: %d", code)
	}

	return nil
}

// Thank stars a repository as a user with thanker's apiToken on Gitlab.
func (t Thanker) Thank(r discover.Repository) error {
	urlString := fmt.Sprintf("%s/projects/%s%%2F%s/star%s", apiEndpoint, r.Owner, r.Name, t.authTokenParams())
	uri, err := url.Parse(urlString)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", uri.String(), nil)
	if err != nil {
		return err
	}
	// Note that you'll need to set Content-Length to zero when calling out to this endpoint.
	req.ContentLength = 0

	resp, err := t.cl.Do(req)
	if err != nil {
		return fmt.Errorf("gitlab starring repo %s failed with an errror: %v", r.URL, err)
	}
	defer resp.Body.Close()
	if code := resp.StatusCode; code != http.StatusCreated && code != http.StatusNotModified {
		return fmt.Errorf("gitlab starring repo %s failed with HTTP status code %d", r.URL, code)
	}

	return nil
}

// CanThank reports, whether Gitlab thanker is capable of thanking to repository r, by checking package prefix.
func (t Thanker) CanThank(r discover.Repository) bool {
	return strings.HasPrefix(r.URL, "gitlab.com")
}

func (t Thanker) authTokenParams() string {
	return fmt.Sprintf("?private_token=%s", t.apiToken)
}
