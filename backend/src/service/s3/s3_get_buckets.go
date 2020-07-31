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
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/indece-official/s3-explorer/backend/src/model"
)

// GetBuckets lists the buckets in the S3 Storage specified by the passed profile.
// If the profile contains pre-defined buckets, these are returned instead and ListBuckets
// is not called (this can be usefull if the profile's user doesn't have the ListBucket-permission)
func (s *Service) GetBuckets(profile *model.ProfileV1) ([]*model.BucketV1, error) {
	buckets := []*model.BucketV1{}

	if len(profile.Buckets) > 0 {
		for _, bucket := range profile.Buckets {
			buckets = append(buckets, &model.BucketV1{
				Name: bucket,
			})
		}

		return buckets, nil
	}

	s3Client, err := s.getClient(profile)
	if err != nil {
		return nil, err
	}

	listBucketsOutput, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	for _, bucket := range listBucketsOutput.Buckets {
		buckets = append(buckets, &model.BucketV1{
			Name: *bucket.Name,
		})
	}

	return buckets, nil
}
