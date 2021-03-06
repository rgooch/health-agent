package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Cloud-Foundations/Dominator/lib/fsutil"
	"github.com/Cloud-Foundations/Dominator/lib/log"
	libprober "github.com/Cloud-Foundations/health-agent/lib/prober"
	"github.com/Cloud-Foundations/health-agent/lib/proberlist"
	dialprober "github.com/Cloud-Foundations/health-agent/probers/dial"
	dnsprober "github.com/Cloud-Foundations/health-agent/probers/dns"
	ldapprober "github.com/Cloud-Foundations/health-agent/probers/ldap"
	pidprober "github.com/Cloud-Foundations/health-agent/probers/pidfile"
	testprogprober "github.com/Cloud-Foundations/health-agent/probers/testprog"
	urlprober "github.com/Cloud-Foundations/health-agent/probers/url"
	"github.com/Cloud-Foundations/tricorder/go/tricorder"
	"github.com/Cloud-Foundations/tricorder/go/tricorder/units"
	"gopkg.in/yaml.v2"
)

type testConfig struct {
	Testtype  string `yaml:"type"`
	Probefreq uint8  `yaml:"probe-freq"`
	Specs     testSpecs
}

type testSpecs struct {
	Address    string `yaml:"address"`
	Hostname   string `yaml:"hostname"`
	Network    string `yaml:"network"`
	Pathname   string
	SssdConfig string `yaml:"sssd-config"`
	Urlpath    string `yaml:"url-path"`
	Urlport    uint   `yaml:"url-port"`
}

func setupHealthchecks(configDir string, pl *proberlist.ProberList,
	logger log.Logger) error {
	topDir := "/health-checks"
	configdir, err := os.Open(path.Join(configDir, "tests.d"))
	defer configdir.Close()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	configfiles, err := configdir.Readdir(0)
	if err != nil {
		return err
	}
	healthCheckers := make(map[string]libprober.HealthChecker)
	for _, configfile := range configfiles {
		if configfile.IsDir() {
			continue
		}
		data, err := ioutil.ReadFile(path.Join(configdir.Name(),
			configfile.Name()))
		if err != nil {
			logger.Printf("Unable to read file %q: %s",
				configfile.Name(), err)
			return err
		}
		c := testConfig{}
		if err := yaml.Unmarshal([]byte(data), &c); err != nil {
			logger.Printf("Error unmarshalling file %s: %q",
				configfile.Name(), err)
			return err
		}
		testname := strings.Split(configfile.Name(), ".")[0]
		if prober := makeProber(testname, &c, logger); prober != nil {
			pl.Add(prober, path.Join(topDir, testname), c.Probefreq)
			if hc, ok := prober.(libprober.HealthChecker); ok {
				healthCheckers[testname] = hc
			}
		}
	}
	var allHealthy bool
	list := tricorder.NewList([]string{}, tricorder.ImmutableSlice)
	group := tricorder.NewGroup()
	err = group.RegisterMetric(path.Join(topDir, "*", "healthy"), &allHealthy,
		units.None, "If true, all health checks are healthy")
	if err != nil {
		return err
	}
	err = tricorder.RegisterMetricInGroup(path.Join(topDir, "*",
		"unhealthy-list"), list, group, units.None,
		"List of failed health checks")
	if err != nil {
		return err
	}
	group.RegisterUpdateFunc(func() time.Time {
		healthy := true
		var unhealthyList []string
		for testname, healthChecker := range healthCheckers {
			if !healthChecker.HealthCheck() {
				healthy = false
				unhealthyList = append(unhealthyList, testname)
			}
		}
		allHealthy = healthy
		list.Change(unhealthyList, tricorder.ImmutableSlice)
		return time.Now()
	})
	return nil
}

func makeProber(testname string, c *testConfig,
	logger log.Logger) proberlist.RegisterProber {
	switch c.Testtype {
	case "dial":
		return dialprober.MakeDialProber(c.Specs.Network, c.Specs.Address)
	case "dns":
		hostname := c.Specs.Hostname
		return dnsprober.New(testname, hostname)
	case "ldap":
		sssd := c.Specs.SssdConfig
		file, err := fsutil.WaitFile(sssd, time.Minute*5)
		if err != nil {
			logger.Println(err)
			return nil
		}
		return ldapprober.Makeldapprober(testname, file, c.Probefreq)
	case "pid":
		pidpath := c.Specs.Pathname
		if pidpath == "" {
			return nil
		}
		return pidprober.Makepidprober(testname, pidpath)
	case "testprog":
		testprogpath := c.Specs.Pathname
		if testprogpath == "" {
			return nil
		}
		return testprogprober.Maketestprogprober(testname, testprogpath)
	case "url":
		urlpath := c.Specs.Urlpath
		urlport := c.Specs.Urlport
		return urlprober.Makeurlprober(testname, urlpath, urlport)
	default:
		logger.Printf("Test type %s not supported", c.Testtype)
		return nil
	}
}
