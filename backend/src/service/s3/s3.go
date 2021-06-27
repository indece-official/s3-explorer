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
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/indece-official/go-gousu"
	"github.com/indece-official/s3-explorer/backend/src/model"
	"gopkg.in/guregu/null.v4"
)

// ServiceName defines the name of the s3 service
const ServiceName = "s3"

// IService defines the interface of s3 service
type IService interface {
	gousu.IService

	GetBuckets(profile *model.ProfileV1) ([]*model.BucketV1, error)
	AddBucket(profile *model.ProfileV1, bucket *model.BucketV1) error
	DeleteBucket(profile *model.ProfileV1, bucketName string) error

	GetObjects(profile *model.ProfileV1, bucket string, continuationToken null.String, size int64) ([]*model.ObjectV1, string, error)
	GetObject(profile *model.ProfileV1, bucket string, objectKey string) (*model.ObjectV1, error)
	DownloadObject(profile *model.ProfileV1, bucket string, objectKey string, out io.Writer) error
	DownloadObjectToFile(profile *model.ProfileV1, bucket string, objectKey string, filename string, progressClb ProgressCallback) error
	AddObject(profile *model.ProfileV1, bucketName string, filename string, file multipart.File) error
	DeleteObject(profile *model.ProfileV1, bucketName string, objectKey string) error
}

// Service is the service used to access the s3 storage server
type Service struct {
	log *gousu.Log
}

var _ IService = (*Service)(nil)

func (s *Service) getClient(profile *model.ProfileV1) (*s3.S3, error) {
	awsSession, err := session.NewSessionWithOptions(profile.ToAwsOptions())
	if err != nil {
		return nil, err
	}

	return s3.New(awsSession), nil
}

func (s *Service) getUploader(profile *model.ProfileV1) (*s3manager.Uploader, error) {
	awsSession, err := session.NewSessionWithOptions(profile.ToAwsOptions())
	if err != nil {
		return nil, err
	}

	return s3manager.NewUploader(awsSession), nil
}

func (s *Service) getDownloader(profile *model.ProfileV1) (*s3manager.Downloader, error) {
	awsSession, err := session.NewSessionWithOptions(profile.ToAwsOptions())
	if err != nil {
		return nil, err
	}

	return s3manager.NewDownloader(awsSession), nil
}

// Name returns the service's name defined by ServiceName
func (s *Service) Name() string {
	return ServiceName
}

// Start does nothing
func (s *Service) Start() error {
	return nil
}

// Health returns always 'ok' (nil)
func (s *Service) Health() error {
	return nil
}

// Stop does nothing
func (s *Service) Stop() error {
	return nil
}

// NewService creates a new instance of s3 service
func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log: gousu.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
	}
}

var _ (gousu.ServiceFactory) = NewService
