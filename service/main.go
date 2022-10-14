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
	"context"
	"errors"
	"log"
	"os"

	"github.com/yaq-cc/ezcx"
)

var (
	PORT = os.Getenv("PORT")
)

func main() {
	parent := context.Background()
	lg := log.Default()
	server := ezcx.NewServer(parent, ":"+PORT, lg)
	server.HandleCx("/inject-callerId", CxCallerIdInjectionHandler)
	server.ListenAndServe(parent)

}

// CxCallerIdInjectionHandler copies the telephony payload into the
// responses SessionInfo.Parameters.  It does not override any other parameters
func CxCallerIdInjectionHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	// Check if the payload is empty.
	payload := req.GetPayload()
	if payload == nil {
		return errors.New("ERROR: no payload found")
	}

	// Check if there was actually a telephony payload.
	telephony, ok := payload["telephony"].(map[string]any)
	if !ok {
		return errors.New("ERROR: no telephony payload found")
	}
	
	// Add the telephony payload to the session parameters.
	err := res.AddSessionParameters(telephony)
	if err != nil {
		return err
	}
	return nil
}
