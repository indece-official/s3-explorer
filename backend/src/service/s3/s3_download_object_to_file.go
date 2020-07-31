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
	"os"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/indece-official/s3-explorer/backend/src/model"
)

// ProgressCallback defines a function that is called on progress
type ProgressCallback func(written int64, total int64)

// ProgressWriter is a WriterAt that can be used to track the progress
// of the underlying WriterAt by calling a ProgressCallback function
type ProgressWriter struct {
	writer      io.WriterAt
	total       int64
	written     int64
	progressClb ProgressCallback
}

// WriteAt writes data to the underlying WriterAt, calling the ProgressCallback
// in a separate go "thread" asynchronously
func (p *ProgressWriter) WriteAt(buf []byte, off int64) (int, error) {
	atomic.AddInt64(&p.written, int64(len(buf)))
	go p.progressClb(p.written, p.total)

	return p.writer.WriteAt(buf, off)
}

// NewProgressWriter constructs a new ProgressWriter for an unerlying WriterAt
// with a total expected size of data being written and a ProgressCallback
func NewProgressWriter(writer io.WriterAt, total int64, progressClb ProgressCallback) *ProgressWriter {
	return &ProgressWriter{
		writer:      writer,
		total:       total,
		written:     int64(0),
		progressClb: progressClb,
	}
}

// DownloadObjectToFile downloads an Object to a file, tracking the download progress
// via a ProgressCallback
func (s *Service) DownloadObjectToFile(profile *model.ProfileV1, bucketName string, objectKey string, filename string, progressClb ProgressCallback) error {
	s3Client, err := s.getClient(profile)
	if err != nil {
		return err
	}

	s3Downloader, err := s.getDownloader(profile)
	if err != nil {
		return err
	}

	headObjectOutput, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := NewProgressWriter(
		file,
		*headObjectOutput.ContentLength,
		progressClb,
	)

	_, err = s3Downloader.Download(writer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
	if err != nil {
		return err
	}

	return nil
}
