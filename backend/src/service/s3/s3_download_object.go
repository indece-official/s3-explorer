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
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/indece-official/s3-explorer/backend/src/model"
)

// FakeWriterAt can be used to pipe an WriterAt to an regular Writer
//
// Important: It must be made sure the WriterAt is always called without "jumping" "At"s!
type FakeWriterAt struct {
	w io.Writer
}

// WriteAt writes to the writer and ignores the offset
func (fw FakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	return fw.w.Write(p)
}

// DownloadObject downloads an object to an writer
//
// Returns nil on success after the download has finished, else error
func (s *Service) DownloadObject(profile *model.ProfileV1, bucket string, objectKey string, out io.Writer) error {
	s3Downloader, err := s.getDownloader(profile)
	if err != nil {
		return err
	}

	// Limit concurrency to 1 so we can use FakeWriterAt
	s3Downloader.Concurrency = 1

	_, err = s3Downloader.Download(FakeWriterAt{out},
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		})
	if err != nil {
		return err
	}

	return nil
}
