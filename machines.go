package octopus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Machines represents a page of Project results.
type Machines struct {
	Items []Machine `json:"Items"`

	PagedResults
}

// Machine represents a machine that Octopus can target for deployment.
type Machine struct {
	ID                string            `json:"Id"`
	Name              string            `json:"Name"`
	Thumbprint        string            `json:"Thumbprint"`
	URI               string            `json:"Uri"`
	IsDisabled        bool              `json:"IsDisabled"`
	EnvironmentIDs    []string          `json:"EnvironmentIds"`
	Roles             []string          `json:"Roles"`
	Status            string            `json:"Status"`
	HasLatestCalamari bool              `json:"HasLatestCalamari"`
	Endpoint          Endpoint          `json:"Endpoint"`
	Links             map[string]string `json:"Links"`
}

// Endpoint represents an Octopus deployment end-point.
type Endpoint struct {
	ID                     *string                `json:"Id,omitempty"`
	CommunicationsStyle    string                 `json:"CommunicationsStyle"`
	URI                    string                 `json:"Uri"`
	Thumbprint             string                 `json:"Thumbprint"`
	TentacleVersionDetails TentacleVersionDetails `json:"TentacleVersionDetails"`
	LastModifiedOn         *string                `json:"LastModifiedOn,omitempty"`
	LastModifiedBy         *string                `json:"LastModifiedBy,omitempty"`
	Links                  map[string]string      `json:"Links"`
}

// TentacleVersionDetails represents version information for an Octopus tentacle.
type TentacleVersionDetails struct {
	Version          string `json:"Version"`
	UpgradeSuggested bool   `json:"UpgradeSuggested"`
	UpgradeRequired  bool   `json:"UpgradeRequired"`
	UpgradeLocked    bool   `json:"UpgradeLocked"`
}

// GetMachines retrieves a page of Octopus machines.
//
// skip indicates the number of results to skip over.
// Call Machines.GetSkipForNextPage() / Machines.GetSkipForPreviousPage() to get the number of items to skip for the next / previous page of results.
func (client *Client) GetMachines(skip int) (machines *Machines, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("machines?skip=%d", skip)
	request, err = client.newRequest(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		errorResponse, err = readAPIErrorResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, errorResponse.ToError("Request to retrieve all machines failed with status code %d.", statusCode)
	}

	machines = &Machines{}
	err = json.Unmarshal(responseBody, &machines)
	if err != nil {
		return
	}

	return
}

// GetMachine retrieves an Octopus machine by Id or slug.
func (client *Client) GetMachine(idOrSlug string) (machine *Machine, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("machines/%s", idOrSlug)
	request, err = client.newRequest(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		err = fmt.Errorf("Error invoking request to read machine '%s': %s", idOrSlug, err.Error())

		return
	}

	if statusCode == http.StatusNotFound {
		// Environment not found.
		return nil, nil
	}

	if statusCode != http.StatusOK {
		errorResponse, err = readAPIErrorResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, errorResponse.ToError("Request to retrieve machine '%s' failed with status code %d.", idOrSlug, statusCode)
	}

	machine = &Machine{}
	err = json.Unmarshal(responseBody, machine)
	if err != nil {
		err = fmt.Errorf("Invalid response detected when retrieving machine '%s': %s", idOrSlug, err.Error())
	}

	return
}

// GetMachineByName retrieves an Octopus machine by name from the set of all machines.
func (client Client) GetMachineByName(name string) (result *Machine, found bool, err error) {
	skip := 0
	ok := true

	for ok {
		var machinesPage *Machines
		machinesPage, err = client.GetMachines(skip)
		if err != nil {
			return
		}

		for _, machine := range machinesPage.Items {
			if machine.HasName(name) {
				result = &machine
				found = true
				return
			}
		}
		skip, ok = machinesPage.GetSkipForNextPage()
	}
	return
}

// DeleteMachine deletes a machine from the Octopus server
func (client Client) DeleteMachine(machine *Machine) (err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := machine.Links["Self"]
	request, err = client.newRequest(requestURI, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		errorResponse, err = readAPIErrorResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return err
		}

		return errorResponse.ToError("Request to delete machine '%s' failed with status code %d", machine.ID, statusCode)
	}

	return nil
}

// HasName determines whether a machine has the specified name (case-insensitive)
func (machine Machine) HasName(name string) bool {
	return strings.ToLower(machine.Name) == strings.ToLower(name)
}
