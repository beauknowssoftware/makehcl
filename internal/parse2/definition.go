package parse2

import (
	"github.com/hashicorp/hcl/v2"
)

type Definition struct {
	Files map[string]*File
}

func (d *Definition) addFile(name string, f *hcl.File) *File {
	if d.Files == nil {
		d.Files = make(map[string]*File)
	}

	newFile := &File{
		Name:    name,
		hclFile: f,
	}
	d.Files[name] = newFile

	return newFile
}
