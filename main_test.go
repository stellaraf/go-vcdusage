package vcdusage_test

import (
	"os"

	"github.com/stellaraf/go-utils/environment"
)

type Environment struct {
	URL      string  `env:"VCD_URL"`
	Username string  `env:"VCD_USERNAME"`
	Password string  `env:"VCD_PASSWORD"`
	OrgID    string  `env:"VCD_ORG_ID"`
	Cores    uint64  `env:"VCD_CORES"`
	Memory   float64 `env:"VCD_MEMORY"`
	Storage  float64 `env:"VCD_STORAGE"`
	VMCount  uint64  `env:"VCD_VM_COUNT"`
	VdcID    string  `env:"VCD_VDC_ID"`
	OrgID2   string  `env:"VCD2_ORG_ID"`
	Cores2   uint64  `env:"VCD2_CORES"`
	Memory2  float64 `env:"VCD2_MEMORY"`
	Storage2 float64 `env:"VCD2_STORAGE"`
	VMCount2 uint64  `env:"VCD2_VM_COUNT"`
	VdcID2   string  `env:"VCD2_VDC_ID"`
}

var Env Environment

func init() {
	err := environment.Load(&Env, &environment.EnvironmentOptions{
		DotEnv: os.Getenv("CI") == "",
	})
	if err != nil {
		panic(err)
	}
}
