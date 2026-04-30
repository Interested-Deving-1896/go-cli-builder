package resolver

import (
	"os"
)

// GetValue returns the value to bind based on priority: CLI > Env > Default.
//
// Example:
//
//	val, err := resolver.GetValue(nil, "ENV_VAR", "default", false)
func GetValue(cliVal *string, envName, defVal string, isFlag bool) (string, error) {
	if cliVal != nil {
		return *cliVal, nil
	}

	if envName != "" {
		if val, ok := os.LookupEnv(envName); ok {
			return val, nil
		}
	}

	if defVal != "" {
		return defVal, nil
	}

	return "", nil
}