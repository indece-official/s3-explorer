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

	"github.com/indece-official/go-gousu/gousuchi/v2"
	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/indece-official/go-gousu/v2/gousu/logger"
	"github.com/indece-official/s3-explorer/backend/src/service/s3"
	"github.com/indece-official/s3-explorer/backend/src/service/session"
	"github.com/indece-official/s3-explorer/backend/src/service/settings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/namsral/flag"
)

// ControllerName defines the name of the api controller used for dependency exception
const ControllerName = "web"

var (
	serverHost = flag.String("server_host", "127.0.0.1", "")
	// ServerPort specifies the port the web server is listening on
	ServerPort  = flag.Int("server_port", 41100, "")
	disableAuth = flag.Bool("server_disable_auth", false, "")
)

// IController is the interface of the api controller
type IController interface {
	gousu.IController
}

// Controller is the admin api controller
type Controller struct {
	baseController  *gousuchi.AbstractController
	log             *logger.Log
	s3Service       s3.IService
	sessionService  session.IService
	settingsService settings.IService
}

var _ IController = (*Controller)(nil)

func (c *Controller) getRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
	)

	router.Get("/*", c.reqStaticFile)

	router.Get("/api/v1/profile", c.baseController.Wrap(c.reqAPIV1GetProfiles))
	router.Post("/api/v1/profile", c.baseController.Wrap(c.reqAPIV1AddProfile))
	router.Put("/api/v1/profile/{profileID}", c.baseController.Wrap(c.reqAPIV1UpdateProfile))
	router.Delete("/api/v1/profile/{profileID}", c.baseController.Wrap(c.reqAPIV1DeleteProfile))

	router.Get("/api/v1/profile/{profileID}/bucket", c.baseController.Wrap(c.reqAPIV1GetProfileBuckets))
	router.Post("/api/v1/profile/{profileID}/bucket", c.baseController.Wrap(c.reqAPIV1AddProfileBucket))
	router.Delete("/api/v1/profile/{profileID}/bucket/{bucketName}", c.baseController.Wrap(c.reqAPIV1DeleteProfileBucket))

	router.Get("/api/v1/profile/{profileID}/bucket/{bucketName}/stats", c.baseController.Wrap(c.reqAPIV1GetProfileBucketStats))
	router.Get("/api/v1/profile/{profileID}/bucket/{bucketName}/object", c.baseController.Wrap(c.reqAPIV1GetProfileBucketObjects))
	router.Post("/api/v1/profile/{profileID}/bucket/{bucketName}/object", c.baseController.Wrap(c.reqAPIV1AddProfileBucketObject))
	router.Get("/api/v1/profile/{profileID}/bucket/{bucketName}/object/{objectKey}", c.reqAPIV1DownloadProfileBucketObject)
	router.Delete("/api/v1/profile/{profileID}/bucket/{bucketName}/object/{objectKey}", c.baseController.Wrap(c.reqAPIV1DeleteProfileBucketObject))

	return router
}

// Name returns the name of the controller defined by ControllerName
func (c *Controller) Name() string {
	return ControllerName
}

// Start starts the api server in a new go-func
func (c *Controller) Start() error {
	c.baseController.UseHost(*serverHost)
	c.baseController.UsePort(*ServerPort)
	c.baseController.UseRouter(c.getRouter())

	return c.baseController.Start()
}

// Health checks if the api server has thrown unresolvable internal errors
func (c *Controller) Health() error {
	return c.baseController.Health()
}

// Stop currently does nothing
func (c *Controller) Stop() error {
	return c.baseController.Stop()
}

// NewController creates a new preinitialized instance of Controller
func NewController(ctx gousu.IContext) gousu.IController {
	log := logger.GetLogger(fmt.Sprintf("controller.%s", ControllerName))

	return &Controller{
		log:             log,
		baseController:  gousuchi.NewAbstractController(log),
		s3Service:       ctx.GetService(s3.ServiceName).(s3.IService),
		sessionService:  ctx.GetService(session.ServiceName).(session.IService),
		settingsService: ctx.GetService(settings.ServiceName).(settings.IService),
	}
}

var _ gousu.ControllerFactory = NewController
