package admission

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	klusterv1alpha1 "github.com/viveksinghggits/kluster/pkg/apis/viveksingh.dev/v1alpha1"
	"gomodules.xyz/jsonpatch/v2"
	admv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"

	do "github.com/viveksinghggits/valkontroller/pkg/digitalocean"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func ServeKlusterMutation(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responsewriters.InternalError(w, r, err)
		fmt.Printf("Error %s, reading the body", err.Error())
	}

	// read body to get an instance of admissionreivew object
	// get the gvk for admission review
	gvk := admv1beta1.SchemeGroupVersion.WithKind("AdmissionReview")
	// var of type admission reiveew
	var admissionReview admv1beta1.AdmissionReview
	_, _, err = codecs.UniversalDeserializer().Decode(body, &gvk, &admissionReview)
	if err != nil {
		fmt.Printf("Error %s, converting req body to admission review type", err.Error())

	}

	// get kluster spec from admission review object
	gvkKluster := klusterv1alpha1.SchemeGroupVersion.WithKind("Kluster")
	var k klusterv1alpha1.Kluster
	_, _, err = codecs.UniversalDeserializer().Decode(admissionReview.Request.Object.Raw, &gvkKluster, &k)
	if err != nil {
		fmt.Printf("Error %s, while getting kluster type from admission review", err.Error())
	}

	// this is k
	// apiVersion: viveksingh.dev/v1alpha1
	// kind: Kluster
	// metadata:
	// name: kluster-1
	// spec:
	// name: kluster-1
	// region: "nyc1"
	// tokenSecret: "default/dosecret"
	// nodePools:
	// 	- count: 3
	// 	name: "dummy-nodepool"
	// 	size: "sizes-2vcpu-2gb"

	newKluster := k.DeepCopy()
	if newKluster.Spec.Version == "" {
		newKluster.Spec.Version = do.LatestKubeVersion(newKluster.Spec)
	}
	// this is newKluster
	// apiVersion: viveksingh.dev/v1alpha1
	// kind: Kluster
	// metadata:
	// name: kluster-1
	// spec:
	// name: kluster-1
	// region: "nyc1"
	// version: "1.21.3-do.0"
	// tokenSecret: "default/dosecret"
	// nodePools:
	// 	- count: 3
	// 	name: "dummy-nodepool"
	// 	size: "sizes-2vcpu-2gb"

	jsonKluster, err := json.Marshal(newKluster)
	if err != nil {
		fmt.Printf("Errro %s, converting new kluster resource to json", err.Error())
		// return errro, write errro to response
	}

	ops, err := jsonpatch.CreatePatch(admissionReview.Request.Object.Raw, jsonKluster)
	if err != nil {
		fmt.Printf("errro %s, creating patch", err.Error())
		// handle error
	}

	patch, err := json.Marshal(ops)
	if err != nil {
		fmt.Printf("error %s converting operations to slice byte")
	}

	fmt.Printf("patch that we have is %s", patch)

	jsonPatchType := admv1beta1.PatchTypeJSONPatch
	response := admv1beta1.AdmissionResponse{
		UID:       admissionReview.Request.UID,
		Allowed:   true,
		PatchType: &jsonPatchType,
		Patch:     patch,
	}

	admissionReview.Response = &response
	fmt.Printf("respoknse that we are trying to return is %+v\n", response)
	res, err := json.Marshal(admissionReview)
	if err != nil {
		fmt.Printf("error %s, while converting response to byte slice", err.Error())
	}

	_, err = w.Write(res)
	if err != nil {
		fmt.Printf("error %s, writing respnse to responsewriter", err.Error())
	}
}
