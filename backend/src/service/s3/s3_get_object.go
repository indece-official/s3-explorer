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
)

// GetObject loads the objects details
//
// Important: Does not load the owner!
func (s *Service) GetObject(profile *model.ProfileV1, bucketName string, objectKey string) (*model.ObjectV1, error) {
	s3Client, err := s.getClient(profile)
	if err != nil {
		return nil, err
	}

	headObjectOutput, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}

	object := &model.ObjectV1{
		Key:          objectKey,
		LastModified: *headObjectOutput.LastModified,
		OwnerName:    "",
		OwnerID:      "",
		Size:         *headObjectOutput.ContentLength,
	}

	return object, nil
}
