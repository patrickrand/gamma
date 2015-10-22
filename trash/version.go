package trash

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major uint8 `json:"major"`
	Minor uint8 `json:"minor"`
	Patch uint8 `json:"patch"`
}

func NewVersion(major, minor, patch uint8) *Version {
	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func ParseSemVer(semver string) (*Version, error) {
	vers := strings.Split(semver, ".")
	if len(vers) != 3 {
		return nil, errors.New("agent.ParseSemVer: invalid octet count")
	}

	octets := make([]uint8, 0)
	for i := range vers {
		v, err := strconv.ParseUint(vers[i], 10, 8)
		if err != nil {
			return nil, errors.New("agent.ParseSemVer: semantic version can only contain octets")
		}
		octets[i] = uint8(v)
	}

	return &Version{octets[0], octets[1], octets[2]}, nil
}
