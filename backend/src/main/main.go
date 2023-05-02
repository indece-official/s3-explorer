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

//go:generate go run ../assets/generate.go
//go:generate sh -c "mkdir -p ../generated/model/webapi && oapi-codegen --package=webapi --generate=types ../../assets/swagger/webapi.yml > ../generated/model/webapi/webapi.gen.go"
package main

import (
	"fmt"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/s3-explorer/backend/src/controller/ui"
	"github.com/indece-official/s3-explorer/backend/src/controller/web"
	"github.com/indece-official/s3-explorer/backend/src/service/s3"
	"github.com/indece-official/s3-explorer/backend/src/service/session"
	"github.com/indece-official/s3-explorer/backend/src/service/settings"
	"github.com/indece-official/s3-explorer/backend/src/utils"
)

// Variables set during build
var (
	ProjectName  = "dummy"
	BuildVersion string
	BuildDate    string
)

func main() {
	utils.Init()

	runner := gousu.NewRunner(ProjectName, fmt.Sprintf("%s (Build %s)", BuildVersion, BuildDate))

	runner.CreateService(settings.NewService)
	runner.CreateService(s3.NewService)
	runner.CreateService(session.NewService)
	runner.CreateController(web.NewController)
	runner.CreateUIController(ui.NewController)

	runner.Run()
}
