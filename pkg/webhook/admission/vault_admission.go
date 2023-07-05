package admission

import (
	"encoding/json"
	"fmt"
	vaultv1alpha1 "github.com/bank-vaults/vault-operator/pkg/apis/vault/v1alpha1"
	"github.com/bank-vaults/vault-operator/pkg/webhook/mutation"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	"net/http"
)

// serveMutateVault returns an admission review with Vault mutations as a json patch
// in the review response
func serveMutateVault(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("uri", r.RequestURI)
	logger.Debug("received mutation request")

	in, err := parseRequest(*r)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out, err := mutateVaultReview(logger, in.Request)
	if err != nil {
		e := fmt.Sprintf("could not generate admission response: %v", err)
		logger.Error(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jout, err := json.Marshal(out)
	if err != nil {
		e := fmt.Sprintf("could not parse admission response: %v", err)
		logger.Error(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	logger.Debug("sending response")
	logger.Debugf("%s", jout)
	_, _ = fmt.Fprintf(w, "%s", jout)
}

// mutateVaultReview takes an admission request and mutates the Vault within,
// it returns an admission review with mutations as a json patch (if any)
func mutateVaultReview(logger *logrus.Entry, req *admissionv1.AdmissionRequest) (*admissionv1.AdmissionReview, error) {
	vault, err := getVault(req)
	if err != nil {
		e := fmt.Sprintf("could not parse vault in admission review request: %v", err)
		return reviewResponse(req.UID, false, http.StatusBadRequest, e), err
	}

	patch, err := mutation.MutateVaultPatch(vault)
	if err != nil {
		e := fmt.Sprintf("could not mutate vault: %v", err)
		return reviewResponse(req.UID, false, http.StatusBadRequest, e), err
	}

	return patchReviewResponse(req.UID, patch)
}

// getVault extracts a Vault from an admission request
func getVault(req *admissionv1.AdmissionRequest) (*vaultv1alpha1.Vault, error) {
	if req.Kind.Kind != "Vault" &&
		req.Kind.Group != vaultv1alpha1.SchemeGroupVersion.Group {
		return nil, fmt.Errorf("only %s.Vault are supported here", vaultv1alpha1.SchemeGroupVersion)
	}

	p := vaultv1alpha1.Vault{}
	if err := json.Unmarshal(req.Object.Raw, &p); err != nil {
		return nil, err
	}

	return &p, nil
}
