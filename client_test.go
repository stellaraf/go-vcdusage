package vcdusage_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.stellar.af/go-vcdusage"
)

func Test_New(t *testing.T) {
	u, err := vcdusage.ParseURL(Env.URL)
	require.NoError(t, err)
	_, err = vcdusage.New(
		vcdusage.Insecure(),
		vcdusage.URL(u),
		vcdusage.Username(Env.Username),
		vcdusage.Password(Env.Password),
	)
	require.NoError(t, err)
}

func Test_ParseURL(t *testing.T) {
	cases := []string{
		"vcd.example.com",
		"vcd.example.com/api",
		"http://vcd.example.com",
		"http://vcd.example.com/api",
		"https://vcd.example.com",
		"https://vcd.example.com/api",
	}
	final, err := url.Parse("https://vcd.example.com/api")
	require.NoError(t, err)
	for i := 0; i < len(cases); i++ {
		us := cases[i]
		t.Run(fmt.Sprintf("%d-%s", i+1, us), func(t *testing.T) {
			t.Parallel()
			u, err := vcdusage.ParseURL(us)
			require.NoError(t, err)
			assert.Equal(t, final, u)
		})
	}
}
