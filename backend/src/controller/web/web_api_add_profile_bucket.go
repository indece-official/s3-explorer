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

	"github.com/indece-official/go-gousu"
	"github.com/indece-official/s3-explorer/backend/src/generated/model/webapi"
	"github.com/indece-official/s3-explorer/backend/src/model"

	"github.com/go-chi/chi"
)

func v1AddProfileBucketJSONRequestBodyToProfileV1(requestBody *webapi.V1AddProfileBucketJSONRequestBody) *model.BucketV1 {
	bucket := &model.BucketV1{}

	bucket.Name = requestBody.Name

	return bucket

}

func (c *Controller) reqAPIV1AddProfileBucket(w http.ResponseWriter, r *http.Request) {
	var err error

	log := c.log

	profileIDStr := chi.URLParam(r, "profileID")
	profileID, err := strconv.ParseInt(profileIDStr, 10, 64)
	if err != nil {
		log.Warnf("%s %s 400 - Invalid profileID '%s': %s", r.Method, r.RequestURI, profileIDStr, err)

		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad request")
		return
	}

	requestBody := &webapi.V1AddProfileBucketJSONRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(requestBody)
	if err != nil {
		log.Warnf("%s %s 400 - Decoding JSON request body failed: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(400)
		fmt.Fprintf(w, "Decoding request body failed")
		return
	}

	profile, err := c.settingsService.GetProfile(profileID)
	if err != nil {
		log.Warnf("%s %s 404 - Can't load profile: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(404)
		fmt.Fprintf(w, "Profile not found")
		return
	}

	bucket := v1AddProfileBucketJSONRequestBodyToProfileV1(requestBody)

	err = c.s3Service.AddBucket(profile, bucket)
	if err != nil {
		log.Errorf("%s %s 500 - Can't add bucket: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "S3 Error: %s", err)
		return
	}

	if len(profile.Buckets) > 0 && !gousu.ContainsString(profile.Buckets, bucket.Name) {
		profile.Buckets = append(profile.Buckets, bucket.Name)

		err = c.settingsService.UpdateProfile(profile)
		if err != nil {
			log.Errorf("%s %s 500 - Can't update profile: %s", r.Method, r.RequestURI, err)

			w.WriteHeader(500)
			fmt.Fprintf(w, "Internal server error")
			return
		}
	}

	log.Infof("%s %s 201 - Created bucket '%s'", r.Method, r.RequestURI, bucket.Name)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "{}")
}
