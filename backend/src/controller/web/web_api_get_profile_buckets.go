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

func (c *Controller) reqAPIV1GetProfileBuckets(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	profileID, errResp := gousuchi.URLParamInt64(r, "profileID")
	if errResp != nil {
		return errResp
	}

	profile, err := c.settingsService.GetProfile(profileID)
	if err != nil {
		return gousuchi.NotFound(r, "Can't load profile: %s", err)
	}

	buckets, err := c.s3Service.GetBuckets(profile)
	if err != nil {
		return gousuchi.InternalServerError(r, "Can't load buckets: %s", err)
	}

	respData := webapi.V1GetProfileBucketsJSONResponseBody{
		Buckets: []webapi.BucketV1{},
	}

	for _, bucket := range buckets {
		respData.Buckets = append(respData.Buckets, c.mapBucketV1ToAPIBucketV1(bucket))
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded %d buckets", len(respData.Buckets))
}
