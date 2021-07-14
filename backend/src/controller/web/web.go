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

package web

import (
	"fmt"
	"net/http"

	"github.com/indece-official/go-gousu"
	"github.com/indece-official/s3-explorer/backend/src/service/s3"
	"github.com/indece-official/s3-explorer/backend/src/service/settings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/namsral/flag"
)

// ControllerName defines the name of the api controller used for dependency exception
const ControllerName = "web"

var (
	serverHost = flag.String("server_host", "127.0.0.1", "")
	// ServerPort specifies the port the web server is listening on
	ServerPort   = flag.Int("server_port", 41100, "")
	cookieSecure = flag.Bool("cookie_secure", true, "")
	baseHREF     = flag.String("base_href", "", "")
)

// IController is the interface of the api controller
type IController interface {
	gousu.IController
}

// Controller is the admin api controller
type Controller struct {
	error           error
	log             *gousu.Log
	s3Service       s3.IService
	settingsService settings.IService
}

var _ IController = (*Controller)(nil)

func (c *Controller) getRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
	)

	router.Get("/*", c.reqStaticFile)

	router.Get("/api/v1/profile", c.reqAPIV1GetProfiles)
	router.Post("/api/v1/profile", c.reqAPIV1AddProfile)
	router.Put("/api/v1/profile/{profileID}", c.reqAPIV1UpdateProfile)
	router.Delete("/api/v1/profile/{profileID}", c.reqAPIV1DeleteProfile)

	router.Get("/api/v1/profile/{profileID}/bucket", c.reqAPIV1GetProfileBuckets)
	router.Post("/api/v1/profile/{profileID}/bucket", c.reqAPIV1AddProfileBucket)
	router.Delete("/api/v1/profile/{profileID}/bucket/{bucketName}", c.reqAPIV1DeleteProfileBucket)

	router.Get("/api/v1/profile/{profileID}/bucket/{bucketName}/stats", c.reqAPIV1GetProfileBucketStats)
	router.Get("/api/v1/profile/{profileID}/bucket/{bucketName}/object", c.reqAPIV1GetProfileBucketObjects)
	router.Post("/api/v1/profile/{profileID}/bucket/{bucketName}/object", c.reqAPIV1AddProfileBucketObject)
	router.Get("/api/v1/profile/{profileID}/bucket/{bucketName}/object/{objectKey}", c.reqAPIV1DownloadProfileBucketObject)
	router.Delete("/api/v1/profile/{profileID}/bucket/{bucketName}/object/{objectKey}", c.reqAPIV1DeleteProfileBucketObject)

	return router
}

// Name returns the name of the controller defined by ControllerName
func (c *Controller) Name() string {
	return ControllerName
}

// Start starts the api server in a new go-func
func (c *Controller) Start() error {
	c.error = nil

	go func() {
		router := c.getRouter()

		err := http.ListenAndServe(fmt.Sprintf("%s:%d", *serverHost, *ServerPort), router)
		if err != nil {
			c.error = err

			c.log.Errorf("Can't start api server: %s", err)
		}
	}()

	c.log.Infof("API server listening on %s:%d", *serverHost, *ServerPort)

	return nil
}

// Health checks if the api server has thrown unresolvable internal errors
func (c *Controller) Health() error {
	return c.error
}

// Stop currently does nothing
func (c *Controller) Stop() error {
	return nil
}

// NewController creates a new preinitialized instance of Controller
func NewController(ctx gousu.IContext) gousu.IController {
	return &Controller{
		log:             gousu.GetLogger(fmt.Sprintf("controller.%s", ControllerName)),
		s3Service:       ctx.GetService(s3.ServiceName).(s3.IService),
		settingsService: ctx.GetService(settings.ServiceName).(settings.IService),
	}
}

var _ gousu.ControllerFactory = NewController
