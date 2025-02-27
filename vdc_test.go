package vcdusage_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.stellar.af/go-vcdusage"
)

func Test_VDCs(t *testing.T) {
	u, err := vcdusage.ParseURL(Env.URL)
	require.NoError(t, err)
	client, err := vcdusage.New(
		vcdusage.Insecure(),
		vcdusage.URL(u),
		vcdusage.Username(Env.Username),
		vcdusage.Password(Env.Password),
	)
	require.NoError(t, err)
	vdcs, err := client.VDCs(Env.OrgID)
	require.NoError(t, err)
	for _, vdc := range vdcs {
		vdc := vdc
		t.Run(vdc.Obj.Vdc.Name, func(t *testing.T) {
			t.Parallel()
			cores := vdc.CoreCount()
			mem := vdc.Memory()
			stor := vdc.Storage()
			vmCount := vdc.VMCount()
			poweredOn := vdc.PoweredOnVMCount()
			assert.NotZero(t, cores, "core count zero")
			assert.NotZero(t, mem.Float64(), "memory zero")
			assert.NotZero(t, stor.Float64(), "storage zero")
			assert.Equal(t, Env.Cores, cores, "mismatching core count: %v != %v", Env.Cores, cores)
			assert.Equal(t, Env.Memory, mem.GB(), "mismatching memory: %v != %v", Env.Memory, mem.GB())
			assert.Equal(t, Env.Storage, stor.GB(), "mismatching storage: %v != %v", Env.Storage, stor.GB())
			total := Env.VMCountOn + Env.VMCountOff
			assert.Equal(t, total, vmCount, "mismatching VM count: %v != %v", total, vmCount)
			assert.Equal(t, Env.VMCountOn, poweredOn, "mismatching powered-on VM count: %v != %v", Env.VMCountOn, poweredOn)
		})
	}
	t.Run("individual VDC", func(t *testing.T) {
		t.Parallel()
		vdc, err := client.VDC(Env.OrgID, Env.VdcID)
		require.NoError(t, err)
		cores := vdc.CoreCount()
		mem := vdc.Memory()
		stor := vdc.Storage()
		vmCount := vdc.VMCount()
		poweredOn := vdc.PoweredOnVMCount()
		assert.NotZero(t, cores, "core count zero")
		assert.NotZero(t, mem.Float64(), "memory zero")
		assert.NotZero(t, stor.Float64(), "storage zero")
		assert.Equal(t, Env.Cores, cores, "mismatching core count: %v != %v", Env.Cores, cores)
		assert.Equal(t, Env.Memory, mem.GB(), "mismatching memory: %v != %v", Env.Memory, mem.GB())
		assert.Equal(t, Env.Storage, stor.GB(), "mismatching storage: %v != %v", Env.Storage, stor.GB())
		total := Env.VMCountOn + Env.VMCountOff
		assert.Equal(t, total, vmCount, "mismatching VM count: %v != %v", total, vmCount)
		assert.Equal(t, Env.VMCountOn, poweredOn, "mismatching powered-on VM count: %v != %v", Env.VMCountOn, poweredOn)
	})
	t.Run("all VDCs", func(t *testing.T) {
		t.Parallel()
		cores := vdcs.CoreCount()
		mem := vdcs.Memory()
		stor := vdcs.Storage()
		vmCount := vdcs.VMCount()
		poweredOn := vdcs.PoweredOnVMCount()
		assert.NotZero(t, cores, "core count zero")
		assert.NotZero(t, mem.Float64(), "memory zero")
		assert.NotZero(t, stor.Float64(), "storage zero")
		assert.Equal(t, Env.Cores, cores, "mismatching core count: %v != %v", Env.Cores, cores)
		assert.Equal(t, Env.Memory, mem.GB(), "mismatching memory: %v != %v", Env.Memory, mem.GB())
		assert.Equal(t, Env.Storage, stor.GB(), "mismatching storage: %v != %v", Env.Storage, stor.GB())
		total := Env.VMCountOn + Env.VMCountOff
		assert.Equal(t, total, vmCount, "mismatching VM count: %v != %v", total, vmCount)
		assert.Equal(t, Env.VMCountOn, poweredOn, "mismatching powered-on VM count: %v != %v", Env.VMCountOn, poweredOn)
	})
	t.Run("guest os", func(t *testing.T) {
		t.Parallel()
		vdc, err := client.VDC(Env.OrgID, Env.VdcID)
		require.NoError(t, err)
		count := vdc.VMCountWithQuery(vcdusage.VMWithGuestOSContaining("windows"))
		assert.NotZero(t, count, "no VMs matching query")
		assert.Equal(t, Env.WindowsCount, count, "mismatching Windows VM count: %v != %v", Env.WindowsCount, count)

		count = vdc.VMCountWithQuery(vcdusage.VMPoweredOn(), vcdusage.VMWithGuestOSMatching(regexp.MustCompile(".*Windows.*")))
		assert.NotZero(t, count, "no VMs matching query")
		assert.Equal(t, Env.WindowsCountOn, count, "mismatching Windows VM count: %v != %v", Env.WindowsCountOn, count)
	})
}
