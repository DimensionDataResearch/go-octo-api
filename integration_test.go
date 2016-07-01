package octopus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
 * Integration test support
 */

// Configuration for a Client integration test.
type ClientTest struct {
	APIKey      string
	ContentType string
	Request     ClientTestRequester
	Respond     ClientTestResponder
}

// Add default configuration to the client test configuration (if not already specified).
func (clientTest *ClientTest) AddDefaultConfiguration() *ClientTest {
	if clientTest.APIKey == "" {
		clientTest.APIKey = "my-test-api-key"
	}

	if clientTest.ContentType == "" {
		clientTest.ContentType = "application/json"
	}

	return clientTest
}

// A function that invokes the request(s) for an integration test.
type ClientTestRequester func(test *testing.T, client *Client)

// A function that handles requests and generates responses for an integration test.
type ClientTestResponder func(test *testing.T, request *http.Request) (statusCode int, responseBody string)

// Respond with HTTP OK (200) and the specified response body.
func testRespondOK(responseBody string) ClientTestResponder {
	return testRespond(http.StatusOK, responseBody)
}

// Respond with HTTP CREATED (201) and the specified response body.
func testRespondCreated(responseBody string) ClientTestResponder {
	return testRespond(http.StatusCreated, responseBody)
}

func testRespond(statusCode int, responseBody string) ClientTestResponder {
	return func(test *testing.T, request *http.Request) (int, string) {
		return statusCode, responseBody
	}
}

func testClientRequest(test *testing.T, clientTest *ClientTest) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		expect := expect(test)
		expect.headerValue(HeaderNameOctopusAPIKey, clientTest.APIKey, request)

		statusCode, response := clientTest.Respond(test, request)

		writer.Header().Set("Content-Type", clientTest.ContentType)
		writer.WriteHeader(statusCode)

		fmt.Fprintln(writer, response)
	}))
	defer testServer.Close()

	client, err := NewClientWithAPIKey(testServer.URL, clientTest.APIKey)
	if err != nil {
		test.Fatal(err)
	}

	clientTest.Request(test, client)
}
