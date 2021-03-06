package packages

import (
	"github.com/Cloud-Foundations/tricorder/go/tricorder"
)

func register(dir *tricorder.DirectorySpec) *prober {
	p := &prober{
		dir:       dir,
		packagers: make(map[string]*packageList),
	}
	return p
}
