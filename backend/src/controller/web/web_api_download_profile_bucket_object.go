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
	"strconv"

	"github.com/go-chi/chi"
)

func (c *Controller) reqAPIV1DownloadProfileBucketObject(w http.ResponseWriter, r *http.Request) {
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

	objectKey := chi.URLParam(r, "objectKey")
	if objectKey == "" {
		log.Warnf("%s %s 400 - Invalid objectKey '%s': %s", r.Method, r.RequestURI, objectKey, err)

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

	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectKey))

	err = c.s3Service.DownloadObject(profile, bucketName, objectKey, w)
	if err != nil {
		log.Errorf("%s %s 500 - Can't load object: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	log.Infof("%s %s 200 - Downloaded object", r.Method, r.RequestURI)
}
