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
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	"github.com/openlyinc/pointy"
)

func TestEventTriggers_List(t *testing.T) {
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

	expected := []EventTrigger{
		{
			ID:           "4c7498dg87d9e6526801572b",
			Name:         "name",
			Type:         "type",
			FunctionID:   "1",
			FunctionName: "name",
			Disabled:     pointy.Bool(false),
		},
	}

	if diff := deep.Equal(triggers, expected); diff != nil {
		t.Error(diff)
	}
}

func TestEventTriggers_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	appID := "5c7498dg87d9e6526801572b"
	triggerID := "4c7498dg87d9e6526801572b"

	path := fmt.Sprintf("/groups/%s/apps/%s/triggers/%s", groupID, appID, triggerID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
		  "_id": "4c7498dg87d9e6526801572b",
		  "name": "name",
		  "type": "type",
		  "function_id": "1",
		  "function_name": "name",
		  "disabled": false
		}`)
	})

	trigger, _, err := client.EventTriggers.Get(ctx, groupID, appID, triggerID)
	if err != nil {
		t.Fatalf("EventTriggers.Get returned error: %v", err)
	}

	expected := EventTrigger{
		ID:           "4c7498dg87d9e6526801572b",
		Name:         "name",
		Type:         "type",
		FunctionID:   "1",
		FunctionName: "name",
		Disabled:     pointy.Bool(false),
	}

	if diff := deep.Equal(trigger, &expected); diff != nil {
		t.Error(diff)
	}
}

func TestEventTriggers_CreateDatabase(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	appID := "5c7498dg87d9e6526801572b"

	createRequest := &EventTriggerRequest{
		Name:       "Test event trigger",
		Type:       "DATABASE",
		FunctionID: "60c2badbbd2b5c170e91292d",
		Disabled:   pointy.Bool(false),
		Config: &EventTriggerConfig{
			OperationTypes:           []string{"INSERT", "UPDATE", "REPLACE", "DELETE"},
			Database:                 "database",
			Collection:               "collection",
			ServiceID:                "60c2badbbd2b5c170e91292c",
			Match:                    `{"updateDescription.updatedFields.status": {"$exists": true}}`,
			Project:                  `{"updateDescription.updatedFields.FieldA": 1, "operationType": 1}`,
			FullDocument:             pointy.Bool(false),
			FullDocumentBeforeChange: pointy.Bool(false),
			Unordered:                pointy.Bool(false),
		},
		EventProcessors: map[string]interface{}{
			"AWS_EVENTBRIDGE": map[string]interface{}{
				"type": "AWS_EVENTBRIDGE",
				"config": map[string]interface{}{
					"account_id": "012345678901",
					"region":     "us-east-1",
				},
			},
		},
	}

	path := fmt.Sprintf("/groups/%s/apps/%s/triggers", groupID, appID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		expected := map[string]interface{}{
			"name":        "Test event trigger",
			"type":        "DATABASE",
			"function_id": "60c2badbbd2b5c170e91292d",
			"disabled":    false,
			"config": map[string]interface{}{
				"operation_types":             []interface{}{"INSERT", "UPDATE", "REPLACE", "DELETE"},
				"database":                    "database",
				"collection":                  "collection",
				"service_id":                  "60c2badbbd2b5c170e91292c",
				"match":                       `{"updateDescription.updatedFields.status": {"$exists": true}}`,
				"project":                     `{"updateDescription.updatedFields.FieldA": 1, "operationType": 1}`,
				"full_document":               false,
				"full_document_before_change": false,
				"unordered":                   false,
			},
			"event_processors": map[string]interface{}{
				"AWS_EVENTBRIDGE": map[string]interface{}{
					"type": "AWS_EVENTBRIDGE",
					"config": map[string]interface{}{
						"account_id": "012345678901",
						"region":     "us-east-1",
					},
				},
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, `
		{
		 "_id": "4c7498dg87d9e6526801572b",
		 "name": "name",
		 "type": "DATABASE",
		 "function_id": "60c2badbbd2b5c170e91292d",
		 "function_name": "name",
		 "disabled": false,
		 "config": {
		  "operation_types": [
		   "INSERT",
		   "UPDATE",
           "REPLACE",
		   "DELETE"
		  ],
		  "database": "sample_airbnb",
		  "collection": "listingsAndReviews",
		  "clusterName": "cluster name",
		  "service_id": "60c2badbbd2b5c170e91292c",
		  "match": {"updateDescription.updatedFields.status": {"$exists": true}},
		  "project": {"updateDescription.updatedFields.FieldA": 1, "operationType": 1},
		  "full_document": false,
		  "full_document_before_change": false,
		  "unordered": false
		 },
		"event_processors": {
		   "AWS_EVENTBRIDGE": {
			  "type": "AWS_EVENTBRIDGE",
			  "config": {
				 "account_id": "012345678901",
				 "region": "us-east-1"
			  }
		   }
		},
		 "last_modified": 1623793011
		}`)
	})

	updatedEventTriggersServiceOp, _, err := client.EventTriggers.Create(ctx, groupID, appID, createRequest)
	if err != nil {
		t.Fatalf("EventTriggers.Create returned error: %v", err)
	}

	expected := EventTrigger{
		ID:           "4c7498dg87d9e6526801572b",
		Name:         "name",
		Type:         "DATABASE",
		FunctionID:   "60c2badbbd2b5c170e91292d",
		FunctionName: "name",
		Disabled:     pointy.Bool(false),
		Config: EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE", "REPLACE", "DELETE"},
			Database:       "sample_airbnb",
			Collection:     "listingsAndReviews",
			ServiceID:      "60c2badbbd2b5c170e91292c",
			Match: map[string]interface{}{
				"updateDescription.updatedFields.status": map[string]interface{}{
					"$exists": true,
				},
			},
			Project: map[string]interface{}{
				"operationType":                          float64(1),
				"updateDescription.updatedFields.FieldA": float64(1),
			},
			FullDocument:             pointy.Bool(false),
			FullDocumentBeforeChange: pointy.Bool(false),
			Unordered:                pointy.Bool(false),
			ClusterName:              "cluster name",
		},
		EventProcessors: map[string]interface{}{
			"AWS_EVENTBRIDGE": map[string]interface{}{
				"type": "AWS_EVENTBRIDGE",
				"config": map[string]interface{}{
					"account_id": "012345678901",
					"region":     "us-east-1",
				},
			},
		},
		LastModified: pointy.Int64(int64(1623793011)),
	}

	if diff := deep.Equal(updatedEventTriggersServiceOp, &expected); diff != nil {
		t.Error(diff)
	}
}

func TestEventTriggers_CreateAuthentication(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	appID := "5c7498dg87d9e6526801572b"

	createRequest := &EventTriggerRequest{
		Name:       "Test event trigger",
		Type:       "AUTHENTICATION",
		FunctionID: "60c2badbbd2b5c170e91292d",
		Disabled:   pointy.Bool(false),
		Config: &EventTriggerConfig{
			OperationType: "LOGIN",
			Providers:     []string{"anon-user", "local-userpass"},
		},
		EventProcessors: map[string]interface{}{
			"AWS_EVENTBRIDGE": map[string]interface{}{
				"type": "AWS_EVENTBRIDGE",
				"config": map[string]interface{}{
					"account_id": "012345678901",
					"region":     "us-east-1",
				},
			},
		},
	}

	path := fmt.Sprintf("/groups/%s/apps/%s/triggers", groupID, appID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		expected := map[string]interface{}{
			"name":        "Test event trigger",
			"type":        "AUTHENTICATION",
			"function_id": "60c2badbbd2b5c170e91292d",
			"disabled":    false,
			"config": map[string]interface{}{
				"operation_type": "LOGIN",
				"providers":      []interface{}{"anon-user", "local-userpass"},
			},
			"event_processors": map[string]interface{}{
				"AWS_EVENTBRIDGE": map[string]interface{}{
					"type": "AWS_EVENTBRIDGE",
					"config": map[string]interface{}{
						"account_id": "012345678901",
						"region":     "us-east-1",
					},
				},
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, `{
		 "_id": "4c7498dg87d9e6526801572b",
		 "name": "name",
		 "type": "AUTHENTICATION",
		 "function_id": "60c2badbbd2b5c170e91292d",
		 "function_name": "name",
		 "disabled": false,
		 "config": {
		   "operation_type": "LOGIN",
		  "providers": [
		   "anon-user",
		   "local-userpass"
		  ],
		  "match": {},
		  "project": {}
		 },
		"event_processors": {
		   "AWS_EVENTBRIDGE": {
			  "type": "AWS_EVENTBRIDGE",
			  "config": {
				 "account_id": "012345678901",
				 "region": "us-east-1"
			  }
		   }
		},
		 "last_modified": 1623793011
		}`)
	})

	updatedEventTriggersServiceOp, _, err := client.EventTriggers.Create(ctx, groupID, appID, createRequest)
	if err != nil {
		t.Fatalf("EventTriggers.Create returned error: %v", err)
	}

	expected := EventTrigger{
		ID:           "4c7498dg87d9e6526801572b",
		Name:         "name",
		Type:         "AUTHENTICATION",
		FunctionID:   "60c2badbbd2b5c170e91292d",
		FunctionName: "name",
		Disabled:     pointy.Bool(false),
		Config: EventTriggerConfig{
			OperationType: "LOGIN",
			Providers:     []string{"anon-user", "local-userpass"},
			Match:         map[string]interface{}{},
			Project:       map[string]interface{}{},
		},
		EventProcessors: map[string]interface{}{
			"AWS_EVENTBRIDGE": map[string]interface{}{
				"type": "AWS_EVENTBRIDGE",
				"config": map[string]interface{}{
					"account_id": "012345678901",
					"region":     "us-east-1",
				},
			},
		},
		LastModified: pointy.Int64(int64(1623793011)),
	}

	if diff := deep.Equal(updatedEventTriggersServiceOp, &expected); diff != nil {
		t.Error(diff)
	}
}

func TestEventTriggers_CreateScheduled(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	appID := "5c7498dg87d9e6526801572b"

	createRequest := &EventTriggerRequest{
		Name:       "Test event trigger",
		Type:       "SCHEDULED",
		FunctionID: "60c2badbbd2b5c170e91292d",
		Disabled:   pointy.Bool(false),
		Config: &EventTriggerConfig{
			Schedule: "* * * * *",
		},
		EventProcessors: map[string]interface{}{
			"AWS_EVENTBRIDGE": map[string]interface{}{
				"type": "AWS_EVENTBRIDGE",
				"config": map[string]interface{}{
					"account_id": "012345678901",
					"region":     "us-east-1",
				},
			},
		},
	}

	path := fmt.Sprintf("/groups/%s/apps/%s/triggers", groupID, appID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		expected := map[string]interface{}{
			"name":        "Test event trigger",
			"type":        "SCHEDULED",
			"function_id": "60c2badbbd2b5c170e91292d",
			"disabled":    false,
			"config": map[string]interface{}{
				"schedule": "* * * * *",
			},
			"event_processors": map[string]interface{}{
				"AWS_EVENTBRIDGE": map[string]interface{}{
					"type": "AWS_EVENTBRIDGE",
					"config": map[string]interface{}{
						"account_id": "012345678901",
						"region":     "us-east-1",
					},
				},
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, `{
		 "_id": "4c7498dg87d9e6526801572b",
		 "name": "name",
		 "type": "SCHEDULED",
		 "function_id": "60c2badbbd2b5c170e91292d",
		 "function_name": "name",
		 "disabled": false,
		 "config": {
		   "schedule": "* * * * *",
		   "schedule_type": "ADVANCED"
		 },
		"event_processors": {
		   "AWS_EVENTBRIDGE": {
			  "type": "AWS_EVENTBRIDGE",
			  "config": {
				 "account_id": "012345678901",
				 "region": "us-east-1"
			  }
		   }
		},
		 "last_modified": 1623793011
		}`)
	})

	updatedEventTriggersServiceOp, _, err := client.EventTriggers.Create(ctx, groupID, appID, createRequest)
	if err != nil {
		t.Fatalf("EventTriggers.Create returned error: %v", err)
	}

	expected := EventTrigger{
		ID:           "4c7498dg87d9e6526801572b",
		Name:         "name",
		Type:         "SCHEDULED",
		FunctionID:   "60c2badbbd2b5c170e91292d",
		FunctionName: "name",
		Disabled:     pointy.Bool(false),
		Config: EventTriggerConfig{
			Schedule:     "* * * * *",
			ScheduleType: "ADVANCED",
		},
		EventProcessors: map[string]interface{}{
			"AWS_EVENTBRIDGE": map[string]interface{}{
				"type": "AWS_EVENTBRIDGE",
				"config": map[string]interface{}{
					"account_id": "012345678901",
					"region":     "us-east-1",
				},
			},
		},
		LastModified: pointy.Int64(int64(1623793011)),
	}

	if diff := deep.Equal(updatedEventTriggersServiceOp, &expected); diff != nil {
		t.Error(diff)
	}
}

func TestEventTriggers_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	appID := "5c7498dg87d9e6526801572b"
	triggerID := "4c7498dg87d9e6526801572b"

	updateRequest := &EventTriggerRequest{
		Name:       "Test event trigger update",
		Type:       "Test event trigger update",
		FunctionID: "1",
		Disabled:   pointy.Bool(false),
		Config: &EventTriggerConfig{
			OperationTypes: []string{"INSERT", "UPDATE"},
			OperationType:  "CREATE",
			Providers:      []string{"anon-user", "local-userpass"},
			Database:       "database2",
			Collection:     "collection2",
			ServiceID:      "3",
			Match:          `{"updateDescription.updatedFields.status": {"$exists": true}}`,
			Project:        `{"updateDescription.updatedFields.FieldA": 1, "operationType": 1}`,
			FullDocument:   pointy.Bool(false),
			Schedule:       "weekday",
		},
		EventProcessors: map[string]interface{}{
			"AWS_EVENTBRIDGE": map[string]interface{}{
				"type": "AWS_EVENTBRIDGE",
				"config": map[string]interface{}{
					"account_id": "012345678901",
					"region":     "us-east-1",
				},
			},
		},
	}

	path := fmt.Sprintf("/groups/%s/apps/%s/triggers/%s", groupID, appID, triggerID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		expected := map[string]interface{}{
			"name":        "Test event trigger update",
			"type":        "Test event trigger update",
			"function_id": "1",
			"disabled":    false,
			"config": map[string]interface{}{
				"operation_types": []interface{}{"INSERT", "UPDATE"},
				"operation_type":  "CREATE",
				"providers":       []interface{}{"anon-user", "local-userpass"},
				"database":        "database2",
				"collection":      "collection2",
				"service_id":      "3",
				"match":           `{"updateDescription.updatedFields.status": {"$exists": true}}`,
				"project":         `{"updateDescription.updatedFields.FieldA": 1, "operationType": 1}`,
				"full_document":   false,
				"schedule":        "weekday",
			},
			"event_processors": map[string]interface{}{
				"AWS_EVENTBRIDGE": map[string]interface{}{
					"type": "AWS_EVENTBRIDGE",
					"config": map[string]interface{}{
						"account_id": "012345678901",
						"region":     "us-east-1",
					},
				},
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, `{
		  "_id": "4c7498dg87d9e6526801572b",
		  "name": "name",
		  "type": "type",
		  "function_id": "1",
		  "function_name": "name",
		  "disabled": false
		}`)
	})

	updatedEventTriggersServiceOp, _, err := client.EventTriggers.Update(ctx, groupID, appID, triggerID, updateRequest)
	if err != nil {
		t.Fatalf("EventTriggers.Update returned error: %v", err)
	}

	expected := EventTrigger{
		ID:           "4c7498dg87d9e6526801572b",
		Name:         "name",
		Type:         "type",
		FunctionID:   "1",
		FunctionName: "name",
		Disabled:     pointy.Bool(false),
	}

	if diff := deep.Equal(updatedEventTriggersServiceOp, &expected); diff != nil {
		t.Error(diff)
	}
}

func TestEventTriggers_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	appID := "5c7498dg87d9e6526801572b"
	triggerID := "4c7498dg87d9e6526801572b"

	path := fmt.Sprintf("/groups/%s/apps/%s/triggers/%s", groupID, appID, triggerID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.EventTriggers.Delete(ctx, groupID, appID, triggerID)
	if err != nil {
		t.Fatalf("EventTriggers.Delete returned error: %v", err)
	}
}
