package memory

import (
	libprober "github.com/Cloud-Foundations/health-agent/lib/prober"
	"github.com/Cloud-Foundations/tricorder/go/tricorder"
	"io"
)

type prober struct {
	available uint64
	free      uint64
	total     uint64
}

func Register(dir *tricorder.DirectorySpec) libprober.Prober {
	return register(dir)
}

func (p *prober) Probe() error {
	return p.probe()
}

func (p *prober) WriteHtml(writer io.Writer) {
	p.writeHtml(writer)
}
