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
	"encoding/json"
	"net/http"

	"github.com/indece-official/go-gousu/gousuchi/v2"
	"github.com/indece-official/s3-explorer/backend/src/generated/model/webapi"
)

func (c *Controller) reqAPIV1AddProfile(w http.ResponseWriter, r *http.Request) gousuchi.IResponse {
	errResp := c.checkAuth(r)
	if errResp != nil {
		return errResp
	}

	requestBody := &webapi.V1AddProfileJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(requestBody)
	if err != nil {
		return gousuchi.BadRequest(r, "Decoding JSON request body failed: %s", err)
	}

	profile := c.mapV1AddProfileJSONRequestBodyToProfileV1(requestBody)

	profileID, err := c.settingsService.AddProfile(profile)
	if err != nil {
		return gousuchi.InternalServerError(r, "Can't add profile: %s", err)
	}

	respData := webapi.V1AddProfileJSONResponseBody{
		ProfileId: profileID,
	}

	return gousuchi.JSON(r, respData).
		WithStatusCode(http.StatusCreated).
		WithDetailedMessage("Added new profile %d", profileID)
}
