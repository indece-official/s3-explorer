// S3 Explorer
// Copyright (C) 2020  indece UG (haftungsbeschränkt)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License or any
// later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/indece-official/s3-explorer/backend/src/model"
	"gopkg.in/guregu/null.v4"
)

// GetObjects loads a list of max. 1000 objects from the given bucket
func (s *Service) GetObjects(profile *model.ProfileV1, bucket string, continuationToken null.String, size int64) ([]*model.ObjectV1, string, error) {
	objects := []*model.ObjectV1{}

	s3Client, err := s.getClient(profile)
	if err != nil {
		return nil, "", err
	}

	listObjectsOutput, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:            aws.String(bucket),
		FetchOwner:        aws.Bool(true),
		ContinuationToken: continuationToken.Ptr(),
		MaxKeys:           aws.Int64(size),
	})
	if err != nil {
		return nil, "", err
	}

	for _, object := range listObjectsOutput.Contents {
		objects = append(objects, &model.ObjectV1{
			Key:          aws.StringValue(object.Key),
			LastModified: aws.TimeValue(object.LastModified),
			OwnerName:    aws.StringValue(object.Owner.DisplayName),
			OwnerID:      aws.StringValue(object.Owner.ID),
			Size:         aws.Int64Value(object.Size),
		})
	}

	nextContinuationToken := ""

	if aws.BoolValue(listObjectsOutput.IsTruncated) {
		nextContinuationToken = aws.StringValue(listObjectsOutput.NextContinuationToken)
	}

	return objects, nextContinuationToken, nil
}
