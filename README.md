# Dialogflow CX Telephony Payload to Session Parameters Injector
The provided webhook copies the `telephony` payload into the WebhookResponse's session parameters.  It does not 'over-write' any existing parameters unless there is an existing `"caller_id"` session parameter.

# As-Is Disclaimer
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License..

# Usage
## Requirements
- Make sure you have sufficient IAM permissions to: 
  1. Edit Dialogflow CX Agent definitions (make and assign webhooks)
  2. Create and Store container images in the Container Registry
  3. Define and deploy Cloud Run services.
- It is recommended to host the Cloud Run service that hosts the webhook in the same project where the Dialogflow Virtual Agent exists.  

## Instructions
Review the provided cloudbuild.yaml's Cloud Run configuration (see step id:`gcloud-run-deploy-cx-callerId-to-session`).  Once reviewed, run `gcloud builds submit` from Cloud Shell.  This will create and store the container within Google Cloud's container registry and then create the Cloud Run service that hosts the Webhook.  

Once the Cloud Run service is ready, copy the provided URL, and append the Handler's path (`/inject-callerId`) to it: 

- `https://auto-generated-spam.run.app/inject-callerId`

This is the fully qualified URL to use as the webhook URL inside of the virtual agent.  

```yaml
# REVIEW THIS!
- id: gcloud-run-deploy-cx-callerId-to-session
  waitFor: ['docker-build-push-cx-callerId-to-session']
  name: gcr.io/google.com/cloudsdktool/cloud-sdk
  entrypoint: bash
  args:
    - -c
    - |
      gcloud run deploy ${_SERVICE} \
        --project $PROJECT_ID \
        --image gcr.io/$PROJECT_ID/${_SERVICE} \
        --timeout 30s \
        --region ${_REGION} \
        --min-instances 0 \
        --max-instances 3 \
        --no-allow-unauthenticated
```

# Source Code
```go
// main.go
package main

import (
	"context"
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

func CxCallerIdInjectionHandler(res *ezcx.WebhookResponse, req *ezcx.WebhookRequest) error {
	payload := req.GetPayload()
	telephony := payload["telephony"].(map[string]any)
	err := res.AddSessionParameters(telephony)
	if err != nil {
		return err
	}
	return nil
}
```