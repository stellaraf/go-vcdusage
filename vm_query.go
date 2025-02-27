package vcdusage

import (
	"fmt"
	"regexp"
)

type VMQuery struct {
	Name      *regexp.Regexp
	GuestOS   *regexp.Regexp
	PoweredOn bool
}

type VMQuerySetter func(*VMQuery)

func VMWithNameContaining(contains string) VMQuerySetter {
	return func(q *VMQuery) {
		q.Name = regexp.MustCompile(fmt.Sprintf("(?i).*%s.*", contains))
	}
}

func VMWithNameMatching(pattern *regexp.Regexp) VMQuerySetter {
	return func(q *VMQuery) {
		q.Name = pattern
	}
}

func VMWithGuestOSContaining(contains string) VMQuerySetter {
	return func(q *VMQuery) {
		q.GuestOS = regexp.MustCompile(fmt.Sprintf("(?i).*%s.*", contains))
	}
}

func VMWithGuestOSMatching(pattern *regexp.Regexp) VMQuerySetter {
	return func(q *VMQuery) {
		q.GuestOS = regexp.MustCompile("(?i)" + pattern.String())
	}
}

func VMPoweredOn() VMQuerySetter {
	return func(q *VMQuery) {
		q.PoweredOn = true
	}
}
