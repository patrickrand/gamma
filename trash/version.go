package trash

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//func (cfg Config) Validate() error {
//	var errs Errors
//
//	if !regexp.MustCompile("^[a-z][a-z_-]*[a-z]$").MatchString(cfg.AgentName) {
//		errs = append(errs, errors.New("engine.Config.Validate: `name` must match `^[a-z][a-z-]*[a-z]$`"))
//	}
//
//	vers := strings.Split(cfg.AgentVersion, ".")
//	if len(vers) != 3 {
//		errs = append(errs, errors.New("engine.Config.ValidateVersion: invalid octet count"))
//	}
//
//	for i := range vers {
//		_, err := strconv.ParseUint(vers[i], 10, 8)
//		if err != nil {
//			errs = append(errs, errors.New("engine.Config.ValidateVersion: semantic version can only contain octets"))
//			break
//		}
//	}
//
//	if len(errs) == 0 {
//		return nil
//	}
//	return errs
//}

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
