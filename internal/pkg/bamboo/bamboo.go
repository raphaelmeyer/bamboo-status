package bamboo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type BuildState uint

const (
	Success BuildState = iota
	Failure
	Building
)

type Result struct {
	Plan  string
	Build uint
	State BuildState
}

func (bs BuildState) String() string {
	switch bs {
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	case Building:
		return "Building"
	}
	panic("undefined build state")
}

type Bamboo struct {
	server   string
	username string
	password string
}

func NewBamboo(server string, username string, password string) *Bamboo {
	return &Bamboo{server, username, password}
}

func (b *Bamboo) Status(plan string) (*Result, error) {
	uri, err := b.createUri(plan)
	if err != nil {
		return nil, err
	}

	response, err := b.doRequest(uri)
	if err != nil {
		return nil, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b *Bamboo) createUri(plan string) (string, error) {
	uri, err := url.Parse(b.server)
	if err != nil {
		return "", err
	}
	uri.Path = path.Join(uri.Path, "rest/api/latest")
	uri.Path = path.Join(uri.Path, "result")
	uri.Path = path.Join(uri.Path, plan)

	query := url.Values{}
	query.Set("os_authType", "basic")
	query.Set("expand", "results[0]")
	query.Set("includeAllStates", "true")

	uri.RawQuery = query.Encode()

	return uri.String(), nil
}

func (b *Bamboo) doRequest(uri string) (string, error) {
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}

	request.SetBasicAuth(b.username, b.password)
	request.Header.Add("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf(response.Status)
	}

	return string(body), nil
}

func parseResponse(response string) (*Result, error) {
	type PlanResult struct {
		Results struct {
			Result []struct {
				State string `json:"state"`
				Build uint   `json:"number"`
				Plan  struct {
					Key string `json:"key"`
				} `json:"plan"`
			} `json:"result"`
		} `json:"results"`
	}

	var parsed PlanResult
	err := json.Unmarshal([]byte(response), &parsed)
	if err != nil {
		return nil, err
	}

	last := parsed.Results.Result[0]
	result := &Result{
		Plan:  last.Plan.Key,
		Build: last.Build,
		State: evaluate(last.State),
	}

	return result, nil
}

func evaluate(state string) BuildState {
	switch state {
	case "Successful":
		return Success

	case "Unknown":
		return Building

	default:
		return Failure

	}
}
