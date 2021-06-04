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
	"context"
	"fmt"
	"net/http"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	triggersBasePath = appsBasePath + "/%s/triggers"
)

// EventTriggersService provides access to the event triggers related functions in the Realm API.
//
// See more: https://docs.mongodb.com/realm/admin/api/v3/#event-trigger-apis
type EventTriggersService interface {
	List(context.Context, string, string) ([]Trigger, *Response, error)
}

// EventTriggersServiceOp provides an implementation of the EventTriggersService interface.
type EventTriggersServiceOp service

var _ EventTriggersService = &EventTriggersServiceOp{}

// List all triggers.
//
// See more: https://docs.mongodb.com/realm/admin/api/v3/#responses-96
func (s EventTriggersServiceOp) List(ctx context.Context, groupID, appID string) ([]Trigger, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupId", "must be set")
	}
	if appID == "" {
		return nil, nil, atlas.NewArgError("appID", "must be set")
	}
	path := fmt.Sprintf(triggersBasePath, groupID, appID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root []Trigger
	resp, err := s.Client.Do(ctx, req, &root)

	return root, resp, err
}

type Trigger struct {
	ID           string `json:"_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	FunctionID   string `json:"function_id,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	Disabled     bool   `json:"disabled"`
}
