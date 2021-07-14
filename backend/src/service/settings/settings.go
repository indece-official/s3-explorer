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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/indece-official/go-gousu"
	"github.com/indece-official/s3-explorer/backend/src/model"
)

// ServiceName defines the name of the settings service
const ServiceName = "settings"

var flagFilename = flag.String("settings_file", "~/.s3explorer/config.enc", "")

// IService defines the interface of the settings service
type IService interface {
	gousu.IService

	GetProfiles() ([]*model.ProfileV1, error)
	GetProfile(id int64) (*model.ProfileV1, error)
	AddProfile(profile *model.ProfileV1) (int64, error)
	UpdateProfile(profile *model.ProfileV1) error
	DeleteProfile(id int64) error
}

// Service provides a service for accesisng the settings file
type Service struct {
	log      *gousu.Log
	filename string
	profiles map[int64]*model.ProfileV1
	loaded   bool
	key      string
}

var _ IService = (*Service)(nil)

func (s *Service) expandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

// store saves all currently configured settings to file, creating the
// containing dir if required
func (s *Service) store() error {
	settings := &model.SettingsV1{
		Profiles: []*model.ProfileV1{},
	}

	for _, profile := range s.profiles {
		settings.Profiles = append(settings.Profiles, profile)
	}

	settingsData, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("Can't encode settings: %s", err)
	}

	encryptedData, err := s.encrypt(settingsData, s.key)
	if err != nil {
		return fmt.Errorf("can't encrypt settings: %s", err)
	}

	// Make sure directory exists
	path := filepath.Dir(s.filename)
	err = os.MkdirAll(path, 0700)
	if err != nil {
		return fmt.Errorf("Can't create directory '%s' for settings file: %s", path, err)
	}

	err = ioutil.WriteFile(s.filename, encryptedData, 0600)
	if err != nil {
		return fmt.Errorf("Can't write settings to '%s': %s", s.filename, err)
	}

	return nil
}

// Name returns the name of the settings service defined by ServiceName
func (s *Service) Name() string {
	return ServiceName
}

func (s *Service) load() error {
	var err error

	settings := &model.SettingsV1{
		Profiles: []*model.ProfileV1{},
	}

	s.profiles = map[int64]*model.ProfileV1{}

	encryptedData, err := ioutil.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			s.loaded = true

			return nil
		}

		return err
	}

	settingsData, err := s.decrypt(encryptedData, s.key)
	if err != nil {
		return fmt.Errorf("can't decrypt settings: %s", err)
	}

	err = json.Unmarshal(settingsData, settings)
	if err != nil {
		return fmt.Errorf("Can't decode settings from '%s': %s", s.filename, err)
	}

	for _, profile := range settings.Profiles {
		s.profiles[profile.ID] = profile
	}

	s.loaded = true

	return nil
}

func (s *Service) UpdateKey(newKey string) error {
	s.key = newKey

	return s.store()
}

func (s *Service) UseKey(key string) error {
	s.key = key

	return s.load()
}

// Start loads an existing settings file if it exists, else the settings are created
// in memory an written to file later on the first change
func (s *Service) Start() error {
	var err error

	s.filename, err = s.expandPath(*flagFilename)
	if err != nil {
		return fmt.Errorf("Can't settings file expand path (%s): %s", *flagFilename, err)
	}

	s.key, err = s.getDefaultKey()
	if err != nil {
		s.log.Errorf("can't get default key: %s", err)
	}

	s.log.Infof("Loading file with key %s", s.key)

	err = s.load()
	if err != nil {
		s.log.Errorf("can't load settings: %s", err)
	}

	return nil
}

// Health returns always nil (healthy)
func (s *Service) Health() error {
	return nil
}

// Stop does nothing
func (s *Service) Stop() error {
	return nil
}

// NewService creates a new instance of the settings Service
func NewService(ctx gousu.IContext) gousu.IService {
	return &Service{
		log: gousu.GetLogger(fmt.Sprintf("service.%s", ServiceName)),
	}
}

var _ (gousu.ServiceFactory) = NewService
