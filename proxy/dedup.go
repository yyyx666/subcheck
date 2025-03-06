package proxies

import (
	"fmt"
)

func DeduplicateProxies(proxies []map[string]any) []map[string]any {
	seen := make(map[string]map[string]any)

	for _, proxy := range proxies {
		server, _ := proxy["server"].(string)
		port, _ := proxy["port"].(int)
		if server == "" {
			continue
		}

		key := fmt.Sprintf("%s:%v", server, port)
		seen[key] = proxy
	}

	result := make([]map[string]any, 0, len(seen))
	for _, proxy := range seen {
		result = append(result, proxy)
	}

	return result
}
