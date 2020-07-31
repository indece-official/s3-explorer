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
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/indece-official/s3-explorer/backend/src/model"
)

// AddObject uploads an object to the bucket in the S3 storage specified by the passed profile.
// If an object with this key already exists in this bucket it is replaced by the new file
//
// Returns nil on success or else error
func (s *Service) AddObject(profile *model.ProfileV1, bucketName string, filename string, file multipart.File) error {
	uploader, err := s.getUploader(profile)
	if err != nil {
		return err
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}
