package m

import (
	"embed"
	"github.com/Vilsol/go-mlog/checker"
)

//go:embed *.go
var PackageM embed.FS

func init() {
	dir, err := PackageM.ReadDir(".")

	if err != nil {
		panic(err)
	}

	files := make(map[string]string)
	for _, entry := range dir {
		if !entry.IsDir() {
			file, err := PackageM.ReadFile(entry.Name())
			if err != nil {
				panic(err)
			}
			files[entry.Name()] = string(file)
		}
	}

	checker.RegisterPackages("github.com/Vilsol/go-mlog/m", files)
}
