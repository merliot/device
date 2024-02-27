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
	TinyGo bool
}

type Targets map[string]Target

func MakeTargets(targets []string) Targets {
	filtered := make(Targets)
	for _, target := range targets {
		if value, ok := AllTargets[target]; ok {
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

func (targets Targets) TinyGoTargets() Targets {
	filtered := make(Targets)
	for key, target := range targets {
		if target.TinyGo {
			filtered[key] = target
		}
	}
	return filtered
}
