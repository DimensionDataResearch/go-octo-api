package octopus

import (
	"strings"
	"testing"
)

// To handle URIs coming from the Links fields
func Test_FullPathURL(t *testing.T) {
	url := "/api/machines/Machines-42"
	if actual, err := normalizeURI(url); err != nil {
		t.Error("Got error: ", err)
	} else if strings.Compare(url, actual.String()) != 0 {
		t.Errorf("Expected \"%s\", got \"%s\"", url, actual)
	}
}

func Test_PartialPathURL(t *testing.T) {
	testCases := []struct {
		url      string
		expected string
	}{
		{"machines?skip=10", "/api/machines?skip=10"},
		{"variables/variableset-Projects-501", "/api/variables/variableset-Projects-501"},
	}
	for _, testCase := range testCases {
		if actual, err := normalizeURI(testCase.url); err != nil {
			t.Error("Got error: ", err)
		} else if strings.Compare(testCase.expected, actual.String()) != 0 {
			t.Errorf("Expected \"%s\", got \"%s\"", testCase.expected, actual)
		}
	}
}
