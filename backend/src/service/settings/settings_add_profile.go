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
	"github.com/indece-official/s3-explorer/backend/src/model"
)

// AddProfile adds a new profile to the settings, stores them to file and returns
// the profile's new ID
func (s *Service) AddProfile(profile *model.ProfileV1) (int64, error) {
	maxProfileID := int64(0)
	for profileID := range s.profiles {
		if profileID > maxProfileID {
			maxProfileID = profileID
		}
	}

	newProfileID := maxProfileID + 1

	profile.ID = newProfileID

	s.profiles[newProfileID] = profile

	err := s.store()
	if err != nil {
		return 0, err
	}

	return newProfileID, nil
}
