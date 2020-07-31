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

package utils

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
//
// static char* _ ( char* s )
// {
//     return s;
// }
//
// static void setWindowIcon(void* window, void* iconData, int iconDataLength)
// {
//     GdkPixbufLoader *loader = gdk_pixbuf_loader_new();
//     gdk_pixbuf_loader_write(loader, iconData, iconDataLength, NULL);
//     GdkPixbuf* pixbuf = gdk_pixbuf_loader_get_pixbuf(loader);
//
//     gtk_window_set_icon(window, pixbuf);
// }
//
// static char* saveFile ( void* parentWindow, char* defaultFilename )
// {
//     GtkWidget *dialog;
//     GtkFileChooser *chooser;
//     GtkFileChooserAction action = GTK_FILE_CHOOSER_ACTION_SAVE;
//     gint res;
//
//     dialog = gtk_file_chooser_dialog_new ("Save File",
//                                           parentWindow,
//                                           action,
//                                           _("_Cancel"),
//                                           GTK_RESPONSE_CANCEL,
//                                           _("_Save"),
//                                           GTK_RESPONSE_ACCEPT,
//                                           NULL);
//     chooser = GTK_FILE_CHOOSER(dialog);
//
//     gtk_file_chooser_set_do_overwrite_confirmation(chooser, TRUE);
//
//     gtk_file_chooser_set_current_name(chooser, defaultFilename);
//
//     res = gtk_dialog_run (GTK_DIALOG (dialog));
//     char *filename = 0;
//     if (res == GTK_RESPONSE_ACCEPT)
//     {
//         filename = gtk_file_chooser_get_filename(chooser);
//     }
//
//     gtk_widget_destroy(dialog);
//     return filename;
// }
import "C"
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"unsafe"

	"github.com/indece-official/s3-explorer/backend/src/assets"
)

// Init is a dummy function called before the window gets created (only required on windows)
func Init() {
}

// SetWindowIcon sets the windows icon loading it from the resources
func SetWindowIcon(win unsafe.Pointer) error {
	iconData, err := assets.ReadFileBinary("icon.ico")
	if err != nil {
		return err
	}

	C.setWindowIcon(win, C.CBytes(iconData), C.int(len(iconData)))

	return nil
}

// SaveFile opens a native Save-File-As dialog and returns the selected filename
// If the dialog is canceled an empty string ("") is returned as filename
func SaveFile(win unsafe.Pointer, defaultFilename string) (string, error) {
	cDefaultFilename := C.CString(defaultFilename)

	cFilename := C.saveFile(win, cDefaultFilename)

	filename := C.GoString(cFilename)

	C.free(unsafe.Pointer(cDefaultFilename))
	C.free(unsafe.Pointer(cFilename))

	return filename, nil
}

// OpenLink opens an url of the systems default browser
func OpenLink(url string) error {
	return exec.Command("xdg-open", url).Start()
}

// GetLicensePath returns the path of the LICENSE.txt
func GetLicensePath() (string, error) {
	// Check if license file is in current binary path (when not installed as package)
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("I don't exist?")
	}
	exPath := filepath.Dir(ex)

	licensePath := filepath.Join(exPath, "LICENSE.txt")
	if _, err := os.Stat(licensePath); err == nil {
		return licensePath, nil
	}

	return "/usr/share/doc/s3-explorer/copyright", nil
}
