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

func TestEventTriggersServiceOp_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	appID := "5c7498dg87d9e6526801572b"

	path := fmt.Sprintf("/groups/%s/apps/%s/triggers", groupID, appID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{
		  "_id": "4c7498dg87d9e6526801572b",
		  "name": "name",
		  "type": "type",
		  "function_id": "1",
		  "function_name": "name",
		  "disabled": false
		}]`)
	})

	triggers, _, err := client.EventTriggers.List(ctx, groupID, appID)
	if err != nil {
		t.Fatalf("EventTriggers.List returned error: %v", err)
	}

	expected := []Trigger{
		{
			ID:           "4c7498dg87d9e6526801572b",
			Name:         "name",
			Type:         "type",
			FunctionID:   "1",
			FunctionName: "name",
			Disabled:     false,
		},
	}

	if diff := deep.Equal(triggers, expected); diff != nil {
		t.Error(diff)
	}
}
