package kernel

import (
	"bufio"
	"os"

	"github.com/Cloud-Foundations/tricorder/go/tricorder"
	"github.com/Cloud-Foundations/tricorder/go/tricorder/units"
)

var versionFilename string = "/proc/version"

func register(dir *tricorder.DirectorySpec) *prober {
	p := new(prober)
	if err := dir.RegisterMetric("hostname", &p.hostname, units.None,
		"system hostname"); err != nil {
		panic(err)
	}
	version := getVersion()
	if err := dir.RegisterMetric("version/raw", &version, units.None,
		"raw kernel version"); err != nil {
		panic(err)
	}
	if err := dir.RegisterMetric("random/entropy-available",
		&p.randomEntropyAvailable, units.Byte,
		"entropy available for the random number generator"); err != nil {
		panic(err)
	}
	return p
}

func getVersion() string {
	file, err := os.Open(versionFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	version, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return version
}
