// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"testing"

	"github.com/yaq-cc/ezcx"
)

// Synthetic unit test for CxCallerIdInjectionHandler
func TestCxCallerIdInjectionHandler(t *testing.T) {
	req, err := sampleWebhookRequest()
	if err != nil {
		t.Fatal(err)
	}
	res := req.InitializeResponse()
	err = CxCallerIdInjectionHandler(res, req)
	if err != nil {
		t.Fatal(err)
	}
	res.WriteResponse(os.Stdout)
}

func sampleWebhookRequest() (*ezcx.WebhookRequest, error) {
	payload := make(map[string]any)
	payload["telephony"] = make(map[string]any)
	payload["telephony"].(map[string]any)["caller_id"] = "+12025553711"
	req, err := ezcx.NewTestingWebhookRequest(nil, payload, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
