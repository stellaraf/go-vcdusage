package vcdusage

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/joomcode/errorx"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type Client struct {
	VCD *govcd.VCDClient
}

type Options struct {
	Insecure bool
	Org      string
	Username string
	Password string
	URL      *url.URL
}

// ErrCfgNoUsername indicates an authentication username was not provided.
var ErrCfgNoUsername = errors.New("username required")

// ErrCfgNoPassword indicates an authentication password was not provided.
var ErrCfgNoPassword = errors.New("password required")

// ErrCfgNoURL indicates a URL was not provided.
var ErrCfgNoURL = errors.New("URL required")

// Validate ensures required parameters are set.
func (opts *Options) Validate() error {
	if opts.Username == "" {
		return ErrCfgNoUsername
	}
	if opts.Password == "" {
		return ErrCfgNoPassword
	}
	if opts.URL == nil {
		return ErrCfgNoURL
	}
	return nil
}

type Option func(*Options)

// Insecure disables SSL certificate validation when communicating with vCloud.
func Insecure() Option {
	return func(opts *Options) {
		opts.Insecure = true
	}
}

// Org sets the organization name used when authenticating with vCloud. If not set, 'system' will
// be used.
func Org(org string) Option {
	return func(opts *Options) {
		opts.Org = org
	}
}

// Username sets the authentication username. This option is required.
func Username(u string) Option {
	return func(opts *Options) {
		opts.Username = u
	}
}

// Password sets the authentication password. This option is required.
func Password(p string) Option {
	return func(opts *Options) {
		opts.Password = p
	}
}

// URL sets the vCloud base URL. This option is required. See ParseURL helper.
func URL(u *url.URL) Option {
	return func(opts *Options) {
		opts.URL = u
	}
}

// ParseURL parses a vCloud URL from a string to a *url.URL, sets the appropriate URI schema, and
// sets the correct path.
func ParseURL(u string) (*url.URL, error) {
	if !strings.HasPrefix(u, "http") {
		u = fmt.Sprintf("https://%s", u)
	}
	pu, err := url.Parse(u)
	if err != nil {
		err = errorx.Decorate(err, "failed to parse URL '%s'", u)
		return nil, err
	}
	pu.Scheme = "https"
	if pu.Path == "/" || pu.Path == "" {
		pu.Path = "/api"
	}
	return pu, nil
}

// Create a new VCD Usage client.
func New(options ...Option) (*Client, error) {
	opts := &Options{
		Insecure: false,
		Org:      "system",
	}
	for _, setter := range options {
		setter(opts)
	}
	err := opts.Validate()
	if err != nil {
		err = errorx.Decorate(err, "option validation failed")
		return nil, err
	}
	vcd := govcd.NewVCDClient(*opts.URL, opts.Insecure)
	err = vcd.Authenticate(opts.Username, opts.Password, opts.Org)
	if err != nil {
		err = errorx.Decorate(err, "failed to authenticate with vCloud host '%s'", opts.URL.String())
		return nil, err
	}
	client := &Client{
		VCD: vcd,
	}
	return client, nil
}
