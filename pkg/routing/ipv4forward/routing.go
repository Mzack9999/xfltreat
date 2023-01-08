package routing

import (
	"github.com/pkg/errors"

	"github.com/Mzack9999/xfltreat/pkg/sysctl"
)

func IsEnabled() (bool, error) {
	ipv4Forward, err := sysctl.Get(sysctl.Linux_IP4_Forward)
	if ipv4Forward == -1 {
		return false, errors.Wrap(err, "unexpected value")
	}
	return ipv4Forward == 1, err
}

func Set(value bool) error {
	var newValue int
	if value {
		newValue = 1
	} else {
		newValue = 0
	}
	return sysctl.Set(sysctl.Linux_IP4_Forward, newValue)
}
