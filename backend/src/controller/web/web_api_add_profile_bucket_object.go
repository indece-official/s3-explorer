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

func (c *Controller) reqAPIV1AddProfileBucketObject(w http.ResponseWriter, r *http.Request) {
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

	r.ParseMultipartForm(10 * 1024 * 1024 * 1024 * 1024) // Max 10GB

	filename := r.FormValue("filename")

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Warnf("%s %s 400 - Can't load file: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad request")
		return
	}
	defer file.Close()

	if filename == "" {
		filename = handler.Filename
	}

	err = c.s3Service.AddObject(profile, bucketName, filename, file)
	if err != nil {
		log.Errorf("%s %s 500 - Can't add bucket: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "S3 Error: %s", err)
		return
	}

	log.Infof("%s %s 201 - Created object '%s'", r.Method, r.RequestURI, handler.Filename)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "{}")
}
