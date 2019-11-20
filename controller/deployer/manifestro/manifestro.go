package manifestro

import (
	"github.com/cloudfoundry-incubator/candiedyaml"
)

type manifestYaml struct {
	Applications []struct {
		Instances *uint16
	}
}

// GetInstances reads a Cloud Foundry manifest as a string and returns the number of instances
// defined in the manifest, if there are any.
//
// Returns a point to a uint16. If instances are not found or less than 1, it returns nil.
func GetInstances(manifest string) *uint16 {
	var m manifestYaml

	err := candiedyaml.Unmarshal([]byte(manifest), &m)
	if err != nil || m.Applications == nil || m.Applications[0].Instances == nil || *m.Applications[0].Instances < 1 {
		return nil
	}

	return m.Applications[0].Instances
}
