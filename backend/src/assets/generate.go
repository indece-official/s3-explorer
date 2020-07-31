// +build ignore

package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var cwd, _ = os.Getwd()
	assets := http.Dir(filepath.Join(cwd, "../../assets"))

	err := vfsgen.Generate(assets, vfsgen.Options{
		Filename:     "../assets/assets_vfsdata.gen.go",
		PackageName:  "assets",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
