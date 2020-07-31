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

func objectV1ToAPIObjectV1(object *model.ObjectV1) webapi.ObjectV1 {
	apiObject := webapi.ObjectV1{}

	apiObject.Key = object.Key
	apiObject.LastModified = object.LastModified
	apiObject.OwnerName = object.OwnerName
	apiObject.OwnerId = object.OwnerID
	apiObject.Size = object.Size

	return apiObject

}

func (c *Controller) reqAPIV1GetProfileBucketObjects(w http.ResponseWriter, r *http.Request) {
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

	bucketName := chi.URLParam(r, "bucketName")
	if bucketName == "" {
		log.Warnf("%s %s 400 - Invalid bucketName '%s': %s", r.Method, r.RequestURI, bucketName, err)

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

	objects, err := c.s3Service.GetObjects(profile, bucketName)
	if err != nil {
		log.Errorf("%s %s 500 - Can't load objects: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "S3 Error: %s", err)
		return
	}

	response := webapi.V1GetProfileBucketObjectsJSONResponseBody{
		Objects: []webapi.ObjectV1{},
	}

	for _, object := range objects {
		response.Objects = append(response.Objects, objectV1ToAPIObjectV1(object))
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Errorf("%s %s 500 - JSON-encoding response failed: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	log.Infof("%s %s 200 - Loaded %d objects", r.Method, r.RequestURI, len(response.Objects))

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseJSON)
}
