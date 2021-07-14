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

package utils

// #cgo LDFLAGS: -lcomdlg32
//
// #include <stdlib.h>
// #include <windows.h>
// #include <commdlg.h>
// #include "../main/resources.h"
//
// static void setWindowIcon ( void *win )
// {
//     SendMessageW(win, WM_SETICON, ICON_BIG, (LPARAM) LoadIconW(GetModuleHandleW(NULL), MAKEINTRESOURCEW(IDI_ICON_1)));
//     SendMessageW(win, WM_SETICON, ICON_SMALL, (LPARAM) LoadIconW(GetModuleHandleW(NULL), MAKEINTRESOURCEW(IDI_ICON_1)));
// }
//
// static char* saveFile ( void* parentWindow, char* defaultFilename )
// {
//     char* filename = (char*) malloc(MAX_PATH+1 * sizeof(char));
//
//     if ( defaultFilename != NULL )
//     {
//         strcpy_s(filename, MAX_PATH+1, defaultFilename);
//     }
//
//     LPOPENFILENAME ofn = (LPOPENFILENAME) malloc(sizeof(OPENFILENAME));
//
//     memset(ofn, 0, sizeof(OPENFILENAME));
//
//     ofn->lStructSize = sizeof(OPENFILENAME);
//     ofn->hwndOwner = parentWindow;
//     ofn->hInstance = (HINSTANCE) GetWindowLongPtr(parentWindow, GWLP_HINSTANCE);
//     ofn->lpstrFilter = (LPCTSTR) NULL;
//     ofn->lpstrCustomFilter = (LPSTR) NULL;
//     ofn->nMaxCustFilter = 0;
//     ofn->nFilterIndex = 0;
//     ofn->lpstrFile = filename;
//     ofn->nMaxFile = MAX_PATH+1;
//     ofn->lpstrFileTitle = NULL;
//     ofn->lpstrTitle = (LPCTSTR) NULL;
//     ofn->Flags = OFN_EXPLORER;
//
//     BOOL result = GetSaveFileName(ofn);
//     if ( result == FALSE )
//     {
//         strcpy_s(filename, MAX_PATH+1, "");
//     }
//
//     free(ofn);
//
//     return filename;
// }
import "C"

import (
	"fmt"
	"github.com/indece-official/go-gousu"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

// Init is called before the window is created and fixes "blury"-issues on high-res screens
func Init() {
	modshcore := syscall.NewLazyDLL("shcore.dll")
	shc := modshcore.NewProc("SetProcessDpiAwareness")
	shc.Call(uintptr(1))

	gousu.DisableLogger()
}

// SetWindowsIcons sets the window's icon loading it from the compiled in resources
func SetWindowIcon(win unsafe.Pointer) error {
	C.setWindowIcon(win)

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
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
}

// GetLicensePath returns the path of the LICENSE.txt
func GetLicensePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("I don't exist?!")
	}
	exPath := filepath.Dir(ex)

	return filepath.Join(exPath, "LICENSE.txt"), nil
}
