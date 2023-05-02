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
	"fmt"
	"net/http"
	"net/url"

	"github.com/indece-official/go-gousu/gousuchi/v2"
)

func (c *Controller) reqAPIV1DownloadProfileBucketObject(w http.ResponseWriter, r *http.Request) {
	errResp := c.checkAuth(r)
	if errResp != nil {
		errResp.Write(w)

		return
	}

	profileID, errResp := gousuchi.URLParamInt64(r, "profileID")
	if errResp != nil {
		errResp.Write(w)

		return
	}

	bucketName, errResp := gousuchi.URLParamString(r, "bucketName")
	if errResp != nil {
		errResp.Write(w)

		return
	}

	objectKeyRaw, errResp := gousuchi.URLParamString(r, "objectKey")
	if errResp != nil {
		errResp.Write(w)

		return
	}

	objectKey, err := url.QueryUnescape(objectKeyRaw)
	if err != nil {
		gousuchi.BadRequest(r, "Can't unescape objectKey: %s", err).Write(w)

		return
	}

	profile, err := c.settingsService.GetProfile(profileID)
	if err != nil {
		gousuchi.NotFound(r, "Can't load profile: %s", err).Write(w)

		return
	}

	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectKey))

	err = c.s3Service.DownloadObject(profile, bucketName, objectKey, w)
	if err != nil {
		gousuchi.InternalServerError(r, "Can't load object: %s", err).Write(w)

		return
	}

	c.log.Infof("%s %s 200 - Downloaded object", r.Method, r.RequestURI)
}
