package octopus

import (
	"testing"
)

/*
 * Integration tests
 */

// Get project group by Id (successful).
func Test_Client_GetProjectGroup_Success(test *testing.T) {
	testClientRequest(test, &ClientTest{
		APIKey: "my-test-api-key",
		Invoke: func(test *testing.T, client *Client) {
			projectGroup, err := client.GetProjectGroup("ProjectGroups-49")
			if err != nil {
				test.Fatal(err)
			}

			verifyGetProjectGroupTestResponse(test, projectGroup)
		},
		Handle: testRespondOK(getProjectGroupTestResponse),
	})
}

/*
 * Test responses.
 */

const getProjectGroupTestResponse = `
	{
		"Id": "ProjectGroups-49",
		"Name": "Platform Current",
		"Description": "Platform R1.5 and R2.0.",
		"EnvironmentIds": [
			"Environments-1",
			"Environments-2"
		],
		"RetentionPolicyId": null,
		"Links": {
			"Self": "/api/projectgroups/ProjectGroups-49",
			"Projects": "/api/projectgroups/ProjectGroups-49/projects"
		}
	}
`

func verifyGetProjectGroupTestResponse(test *testing.T, projectGroup *ProjectGroup) {
	expect := expect(test)

	expect.NotNil("ProjectGroup", projectGroup)
	expect.EqualsString("ProjectGroup.ID", "ProjectGroups-49", projectGroup.ID)
	expect.EqualsString("ProjectGroup.Name", "Platform Current", projectGroup.Name)
	expect.EqualsString("ProjectGroup.Description", "Platform R1.5 and R2.0.", projectGroup.Description)
	expect.EqualsInt("ProjectGroup.EnvironmentIDs.Length", 2, len(projectGroup.EnvironmentIDs))
}
