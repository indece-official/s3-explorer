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

// +build linux

package settings

import (
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (s *Service) createDefaultKey(keyFilename string) (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	keyStr := base64.StdEncoding.EncodeToString(key)

	err = ioutil.WriteFile(keyFilename, []byte(keyStr), 0600)
	if err != nil {
		return "", err
	}

	return keyStr, nil
}

func (s *Service) getDefaultKey() (string, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	keyFilename := filepath.Join(confDir, "s3explorer.dat")

	data, err := ioutil.ReadFile(keyFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return s.createDefaultKey(keyFilename)
		}

		return "", err
	}

	return string(data), nil
}
