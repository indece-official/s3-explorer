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
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/indece-official/s3-explorer/backend/src/assets"
)

func (c *Controller) reqStaticFile(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimLeft(r.RequestURI, "/")

	filename = strings.Replace(filename, "../", "", -1)
	filename = strings.Replace(filename, "./", "", -1)

	if filename == "" {
		filename = "index.html"
	}

	file, err := assets.Assets.Open(path.Join("/www", filename))
	if err != nil {
		filename = "index.html"

		file, err = assets.Assets.Open(path.Join("/www", filename))
		if err != nil {
			c.log.Infof("GET /%s - 404 Not Found (%s)", filename, err)

			w.WriteHeader(404)
			fmt.Fprintf(w, "File not found")
			return
		}
	}
	defer file.Close()

	contentType := "application/octet-stream"

	switch path.Ext(filename) {
	case ".js":
		contentType = "text/javascript"
	case ".css":
		contentType = "text/css"
	case ".html":
		contentType = "text/html"
	case ".ico":
		contentType = "image/x-icon"
	}

	// Get the file size
	stat, _ := file.Stat()                     // Get info from file
	size := strconv.FormatInt(stat.Size(), 10) // Get file size as a string

	c.log.Infof("GET /%s - 200 OK", filename)

	// Send the headers
	//w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", size)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Send the file
	// We read 512 bytes from the file already, so we reset the offset back to 0
	file.Seek(0, 0)
	io.Copy(w, file) // 'Copy' the file to the client
}
