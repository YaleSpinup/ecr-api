package ecr

import (
	"testing"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/pkg/errors"
)

func TestErrCode(t *testing.T) {
	apiErrorTestCases := map[string]string{
		"": apierror.ErrBadRequest,

		"AccessDenied": apierror.ErrForbidden,

		ecr.ErrCodeServerException: apierror.ErrInternalError,

		ecr.ErrCodeImageTagAlreadyExistsException: apierror.ErrConflict,

		ecr.ErrCodeLayerAlreadyExistsException:               apierror.ErrConflict,
		ecr.ErrCodeLifecyclePolicyPreviewInProgressException: apierror.ErrConflict,
		ecr.ErrCodeImageAlreadyExistsException:               apierror.ErrConflict,
		ecr.ErrCodeRepositoryAlreadyExistsException:          apierror.ErrConflict,

		ecr.ErrCodeEmptyUploadException:             apierror.ErrBadRequest,
		ecr.ErrCodeImageDigestDoesNotMatchException: apierror.ErrBadRequest,
		ecr.ErrCodeInvalidLayerException:            apierror.ErrBadRequest,
		ecr.ErrCodeInvalidLayerPartException:        apierror.ErrBadRequest,
		ecr.ErrCodeInvalidParameterException:        apierror.ErrBadRequest,
		ecr.ErrCodeInvalidTagParameterException:     apierror.ErrBadRequest,
		ecr.ErrCodeLayerInaccessibleException:       apierror.ErrBadRequest,
		ecr.ErrCodeLayerPartTooSmallException:       apierror.ErrBadRequest,
		ecr.ErrCodeRepositoryNotEmptyException:      apierror.ErrBadRequest,
		ecr.ErrCodeUnsupportedImageTypeException:    apierror.ErrBadRequest,

		ecr.ErrCodeLayersNotFoundException:                 apierror.ErrNotFound,
		ecr.ErrCodeLifecyclePolicyNotFoundException:        apierror.ErrNotFound,
		ecr.ErrCodeLifecyclePolicyPreviewNotFoundException: apierror.ErrNotFound,
		ecr.ErrCodeReferencedImagesNotFoundException:       apierror.ErrNotFound,
		ecr.ErrCodeRepositoryNotFoundException:             apierror.ErrNotFound,
		ecr.ErrCodeRepositoryPolicyNotFoundException:       apierror.ErrNotFound,
		ecr.ErrCodeScanNotFoundException:                   apierror.ErrNotFound,
		ecr.ErrCodeImageNotFoundException:                  apierror.ErrNotFound,
		ecr.ErrCodeUploadNotFoundException:                 apierror.ErrNotFound,

		ecr.ErrCodeLimitExceededException: apierror.ErrLimitExceeded,
		ecr.ErrCodeTooManyTagsException:   apierror.ErrLimitExceeded,
	}

	for awsErr, apiErr := range apiErrorTestCases {
		err := ErrCode("test error", awserr.New(awsErr, awsErr, nil))
		if aerr, ok := errors.Cause(err).(apierror.Error); ok {
			t.Logf("got apierror '%s'", aerr)
		} else {
			t.Errorf("expected cloudwatch error %s to be an apierror.Error %s, got %s", awsErr, apiErr, err)
		}
	}

	err := ErrCode("test error", errors.New("Unknown"))
	if aerr, ok := errors.Cause(err).(apierror.Error); ok {
		t.Logf("got apierror '%s'", aerr)
	} else {
		t.Errorf("expected unknown error to be an apierror.ErrInternalError, got %s", err)
	}
}
