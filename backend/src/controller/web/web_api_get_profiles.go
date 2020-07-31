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

func profileV1ToAPIProfileV1(profile *model.ProfileV1) webapi.ProfileV1 {
	apiProfile := webapi.ProfileV1{}

	apiProfile.Id = profile.ID
	apiProfile.Name = profile.Name
	apiProfile.AccessKey = profile.AccessKey
	apiProfile.SecretKey = profile.SecretKey
	apiProfile.Region = profile.Region
	apiProfile.Endpoint = profile.Endpoint
	apiProfile.Ssl = profile.SSL
	apiProfile.PathStyle = profile.PathStyle
	apiProfile.Buckets = profile.Buckets

	return apiProfile

}

func (c *Controller) reqAPIV1GetProfiles(w http.ResponseWriter, r *http.Request) {
	var err error

	log := c.log

	profiles, err := c.settingsService.GetProfiles()
	if err != nil {
		log.Errorf("%s %s 500 - Can't load profiles: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	response := webapi.V1GetProfilesJSONResponseBody{
		Profiles: []webapi.ProfileV1{},
	}

	for _, profile := range profiles {
		response.Profiles = append(response.Profiles, profileV1ToAPIProfileV1(profile))
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Errorf("%s %s 500 - JSON-encoding response failed: %s", r.Method, r.RequestURI, err)

		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error")
		return
	}

	log.Infof("%s %s 200 - Loaded %d profiles", r.Method, r.RequestURI, len(response.Profiles))

	w.Header().Add("Content-Type", "application/json")
	w.Write(responseJSON)
}
