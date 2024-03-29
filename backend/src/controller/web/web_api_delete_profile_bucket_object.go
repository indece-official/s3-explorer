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
	"strconv"

	"github.com/go-chi/chi"
)

func (c *Controller) reqAPIV1DeleteProfileBucketObject(w http.ResponseWriter, r *http.Request) {
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

	objectKeyRaw := chi.URLParam(r, "objectKey")
	if objectKeyRaw == "" {
		log.Warnf("%s %s 400 - Invalid objectKey '%s': %s", r.Method, r.RequestURI, objectKeyRaw, err)

		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad request")
		return
	}

	objectKey, err := url.QueryUnescape(objectKeyRaw)
	if err != nil {
		log.Warnf("%s %s 400 - Can't unescape objectKey: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad request")
		return
	}

	log.Infof("Object: %s", objectKey)

	profile, err := c.settingsService.GetProfile(profileID)
	if err != nil {
		log.Warnf("%s %s 404 - Can't load profile: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(404)
		fmt.Fprintf(w, "Profile not found")
		return
	}

	err = c.s3Service.DeleteObject(profile, bucketName, objectKey)
	if err != nil {
		log.Errorf("%s %s 500 - Can't delete object: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "S3 Error: %s", err)
		return
	}

	log.Infof("%s %s 200 - Deleted object", r.Method, r.RequestURI)

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "{}")
}
