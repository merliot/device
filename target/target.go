package target

import (
	"sort"
	"strings"
)

type GpioPin int
type GpioPins map[string]GpioPin

type Target struct {
	FullName string
	GpioPins
}

type Targets map[string]Target

func MakeTargets(targets []string) Targets {
	filtered := make(Targets)
	for _, target := range targets {
		if value, ok := supportedTargets[target]; ok {
			filtered[target] = value
		}
	}
	return filtered
}

func (targets Targets) FullNames() string {
	var fullNames []string
	for _, t := range targets {
		fullNames = append(fullNames, t.FullName)
	}
	// Sort FullNames alpha-numeric
	sort.Strings(fullNames)
	// Concatenate FullNames with commas
	return strings.Join(fullNames, ", ")
}
