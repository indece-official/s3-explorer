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
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/indece-official/s3-explorer/backend/src/generated/model/webapi"
	"github.com/indece-official/s3-explorer/backend/src/model"
)

func v1UpdateProfileJSONRequestBodyToProfileV1(profileID int64, requestBody *webapi.V1UpdateProfileJSONRequestBody) *model.ProfileV1 {
	profile := &model.ProfileV1{}

	profile.ID = profileID
	profile.Name = requestBody.Name
	profile.AccessKey = requestBody.AccessKey
	profile.SecretKey = requestBody.SecretKey
	profile.Region = requestBody.Region
	profile.Endpoint = requestBody.Endpoint
	profile.SSL = requestBody.Ssl
	profile.PathStyle = requestBody.PathStyle
	profile.Buckets = requestBody.Buckets

	return profile

}

func (c *Controller) reqAPIV1UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var err error

	log := c.log

	requestBody := &webapi.V1UpdateProfileJSONRequestBody{}

	profileIDStr := chi.URLParam(r, "profileID")
	profileID, err := strconv.ParseInt(profileIDStr, 10, 64)
	if err != nil {
		log.Warnf("%s %s 400 - Invalid profileID '%s': %s", r.Method, r.RequestURI, profileIDStr, err)

		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad request")
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(requestBody)
	if err != nil {
		log.Warnf("%s %s 400 - Decoding JSON request body failed: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(400)
		fmt.Fprintf(w, "Decoding request body failed")
		return
	}

	profile := v1UpdateProfileJSONRequestBodyToProfileV1(profileID, requestBody)

	err = c.settingsService.UpdateProfile(profile)
	if err != nil {
		log.Errorf("%s %s 500 - Can't update profile: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	log.Infof("%s %s 201 - Updated profile %d", r.Method, r.RequestURI, profileID)

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "{}")
}
