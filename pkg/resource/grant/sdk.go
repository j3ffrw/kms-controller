// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package grant

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/kms"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/kms-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.KMS{}
	_ = &svcapitypes.Grant{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadManyInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newListRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.ListGrantsResponse
	resp, err = rm.sdkapi.ListGrantsWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_MANY", "ListGrants", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "UNKNOWN" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	found := false
	for _, elem := range resp.Grants {
		if elem.Constraints != nil {
			f0 := &svcapitypes.GrantConstraints{}
			if elem.Constraints.EncryptionContextEquals != nil {
				f0f0 := map[string]*string{}
				for f0f0key, f0f0valiter := range elem.Constraints.EncryptionContextEquals {
					var f0f0val string
					f0f0val = *f0f0valiter
					f0f0[f0f0key] = &f0f0val
				}
				f0.EncryptionContextEquals = f0f0
			}
			if elem.Constraints.EncryptionContextSubset != nil {
				f0f1 := map[string]*string{}
				for f0f1key, f0f1valiter := range elem.Constraints.EncryptionContextSubset {
					var f0f1val string
					f0f1val = *f0f1valiter
					f0f1[f0f1key] = &f0f1val
				}
				f0.EncryptionContextSubset = f0f1
			}
			ko.Spec.Constraints = f0
		} else {
			ko.Spec.Constraints = nil
		}
		if elem.GrantId != nil {
			ko.Status.GrantID = elem.GrantId
		} else {
			ko.Status.GrantID = nil
		}
		if elem.GranteePrincipal != nil {
			ko.Spec.GranteePrincipal = elem.GranteePrincipal
		} else {
			ko.Spec.GranteePrincipal = nil
		}
		if elem.KeyId != nil {
			ko.Spec.KeyID = elem.KeyId
		} else {
			ko.Spec.KeyID = nil
		}
		if elem.Name != nil {
			ko.Spec.Name = elem.Name
		} else {
			ko.Spec.Name = nil
		}
		if elem.Operations != nil {
			f7 := []*string{}
			for _, f7iter := range elem.Operations {
				var f7elem string
				f7elem = *f7iter
				f7 = append(f7, &f7elem)
			}
			ko.Spec.Operations = f7
		} else {
			ko.Spec.Operations = nil
		}
		if elem.RetiringPrincipal != nil {
			ko.Spec.RetiringPrincipal = elem.RetiringPrincipal
		} else {
			ko.Spec.RetiringPrincipal = nil
		}
		found = true
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadManyInput returns true if there are any fields
// for the ReadMany Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadManyInput(
	r *resource,
) bool {
	return r.ko.Status.GrantID == nil

}

// newListRequestPayload returns SDK-specific struct for the HTTP request
// payload of the List API call for the resource
func (rm *resourceManager) newListRequestPayload(
	r *resource,
) (*svcsdk.ListGrantsInput, error) {
	res := &svcsdk.ListGrantsInput{}

	if r.ko.Status.GrantID != nil {
		res.SetGrantId(*r.ko.Status.GrantID)
	}
	if r.ko.Spec.GranteePrincipal != nil {
		res.SetGranteePrincipal(*r.ko.Spec.GranteePrincipal)
	}
	if r.ko.Spec.KeyID != nil {
		res.SetKeyId(*r.ko.Spec.KeyID)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateGrantOutput
	_ = resp
	resp, err = rm.sdkapi.CreateGrantWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateGrant", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if resp.GrantId != nil {
		ko.Status.GrantID = resp.GrantId
	} else {
		ko.Status.GrantID = nil
	}
	if resp.GrantToken != nil {
		ko.Status.GrantToken = resp.GrantToken
	} else {
		ko.Status.GrantToken = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateGrantInput, error) {
	res := &svcsdk.CreateGrantInput{}

	if r.ko.Spec.Constraints != nil {
		f0 := &svcsdk.GrantConstraints{}
		if r.ko.Spec.Constraints.EncryptionContextEquals != nil {
			f0f0 := map[string]*string{}
			for f0f0key, f0f0valiter := range r.ko.Spec.Constraints.EncryptionContextEquals {
				var f0f0val string
				f0f0val = *f0f0valiter
				f0f0[f0f0key] = &f0f0val
			}
			f0.SetEncryptionContextEquals(f0f0)
		}
		if r.ko.Spec.Constraints.EncryptionContextSubset != nil {
			f0f1 := map[string]*string{}
			for f0f1key, f0f1valiter := range r.ko.Spec.Constraints.EncryptionContextSubset {
				var f0f1val string
				f0f1val = *f0f1valiter
				f0f1[f0f1key] = &f0f1val
			}
			f0.SetEncryptionContextSubset(f0f1)
		}
		res.SetConstraints(f0)
	}
	if r.ko.Spec.GrantTokens != nil {
		f1 := []*string{}
		for _, f1iter := range r.ko.Spec.GrantTokens {
			var f1elem string
			f1elem = *f1iter
			f1 = append(f1, &f1elem)
		}
		res.SetGrantTokens(f1)
	}
	if r.ko.Spec.GranteePrincipal != nil {
		res.SetGranteePrincipal(*r.ko.Spec.GranteePrincipal)
	}
	if r.ko.Spec.KeyID != nil {
		res.SetKeyId(*r.ko.Spec.KeyID)
	}
	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}
	if r.ko.Spec.Operations != nil {
		f5 := []*string{}
		for _, f5iter := range r.ko.Spec.Operations {
			var f5elem string
			f5elem = *f5iter
			f5 = append(f5, &f5elem)
		}
		res.SetOperations(f5)
	}
	if r.ko.Spec.RetiringPrincipal != nil {
		res.SetRetiringPrincipal(*r.ko.Spec.RetiringPrincipal)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	// TODO(jaypipes): Figure this out...
	return nil, ackerr.NotImplemented
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	// TODO(jaypipes): Figure this out...
	return nil, nil

}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Grant,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	// No terminal_errors specified for this resource in generator config
	return false
}
