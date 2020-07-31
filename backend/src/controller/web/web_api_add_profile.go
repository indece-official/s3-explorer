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

	"github.com/indece-official/s3-explorer/backend/src/generated/model/webapi"
	"github.com/indece-official/s3-explorer/backend/src/model"
)

func v1AddProfileJSONRequestBodyToProfileV1(requestBody *webapi.V1AddProfileJSONRequestBody) *model.ProfileV1 {
	profile := &model.ProfileV1{}

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

func (c *Controller) reqAPIV1AddProfile(w http.ResponseWriter, r *http.Request) {
	var err error

	log := c.log

	requestBody := &webapi.V1AddProfileJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(requestBody)
	if err != nil {
		log.Warnf("%s %s 400 - Decoding JSON request body failed: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(400)
		fmt.Fprintf(w, "Decoding request body failed")
		return
	}

	profile := v1AddProfileJSONRequestBodyToProfileV1(requestBody)

	profileID, err := c.settingsService.AddProfile(profile)
	if err != nil {
		log.Errorf("%s %s 500 - Can't add profile: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	response := webapi.V1AddProfileJSONResponseBody{
		ProfileId: profileID,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Errorf("%s %s 500 - JSON-encoding response failed: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	log.Infof("%s %s 201 - Added new profile %d", r.Method, r.RequestURI, profileID)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(responseJSON)
}
