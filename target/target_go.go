//go:build !tinygo

package target

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// GenTargetJS generates JS const deployTargetGpios
func GenTargetJS(name string) {
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("const deployTargetGpios = {\n")
	for k, v := range AllTargets {
		gpioKeys := make([]string, 0, len(v.GpioPins))
		for gpioKey := range v.GpioPins {
			gpioKeys = append(gpioKeys, fmt.Sprintf(`"%s"`, gpioKey))
		}
		sort.Slice(gpioKeys, func(i, j int) bool {
			numA := extractNumber(gpioKeys[i])
			numB := extractNumber(gpioKeys[j])
			return numA < numB
		})

		fmt.Fprintf(file, "    \"%s\": [%s],\n", k, strings.Join(gpioKeys, ", "))
	}
	file.WriteString("};\n")
}

// extractNumber extracts the numeric part of a string and converts it to an integer
func extractNumber(s string) int {
	numStr := strings.TrimFunc(s, func(r rune) bool {
		return r < '0' || r > '9'
	})
	num, _ := strconv.Atoi(numStr)
	return num
}
