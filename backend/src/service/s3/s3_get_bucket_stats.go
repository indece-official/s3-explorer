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

func (s *Service) GetBucketStats(profile *model.ProfileV1, bucket string, force bool) (*model.BucketStatsV1, error) {
	s3Client, err := s.getClient(profile)
	if err != nil {
		return nil, err
	}

	stats := &model.BucketStatsV1{}

	continuationToken := null.String{}

	for {
		listObjectsOutput, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			FetchOwner:        aws.Bool(false),
			ContinuationToken: continuationToken.Ptr(),
			MaxKeys:           aws.Int64(1000),
		})
		if err != nil {
			return nil, err
		}

		for _, object := range listObjectsOutput.Contents {
			stats.Files++
			stats.Size += aws.Int64Value(object.Size)
		}

		if !aws.BoolValue(listObjectsOutput.IsTruncated) {
			stats.Complete = true
			break
		}

		if !force {
			break
		}

		continuationToken.Scan(aws.StringValue(listObjectsOutput.NextContinuationToken))
	}

	return stats, nil
}
