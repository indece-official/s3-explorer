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

package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// ProfileV1 contains one user profile for accessing a S3 storage
type ProfileV1 struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	AccessKey string   `json:"access_key"`
	SecretKey string   `json:"secret_key"`
	Endpoint  string   `json:"endpoint"`
	SSL       bool     `json:"ssl"`
	PathStyle bool     `json:"path_style"`
	Region    string   `json:"region"`
	Buckets   []string `json:"buckets"`
}

// ToAwsOptions generates the Options from the profile required for
// creating a new aws client connection
func (p *ProfileV1) ToAwsOptions() session.Options {
	config := aws.Config{}
	config.Credentials = credentials.NewStaticCredentials(p.AccessKey, p.SecretKey, "")

	if p.Endpoint != "" {
		config.Endpoint = aws.String(p.Endpoint)
	}

	if p.Region != "" {
		config.Region = aws.String(p.Region)
	} else {
		config.Region = aws.String("us-west-2")
	}

	config.DisableSSL = aws.Bool(!p.SSL)
	config.S3ForcePathStyle = aws.Bool(p.PathStyle)

	return session.Options{
		Config:  config,
		Profile: p.Name,
	}
}

// ByProfileV1ID is the accessor for sorting a list of profiles by their ID
type ByProfileV1ID []*ProfileV1

func (a ByProfileV1ID) Len() int           { return len(a) }
func (a ByProfileV1ID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByProfileV1ID) Less(i, j int) bool { return a[i].ID < a[j].ID }
