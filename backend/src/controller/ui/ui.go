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

import (
	"fmt"
	"os"
	"time"

	"github.com/zserge/webview"

	"github.com/indece-official/go-gousu"
	"github.com/indece-official/s3-explorer/backend/src/controller/web"
	"github.com/indece-official/s3-explorer/backend/src/service/s3"
	"github.com/indece-official/s3-explorer/backend/src/service/settings"
	"github.com/indece-official/s3-explorer/backend/src/utils"
	"github.com/namsral/flag"
)

var debug = flag.Bool("debug", false, "")

// ControllerName defines the name of the ui controller
const ControllerName = "ui"

// DownloadStatus contains the status & progress of a download
type DownloadStatus struct {
	UID        string `json:"uid"`
	ProfileID  int64  `json:"profile_id"`
	BucketName string `json:"bucket_name"`
	ObjectKey  string `json:"object_key"`
	Loaded     int64  `json:"loaded"`
	Total      int64  `json:"total"`
	Finished   bool   `json:"finished"`
	Error      error  `json:"error"`
}

// IController defines the interface of UIController
type IController interface {
	gousu.IUIController
}

// Controller is the core controller running the ui
type Controller struct {
	log             *gousu.Log
	s3Service       s3.IService
	settingsService settings.IService
	view            webview.WebView
	downloads       map[string]*DownloadStatus
}

var _ (IController) = (*Controller)(nil)

// Name returns the name of the controller defined by ControllerName
func (c *Controller) Name() string {
	return ControllerName
}

// Start creates a new window making it ready for being displayed
func (c *Controller) Start() error {
	c.downloads = map[string]*DownloadStatus{}

	c.view = webview.New(*debug)
	c.view.SetTitle("S3 Explorer")
	utils.SetWindowIcon(c.view.Window())
	c.view.SetSize(800, 600, webview.HintNone)
	c.view.Bind("s3DownloadFile", c.downloadFile)
	c.view.Bind("s3DownloadStatus", c.getDownloadStatus)
	c.view.Bind("systemOpenLink", c.systemOpenLink)
	c.view.Bind("showLicense", c.showLicense)
	c.view.Navigate(fmt.Sprintf("http://127.0.0.1:%d/index.html", *web.ServerPort))

	return nil
}

// Health always returns nil (healthy)
func (c *Controller) Health() error {
	return nil
}

// Run is the core routine running the ui
func (c *Controller) Run(chan os.Signal) error {
	// Sleep for 1 second to make sure the web server is up
	time.Sleep(1 * time.Second)

	c.view.Run()

	return nil
}

// Stop destroys the window causing the program to exit
func (c *Controller) Stop() error {
	c.view.Destroy()

	return nil
}

// NewController creates a new instance of the ui controller
func NewController(ctx gousu.IContext) gousu.IUIController {
	return &Controller{
		log:             gousu.GetLogger(fmt.Sprintf("controller.%s", ControllerName)),
		s3Service:       ctx.GetService(s3.ServiceName).(s3.IService),
		settingsService: ctx.GetService(settings.ServiceName).(settings.IService),
	}
}

var _ (gousu.UIControllerFactory) = NewController
