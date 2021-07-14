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
	"testing"

	"github.com/indece-official/go-gousu"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDescript(t *testing.T) {
	ctx := gousu.NewContext()
	service := NewService(ctx).(*Service)

	rawData := []byte("This is a test")

	encData, err := service.encrypt(rawData, "abcdef123456")
	assert.NoError(t, err)
	assert.NotEmpty(t, encData)

	decData, err := service.decrypt(encData, "abcdef123456")
	assert.NoError(t, err)
	assert.Equal(t, rawData, decData)
}
