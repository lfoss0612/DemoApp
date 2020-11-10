package util

import (
	"os"
)

func GetMachineID() string {
	if dockerShortID, err := DockerShortID(); err == nil {
		return dockerShortID
	} else if name, hostnameErr := os.Hostname(); hostnameErr == nil {
		return name
	}
	return ""
}
