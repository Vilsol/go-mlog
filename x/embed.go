package x

import (
	"embed"
	"github.com/Vilsol/go-mlog/checker"
)

//go:embed *.go
var PackageX embed.FS

func init() {
	dir, err := PackageX.ReadDir(".")

	if err != nil {
		panic(err)
	}

	files := make(map[string]string)
	for _, entry := range dir {
		if !entry.IsDir() {
			file, err := PackageX.ReadFile(entry.Name())
			if err != nil {
				panic(err)
			}
			files[entry.Name()] = string(file)
		}
	}

	checker.RegisterPackages("github.com/Vilsol/go-mlog/x", files)
}
