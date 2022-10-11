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
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

// Synthetic unit test for Handler via HTTP Path
func TestServer(t *testing.T) {
	var buf bytes.Buffer
	req, err := sampleWebhookRequest()
	if err != nil {
		t.Fatal(err)
	}
	err = req.WriteRequest(&buf)
	if err != nil {
		t.Fatal(err)
	}
	httpReq := httptest.NewRequest(http.MethodPost, "/inject-callerId", &buf)
	w := httptest.NewRecorder()
	hf := ezcx.HandlerFunc(CxCallerIdInjectionHandler)
	hf.ServeHTTP(w, httpReq)
	res := w.Result()
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}

func TestProtoUnmarshalTest(t *testing.T) {
	var whreq ezcx.WebhookRequest
	rd := strings.NewReader(sampleWebhookString)
	err := whreq.ReadReader(rd)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(whreq.SessionInfo)
}

func TestServerFromRawSample(t *testing.T) {
	r := strings.NewReader(sampleWebhookString)
	httpReq := httptest.NewRequest(http.MethodPost, "/inject-callerId", r)
	w := httptest.NewRecorder()
	hf := ezcx.HandlerFunc(CxCallerIdInjectionHandler)
	hf.ServeHTTP(w, httpReq)
	res := w.Result()
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
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

var sampleWebhookString = `{
	"payload": {
		"telephony": {
			"caller_id": "+12025553711"
		}
	},
	"detectIntentResponseId": "d5c7e9b7-9c9f-4056-a553-06bd1e9b006c",
	"intentInfo": {
	  "lastMatchedIntent": "projects/vocal-etching-343420/locations/global/agents/bde16cf9-d795-4e40-828b-8019c2c9dbe5/intents/00000000-0000-0000-0000-000000000000",
	  "displayName": "Default Welcome Intent",
	  "confidence": 1.0
	},
	"pageInfo": {
	  "currentPage": "projects/vocal-etching-343420/locations/global/agents/bde16cf9-d795-4e40-828b-8019c2c9dbe5/flows/00000000-0000-0000-0000-000000000000/pages/cddd1fc1-89d7-4f0e-af8f-68402883ce0b",
	  "formInfo": {
	  },
	  "displayName": "page.get-caller-id"
	},
	"sessionInfo": {
	  "session": "projects/vocal-etching-343420/locations/global/agents/bde16cf9-d795-4e40-828b-8019c2c9dbe5/sessions/8adaaf-854-94a-485-da05fdae4"
	},
	"fulfillmentInfo": {
	  "tag": "testing-webhook"
	},
	"messages": [{
	  "text": {
		"text": ["Hi there. "],
		"redactedText": ["Hi there. "]
	  },
	  "responseType": "HANDLER_PROMPT",
	  "source": "VIRTUAL_AGENT"
	}, {
	  "text": {
		"text": ["My only role in life is to get your phone number and read it to you. "],
		"redactedText": ["My only role in life is to get your phone number and read it to you. "]
	  },
	  "responseType": "ENTRY_PROMPT",
	  "source": "VIRTUAL_AGENT"
	}],
	"text": "hi",
	"languageCode": "en"
  }`
