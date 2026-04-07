package directives

import (
	"strings"
)

// configStack manages cascading config scopes. Each scope level inherits
// from its parent and can override individual keys.
type configStack struct {
	scopes []map[string]string
}

// push adds a new empty scope level.
func (cs *configStack) push() {
	cs.scopes = append(cs.scopes, make(map[string]string))
}

// pop removes the innermost scope level.
func (cs *configStack) pop() {
	if len(cs.scopes) > 0 {
		cs.scopes = cs.scopes[:len(cs.scopes)-1]
	}
}

// set stores a key=value pair at the current (innermost) scope level.
// If no scopes exist, a root scope is created.
func (cs *configStack) set(key, value string) {
	if len(cs.scopes) == 0 {
		cs.scopes = append(cs.scopes, make(map[string]string))
	}
	cs.scopes[len(cs.scopes)-1][key] = value
}

// merged returns a flat map with all config keys, with deeper scopes
// overriding shallower ones.
func (cs *configStack) merged() map[string]string {
	result := make(map[string]string)
	for _, scope := range cs.scopes {
		for k, v := range scope {
			result[k] = v
		}
	}
	return result
}

// depth returns the number of active scopes.
func (cs *configStack) depth() int {
	return len(cs.scopes)
}

// parseConfigPairs parses "key=value, key=value, ..." into a map.
func parseConfigPairs(s string) map[string]string {
	result := make(map[string]string)
	if s == "" {
		return result
	}

	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		idx := strings.IndexByte(pair, '=')
		if idx == -1 {
			continue // malformed pair, skip
		}
		key := strings.TrimSpace(pair[:idx])
		value := strings.TrimSpace(pair[idx+1:])
		if key != "" {
			result[key] = value
		}
	}
	return result
}
