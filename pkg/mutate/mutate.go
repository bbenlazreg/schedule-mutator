// Package mutate dealing with AdmissionReview requests and responses
package mutate

import (
	"encoding/json"
	"fmt"
	"log"

	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mutate function
func Mutate(body []byte) ([]byte, error) {

	log.Printf("recieved: %s\n", string(body))


	// unmarshal request into AdmissionReview struct
	admReview := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request with errors %s", err)
	}

	var err error
	var pod *corev1.Pod

	responseBody := []byte{}
	ar := admReview.Request
	resp := v1beta1.AdmissionResponse{}

	if ar != nil {

		// get the Pod object and unmarshal
		if err := json.Unmarshal(ar.Object.Raw, &pod); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pod json %v", err)
		}
		// Copy UID to response ad set patch type
		resp.Allowed = true
		resp.UID = ar.UID
		pT := v1beta1.PatchTypeJSONPatch
		resp.PatchType = &pT

		//audit annotations for debug
		resp.AuditAnnotations = map[string]string{
			"gl-apps-schedule": "done",
		}


		// Patches
		type patchMapType map[string]interface{}
		podPatchMap := make([]patchMapType, 0, 0)
        // patch1: toleration patch
		patch1 := map[string]interface{}{
			"op":    "add",
			"path":  "/spec/tolerations/-",
			"value": map[string]string{"key": "AppsOnly", "operator": "Equal", "value": "true"},
		}
	

		podPatchMap = append(podPatchMap, patch1)
		
		// patch2: safe-to-evict patch (should this apply to all pods ? TBA )
		patch2 := map[string]interface{}{
			"op":    "add",
			"path":  "/metadata/annotations/cluster-autoscaler.kubernetes.io~1safe-to-evict",
			"value": "true",
		}			

		podPatchMap = append(podPatchMap, patch2)

        // patch3: nodeSelector patch
		patch3 := map[string]interface{}{
			"op":    "add",
			"path":  "/spec/nodeSelector",
			"value": map[string]string{"AppsOnly": "true"},
		}			

		podPatchMap = append(podPatchMap, patch3)	
		
		resp.Patch, err = json.Marshal(podPatchMap)

		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp
		// ready to send response to the api server convert back to json
		responseBody, err = json.Marshal(admReview)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("response: %s\n", string(responseBody))


	return responseBody, nil
}
