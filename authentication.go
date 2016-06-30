package octopus

import (
	"net/http"
)

const (
	// HeaderNameOctopusAPIKey is the name of the "X-Octopus-ApiKey" header used for authenticating to Octopus Deploy.
	HeaderNameOctopusAPIKey = "X-Octopus-ApiKey"
)

type apiKeyAuthenticator struct {
	apiKey    string
	transport *http.Transport
}

func (authenticator *apiKeyAuthenticator) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set(HeaderNameOctopusAPIKey, authenticator.apiKey)

	return authenticator.transport.RoundTrip(request)
}
