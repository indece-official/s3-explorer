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

// +build windows

package settings

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/sys/windows/registry"
)

func (s *Service) createDefaultKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	keyStr := base64.StdEncoding.EncodeToString(key)

	regKey, _, err := registry.CreateKey(registry.CURRENT_USER, `SOFTWARE\indece\S3Explorer\v1`, registry.WRITE)
	if err != nil {
		return "", err
	}
	defer regKey.Close()

	err = regKey.SetStringValue("DefaultKey", keyStr)
	if err != nil {
		return "", err
	}

	return keyStr, nil
}

func (s *Service) getDefaultKey() (string, error) {
	regKey, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\indece\S3Explorer\v1`, registry.QUERY_VALUE)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return s.createDefaultKey()
		}

		return "", err
	}
	defer regKey.Close()

	regVal, _, err := regKey.GetStringValue("DefaultKey")
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return s.createDefaultKey()
		}

		return "", err
	}

	if regVal == "" {
		return s.createDefaultKey()
	}

	return regVal, nil
}
