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

func (c *Controller) reqAPIV1GetProfileBucketStats(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
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

	force, errResp := gousuchi.OptionalQueryParamBool(r, "force")
	if errResp != nil {
		return errResp
	}

	profile, err := c.settingsService.GetProfile(profileID)
	if err != nil {
		return gousuchi.NotFound(r, "Can't load profile: %s", err)
	}

	stats, err := c.s3Service.GetBucketStats(profile, bucketName, force.Valid && force.Bool)
	if err != nil {
		return gousuchi.InternalServerError(r, "Can't load bucket stats: %s", err)
	}

	respData := webapi.V1GetProfileBucketStatsJSONResponseBody{
		Stats: c.mapBucketStatsV1ToAPIBucketStatsV1(stats),
	}

	return gousuchi.JSON(r, respData).
		WithDetailedMessage("Loaded stats")
}
