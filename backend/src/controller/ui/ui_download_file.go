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

package ui

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/indece-official/s3-explorer/backend/src/utils"
)

func (c *Controller) downloadFile(profileID int64, bucketName string, objectKey string) (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("Can't generate uuid for download: %s", err)
	}

	downloadUID := uid.String()

	profile, err := c.settingsService.GetProfile(profileID)
	if err != nil {
		return "", fmt.Errorf("Can't load profile: %s", err)
	}

	object, err := c.s3Service.GetObject(profile, bucketName, objectKey)
	if err != nil {
		return "", fmt.Errorf("Can't load object: %s", err)
	}

	filename, err := utils.SaveFile(c.view.Window(), filepath.Base(object.Key))
	if err != nil {
		return "", fmt.Errorf("Error opening Save-File-Dialog: %s", err)
	}

	if filename == "" {
		// No filename selected / cancel

		return "", nil
	}

	var muxDownloadStatus sync.Mutex

	downloadStatus := &DownloadStatus{
		UID:        downloadUID,
		ProfileID:  profileID,
		BucketName: bucketName,
		ObjectKey:  objectKey,
		Loaded:     0,
		Total:      0,
		Finished:   false,
		Error:      nil,
	}

	c.downloads[downloadUID] = downloadStatus

	go func() {
		onProgress := func(written int64, total int64) {
			muxDownloadStatus.Lock()
			downloadStatus.Loaded = written
			downloadStatus.Total = total
			muxDownloadStatus.Unlock()
		}

		err := c.s3Service.DownloadObjectToFile(profile, bucketName, objectKey, filename, onProgress)
		if err != nil {
			muxDownloadStatus.Lock()
			downloadStatus.Finished = true
			downloadStatus.Error = err
			muxDownloadStatus.Unlock()
			return
		}

		muxDownloadStatus.Lock()
		downloadStatus.Finished = true
		downloadStatus.Error = nil
		downloadStatus.Loaded = downloadStatus.Total
		muxDownloadStatus.Unlock()
	}()

	return downloadUID, nil
}
