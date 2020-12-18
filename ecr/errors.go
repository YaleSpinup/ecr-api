package ecr

import (
	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func ErrCode(msg string, err error) error {
	log.Debugf("processing error code with message '%s' and error '%s'", msg, err)

	if aerr, ok := errors.Cause(err).(awserr.Error); ok {
		switch aerr.Code() {
		case
			"AccessDenied":

			return apierror.New(apierror.ErrForbidden, msg, aerr)
		case
			// ecr.ErrCodeServerException for service response error code
			// "ServerException".
			//
			// These errors are usually caused by a server-side issue.
			ecr.ErrCodeServerException:

			return apierror.New(apierror.ErrInternalError, msg, err)
		case
			// ecr.ErrCodeImageTagAlreadyExistsException for service response error code
			// "ImageTagAlreadyExistsException".
			//
			// The specified image is tagged with a tag that already exists. The repository
			// is configured for tag immutability.
			ecr.ErrCodeImageTagAlreadyExistsException,

			// ecr.ErrCodeLayerAlreadyExistsException for service response error code
			// "LayerAlreadyExistsException".
			//
			// The image layer already exists in the associated repository.
			ecr.ErrCodeLayerAlreadyExistsException,

			// ecr.ErrCodeLifecyclePolicyPreviewInProgressException for service response error code
			// "LifecyclePolicyPreviewInProgressException".
			//
			// The previous lifecycle policy preview request has not completed. Wait and
			// try again.
			ecr.ErrCodeLifecyclePolicyPreviewInProgressException,

			// ecr.ErrCodeImageAlreadyExistsException for service response error code
			// "ImageAlreadyExistsException".
			//
			// The specified image has already been pushed, and there were no changes to
			// the manifest or image tag after the last push.
			ecr.ErrCodeImageAlreadyExistsException,

			// ecr.ErrCodeRepositoryAlreadyExistsException for service response error code
			// "RepositoryAlreadyExistsException".
			//
			// The specified repository already exists in the specified registry.
			ecr.ErrCodeRepositoryAlreadyExistsException:

			return apierror.New(apierror.ErrConflict, msg, aerr)
		case
			// ecr.ErrCodeEmptyUploadException for service response error code
			// "EmptyUploadException".
			//
			// The specified layer upload does not contain any layer parts.
			ecr.ErrCodeEmptyUploadException,

			// ecr.ErrCodeImageDigestDoesNotMatchException for service response error code
			// "ImageDigestDoesNotMatchException".
			//
			// The specified image digest does not match the digest that Amazon ecr.ErrCodecalculated
			// for the image.
			ecr.ErrCodeImageDigestDoesNotMatchException,

			// ecr.ErrCodeInvalidLayerException for service response error code
			// "InvalidLayerException".
			//
			// The layer digest calculation performed by Amazon ecr.ErrCodeupon receipt of the
			// image layer does not match the digest specified.
			ecr.ErrCodeInvalidLayerException,

			// ecr.ErrCodeInvalidLayerPartException for service response error code
			// "InvalidLayerPartException".
			//
			// The layer part size is not valid, or the first byte specified is not consecutive
			// to the last byte of a previous layer part upload.
			ecr.ErrCodeInvalidLayerPartException,

			// ecr.ErrCodeInvalidParameterException for service response error code
			// "InvalidParameterException".
			//
			// The specified parameter is invalid. Review the available parameters for the
			// API request.
			ecr.ErrCodeInvalidParameterException,

			// ecr.ErrCodeInvalidTagParameterException for service response error code
			// "InvalidTagParameterException".
			//
			// An invalid parameter has been specified. Tag keys can have a maximum character
			// length of 128 characters, and tag values can have a maximum length of 256
			// characters.
			ecr.ErrCodeInvalidTagParameterException,

			// ecr.ErrCodeLayerInaccessibleException for service response error code
			// "LayerInaccessibleException".
			//
			// The specified layer is not available because it is not associated with an
			// image. Unassociated image layers may be cleaned up at any time.
			ecr.ErrCodeLayerInaccessibleException,

			// ecr.ErrCodeLayerPartTooSmallException for service response error code
			// "LayerPartTooSmallException".
			//
			// Layer parts must be at least 5 MiB in size.
			ecr.ErrCodeLayerPartTooSmallException,

			// ecr.ErrCodeRepositoryNotEmptyException for service response error code
			// "RepositoryNotEmptyException".
			//
			// The specified repository contains images. To delete a repository that contains
			// images, you must force the deletion with the force parameter.
			ecr.ErrCodeRepositoryNotEmptyException,

			// ecr.ErrCodeUnsupportedImageTypeException for service response error code
			// "UnsupportedImageTypeException".
			//
			// The image is of a type that cannot be scanned.
			ecr.ErrCodeUnsupportedImageTypeException:

			return apierror.New(apierror.ErrBadRequest, msg, aerr)
		case
			// ecr.ErrCodeLayersNotFoundException for service response error code
			// "LayersNotFoundException".
			//
			// The specified layers could not be found, or the specified layer is not valid
			// for this repository.
			ecr.ErrCodeLayersNotFoundException,

			// ecr.ErrCodeLifecyclePolicyNotFoundException for service response error code
			// "LifecyclePolicyNotFoundException".
			//
			// The lifecycle policy could not be found, and no policy is set to the repository.
			ecr.ErrCodeLifecyclePolicyNotFoundException,

			// ecr.ErrCodeLifecyclePolicyPreviewNotFoundException for service response error code
			// "LifecyclePolicyPreviewNotFoundException".
			//
			// There is no dry run for this repository.
			ecr.ErrCodeLifecyclePolicyPreviewNotFoundException,

			// ecr.ErrCodeReferencedImagesNotFoundException for service response error code
			// "ReferencedImagesNotFoundException".
			//
			// The manifest list is referencing an image that does not exist.
			ecr.ErrCodeReferencedImagesNotFoundException,

			// ecr.ErrCodeRepositoryNotFoundException for service response error code
			// "RepositoryNotFoundException".
			//
			// The specified repository could not be found. Check the spelling of the specified
			// repository and ensure that you are performing operations on the correct registry.
			ecr.ErrCodeRepositoryNotFoundException,

			// ecr.ErrCodeRepositoryPolicyNotFoundException for service response error code
			// "RepositoryPolicyNotFoundException".
			//
			// The specified repository and registry combination does not have an associated
			// repository policy.
			ecr.ErrCodeRepositoryPolicyNotFoundException,

			// ecr.ErrCodeScanNotFoundException for service response error code
			// "ScanNotFoundException".
			//
			// The specified image scan could not be found. Ensure that image scanning is
			// enabled on the repository and try again.
			ecr.ErrCodeScanNotFoundException,

			// ecr.ErrCodeImageNotFoundException for service response error code
			// "ImageNotFoundException".
			//
			// The image requested does not exist in the specified repository.
			ecr.ErrCodeImageNotFoundException,

			// ecr.ErrCodeUploadNotFoundException for service response error code
			// "UploadNotFoundException".
			//
			// The upload could not be found, or the specified upload ID is not valid for
			// this repository.
			ecr.ErrCodeUploadNotFoundException:

			return apierror.New(apierror.ErrNotFound, msg, aerr)
		case
			// ecr.ErrCodeLimitExceededException for service response error code
			// "LimitExceededException".
			//
			// The operation did not succeed because it would have exceeded a service limit
			// for your account. For more information, see Amazon ecr.ErrCodeService Quotas (https://docs.aws.amazon.com/Amazonecr.ErrCodelatest/userguide/service-quotas.html)
			// in the Amazon Elastic Container Registry User Guide.
			ecr.ErrCodeLimitExceededException,

			// ecr.ErrCodeTooManyTagsException for service response error code
			// "TooManyTagsException".
			//
			// The list of tags on the repository is over the limit. The maximum number
			// of tags that can be applied to a repository is 50.
			ecr.ErrCodeTooManyTagsException:

			return apierror.New(apierror.ErrLimitExceeded, msg, aerr)
		default:
			m := msg + ": " + aerr.Message()
			return apierror.New(apierror.ErrBadRequest, m, aerr)
		}
	}

	return apierror.New(apierror.ErrInternalError, msg, err)
}
