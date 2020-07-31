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

	"github.com/indece-official/s3-explorer/backend/src/generated/model/webapi"
	"github.com/indece-official/s3-explorer/backend/src/model"

	"github.com/go-chi/chi"
)

func bucketV1ToAPIBucketV1(bucket *model.BucketV1) webapi.BucketV1 {
	apiBucket := webapi.BucketV1{}

	apiBucket.Name = bucket.Name

	return apiBucket

}

func (c *Controller) reqAPIV1GetProfileBuckets(w http.ResponseWriter, r *http.Request) {
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

	profile, err := c.settingsService.GetProfile(profileID)
	if err != nil {
		log.Warnf("%s %s 404 - Can't load profile: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(404)
		fmt.Fprintf(w, "Profile not found")
		return
	}

	buckets, err := c.s3Service.GetBuckets(profile)
	if err != nil {
		log.Errorf("%s %s 500 - Can't load buckets: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "S3 Error: %s", err)
		return
	}

	response := webapi.V1GetProfileBucketsJSONResponseBody{
		Buckets: []webapi.BucketV1{},
	}

	for _, bucket := range buckets {
		response.Buckets = append(response.Buckets, bucketV1ToAPIBucketV1(bucket))
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Errorf("%s %s 500 - JSON-encoding response failed: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	log.Infof("%s %s 200 - Loaded %d buckets", r.Method, r.RequestURI, len(response.Buckets))

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseJSON)
}
