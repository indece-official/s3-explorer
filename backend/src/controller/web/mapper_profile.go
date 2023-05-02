package web

import (
	"github.com/indece-official/s3-explorer/backend/src/generated/model/webapi"
	"github.com/indece-official/s3-explorer/backend/src/model"
)

func (c *Controller) mapV1AddProfileJSONRequestBodyToProfileV1(requestBody *webapi.V1AddProfileJSONRequestBody) *model.ProfileV1 {
	profile := &model.ProfileV1{}

	profile.Name = requestBody.Name
	profile.AccessKey = requestBody.AccessKey
	profile.SecretKey = requestBody.SecretKey
	profile.Region = requestBody.Region
	profile.Endpoint = requestBody.Endpoint
	profile.SSL = requestBody.Ssl
	profile.PathStyle = requestBody.PathStyle
	profile.Buckets = requestBody.Buckets

	return profile
}

func (c *Controller) mapProfileV1ToAPIProfileV1(profile *model.ProfileV1) webapi.ProfileV1 {
	apiProfile := webapi.ProfileV1{}

	apiProfile.Id = profile.ID
	apiProfile.Name = profile.Name
	apiProfile.AccessKey = profile.AccessKey
	apiProfile.SecretKey = "***"
	apiProfile.Region = profile.Region
	apiProfile.Endpoint = profile.Endpoint
	apiProfile.Ssl = profile.SSL
	apiProfile.PathStyle = profile.PathStyle
	apiProfile.Buckets = profile.Buckets

	return apiProfile
}

func (c *Controller) mapV1UpdateProfileJSONRequestBodyToProfileV1(requestBody *webapi.V1UpdateProfileJSONRequestBody, oldProfile *model.ProfileV1) *model.ProfileV1 {
	// Clone old profile
	profile := *oldProfile

	profile.ID = oldProfile.ID
	profile.Name = requestBody.Name

	if requestBody.AccessKey != nil {
		profile.AccessKey = *requestBody.AccessKey
	}

	if requestBody.SecretKey != nil {
		profile.SecretKey = *requestBody.SecretKey
	}

	profile.Region = requestBody.Region
	profile.Endpoint = requestBody.Endpoint
	profile.SSL = requestBody.Ssl
	profile.PathStyle = requestBody.PathStyle
	profile.Buckets = requestBody.Buckets

	return &profile
}
