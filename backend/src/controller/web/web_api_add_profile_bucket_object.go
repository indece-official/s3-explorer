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

package web

import (
	"net/http"

	"github.com/indece-official/go-gousu/gousuchi/v2"
)

func (c *Controller) reqAPIV1AddProfileBucketObject(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	profileID, errResp := gousuchi.URLParamInt64(r, "profileID")
	if errResp != nil {
		return errResp
	}

	bucketName, errResp := gousuchi.URLParamString(r, "bucketName")
	if errResp != nil {
		return errResp
	}

	profile, err := c.settingsService.GetProfile(profileID)
	if err != nil {
		return gousuchi.NotFound(r, "Can't load profile: %s", err)
	}

	r.ParseMultipartForm(10 * 1024 * 1024 * 1024 * 1024) // Max 10GB

	filename := r.FormValue("filename")

	file, handler, err := r.FormFile("file")
	if err != nil {
		return gousuchi.BadRequest(r, "Can't load file: %s", err)
	}
	defer file.Close()

	if filename == "" {
		filename = handler.Filename
	}

	err = c.s3Service.AddObject(profile, bucketName, filename, file)
	if err != nil {
		return gousuchi.InternalServerError(r, "Can't add bucket: %s", err)
	}

	return gousuchi.JSON(r, map[string]interface{}{}).
		WithStatusCode(http.StatusCreated).
		WithDetailedMessage("Created object '%s'", handler.Filename)
}
