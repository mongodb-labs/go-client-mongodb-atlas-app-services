// Copyright 2021 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package realm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
)

func TestAppsServiceOp_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"

	path := fmt.Sprintf("/groups/%s/apps", groupID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{
		"_id": "1",
		"client_app_id": "2",
		"name": "name",
		"location": "location",
		"deployment_model": "model",
		"domain_id": "3",
		"group_id": "6c7498dg87d9e6526801572b"
	  }]`)
	})

	snapshots, _, err := client.Apps.List(ctx, groupID, nil)
	if err != nil {
		t.Fatalf("Apps.List returned error: %v", err)
	}

	expected := []Application{
		{
			ID:              "1",
			ClientAppID:     "2",
			Name:            "name",
			Location:        "location",
			DeploymentModel: "model",
			DomainID:        "3",
			GroupID:         "6c7498dg87d9e6526801572b",
		},
	}

	if diff := deep.Equal(snapshots, expected); diff != nil {
		t.Error(diff)
	}
}
