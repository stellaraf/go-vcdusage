package vcdusage_test

import (
	"testing"

	"github.com/stellaraf/go-vcdusage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			assert.NotZero(t, cores, "core count zero")
			assert.NotZero(t, mem.Float64(), "memory zero")
			assert.NotZero(t, stor.Float64(), "storage zero")
			assert.Equal(t, Env.Cores, cores, "mismatching core count: %v != %v", Env.Cores, cores)
			assert.Equal(t, Env.Memory, mem.GB(), "mismatching memory: %v != %v", Env.Memory, mem.GB())
			assert.Equal(t, Env.Storage, stor.GB(), "mismatching storage: %v != %v", Env.Storage, stor.GB())
			assert.Equal(t, Env.VMCount, vmCount, "mismatching VM count: %v != %v", Env.VMCount, vmCount)
		})
	}
}