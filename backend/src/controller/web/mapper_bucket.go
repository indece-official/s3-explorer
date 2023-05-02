package web

import (
	"github.com/indece-official/s3-explorer/backend/src/generated/model/webapi"
	"github.com/indece-official/s3-explorer/backend/src/model"
)

func (c *Controller) mapV1AddProfileBucketJSONRequestBodyToProfileV1(requestBody *webapi.V1AddProfileBucketJSONRequestBody) *model.BucketV1 {
	bucket := &model.BucketV1{}

	bucket.Name = requestBody.Name

	return bucket
}

func (c *Controller) mapBucketV1ToAPIBucketV1(bucket *model.BucketV1) webapi.BucketV1 {
	apiBucket := webapi.BucketV1{}

	apiBucket.Name = bucket.Name

	return apiBucket
}

func (c *Controller) mapBucketStatsV1ToAPIBucketStatsV1(bucketStats *model.BucketStatsV1) webapi.BucketStatsV1 {
	apiBucketStats := webapi.BucketStatsV1{}

	apiBucketStats.Files = bucketStats.Files
	apiBucketStats.Size = bucketStats.Size
	apiBucketStats.Complete = bucketStats.Complete

	return apiBucketStats
}
