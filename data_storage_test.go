package vcdusage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.stellar.af/go-vcdusage"
)

func Test_DataStorage(t *testing.T) {
	t.Parallel()
	base := vcdusage.DataStorage(549_755_813_888)
	assert.Equal(t, float64(549_755_813_888), base.Float64(), "float64 mismatch")
	assert.Equal(t, int64(549_755_813_888), base.Int64(), "int64 mismatch")
	assert.Equal(t, uint64(549_755_813_888), base.Uint64(), "uint64 mismatch")
	assert.Equal(t, "549755813888", base.String(), "String mismatch")
	assert.Equal(t, float64(536_870_912), base.KB(), "KB mismatch")
	assert.Equal(t, float64(524_288), base.MB(), "MB mismatch")
	assert.Equal(t, float64(512), base.GB(), "GB mismatch")
	assert.Equal(t, float64(0.5), base.TB(), "TB mismatch")
}
