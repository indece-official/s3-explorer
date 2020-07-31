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

package settings

import (
	"sort"

	"github.com/indece-official/s3-explorer/backend/src/model"
)

// GetProfiles returns all known profiles
func (s *Service) GetProfiles() ([]*model.ProfileV1, error) {
	profilesList := []*model.ProfileV1{}

	for _, profile := range s.profiles {
		profilesList = append(profilesList, profile)
	}

	sort.Sort(model.ByProfileV1ID(profilesList))

	return profilesList, nil
}
