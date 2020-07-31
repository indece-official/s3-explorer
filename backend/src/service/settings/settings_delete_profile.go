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
	"fmt"
)

// DeleteProfile deletes an existing profile and stores all settings to file
//
// Returns an error if the profile doesn't exist
func (s *Service) DeleteProfile(id int64) error {
	_, ok := s.profiles[id]
	if !ok {
		return fmt.Errorf("Profile %d not found", id)
	}

	delete(s.profiles, id)

	return s.store()
}
