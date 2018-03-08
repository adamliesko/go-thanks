package github

import (
	"fmt"
	"net/http"
	"strings"

	"net/url"

	"time"

	"github.com/adamliesko/go-thanks/discover"
)

const apiEndpoint = "https://api.github.com"

// Thanker ...
type Thanker struct {
	apiToken string
}

// New ..
func New(token string) Thanker {
	return Thanker{apiToken: token}
}

// Auth ..
func (t Thanker) Auth() error {
	uri := fmt.Sprintf("%s/%s", apiEndpoint, t.authTokenParams())

	resp, err := http.Get(uri)
	if err != nil {
		return fmt.Errorf("github authentication failed with an errror: %v", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return fmt.Errorf("github authentication failed with HTTP status code: %d", code)
	}

	return nil
}

// Thank ..
func (t Thanker) Thank(r discover.Repository) error {
	urlString := fmt.Sprintf("%s/user/starred/%s/%s%s", apiEndpoint, r.Owner, r.Name, t.authTokenParams())
	uri, err := url.Parse(urlString)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("PUT", uri.String(), nil)
	if err != nil {
		return err
	}
	// Note that you'll need to set Content-Length to zero when calling out to this endpoint.
	req.ContentLength = 0

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("github starring repo %s failed with an errror: %v", r.URL, err)
	}
	defer resp.Body.Close()
	if code := resp.StatusCode; code != http.StatusNoContent {
		// TODO repo doesn't exist and bad token both might be 404s - deal wit it?
		return fmt.Errorf("github starring repo %s failed with HTTP status code %d", r.URL, code)
	}

	return nil
}

// CanThank ...
func (t Thanker) CanThank(r discover.Repository) (bool, error) {
	return strings.HasPrefix(r.URL, "github.com"), nil
}

func (t Thanker) authTokenParams() string {
	return fmt.Sprintf("?access_token=%s", t.apiToken)
}
