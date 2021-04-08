package checker

import (
	"fmt"
	"go/types"
)

type Importer struct {
	packages map[string]*types.Package
}

func New() *Importer {
	return &Importer{
		packages: make(map[string]*types.Package),
	}
}

var importing types.Package

func (p *Importer) Import(path string) (*types.Package, error) {
	pkg := p.packages[path]
	if pkg != nil {
		if pkg == &importing {
			return nil, fmt.Errorf("import cycle through package %q", path)
		}
		if !pkg.Complete() {
			return pkg, fmt.Errorf("reimported partially imported package %q", path)
		}
		return pkg, nil
	}

	p.packages[path] = &importing
	defer func() {
		if p.packages[path] == &importing {
			p.packages[path] = nil
		}
	}()

	var firstHardErr error
	conf := types.Config{
		IgnoreFuncBodies: true,
		Error: func(err error) {
			if firstHardErr == nil && !err.(types.Error).Soft {
				firstHardErr = err
			}
		},
		Importer: p,
	}

	pack := FindPackage(path)

	if pack == nil {
		return nil, fmt.Errorf("package not found %q", path)
	}

	pkg, err := conf.Check(path, GetFset(), pack, nil)
	if err != nil {
		if firstHardErr != nil {
			pkg = nil
			err = firstHardErr
		}
		return pkg, fmt.Errorf("type-checking package %q failed (%v)", path, err)
	}

	if firstHardErr != nil {
		panic("package is not safe yet no error was returned")
	}

	p.packages[path] = pkg
	return pkg, nil
}
