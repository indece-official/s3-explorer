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

package ui

import "fmt"

func (c *Controller) getDownloadStatus(downloadUID string) (*DownloadStatus, error) {
	downloadStatus, ok := c.downloads[downloadUID]
	if !ok {
		return nil, fmt.Errorf("Download with uid '%s' not found", downloadUID)
	}

	return downloadStatus, nil
}
