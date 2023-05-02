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
	"github.com/indece-official/s3-explorer/backend/src/generated/model/webapi"
)

func (c *Controller) reqAPIV1GetProfileBucketObjects(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
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

	size, errResp := gousuchi.OptionalQueryParamInt64(r, "size")
	if errResp != nil {
		return errResp
	}

	if !size.Valid {
		size.Scan(100)
	}

	continuationToken, errResp := gousuchi.OptionalQueryParamString(r, "continuation_token")
	if errResp != nil {
		return errResp
	}

	profile, err := c.settingsService.GetProfile(profileID)
	if err != nil {
		return gousuchi.NotFound(r, "Can't load profile: %s", err)
	}

	objects, newContinuationToken, err := c.s3Service.GetObjects(profile, bucketName, continuationToken, size.Int64)
	if err != nil {
		return gousuchi.InternalServerError(r, "Can't load objects: %s", err)
	}

	respData := webapi.V1GetProfileBucketObjectsJSONResponseBody{
		Objects:           []webapi.ObjectV1{},
		ContinuationToken: newContinuationToken,
	}

	for _, object := range objects {
		respData.Objects = append(respData.Objects, c.mapObjectV1ToAPIObjectV1(object))
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded %d objects", len(respData.Objects))
}
