package proxies

import (
	"net"
	"strings"
)

// Cloudflare 官方 IP 段（https://www.cloudflare.com/ips-v4 + ips-v6）
// 网段变更极少，随发版更新即可。
var cloudflareCIDRs = []string{
	// IPv4
	"173.245.48.0/20",
	"103.21.244.0/22",
	"103.22.200.0/22",
	"103.31.4.0/22",
	"141.101.64.0/18",
	"108.162.192.0/18",
	"190.93.240.0/20",
	"188.114.96.0/20",
	"197.234.240.0/22",
	"198.41.128.0/17",
	"162.158.0.0/15",
	"104.16.0.0/13",
	"104.24.0.0/14",
	"172.64.0.0/13",
	"131.0.72.0/22",
	// IPv6
	"2400:cb00::/32",
	"2606:4700::/32",
	"2803:f800::/32",
	"2405:b500::/32",
	"2405:8100::/32",
	"2a06:98c0::/29",
	"2c0f:f248::/32",
}

var cloudflareIPNets []*net.IPNet

func init() {
	cloudflareIPNets = make([]*net.IPNet, 0, len(cloudflareCIDRs))
	for _, cidr := range cloudflareCIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		cloudflareIPNets = append(cloudflareIPNets, ipNet)
	}
}

// IsCloudflareIP 判断出口 IP 是否落在 Cloudflare 官方网段。
// 空字符串或非法 IP 返回 false（fail-open）。
func IsCloudflareIP(ipStr string) bool {
	if ipStr == "" {
		return false
	}
	ip := net.ParseIP(strings.TrimSpace(ipStr))
	if ip == nil {
		return false
	}
	for _, ipNet := range cloudflareIPNets {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

// IsCloudflareWorkersProxy 判断节点连接地址或 SNI 是否包含 workers.dev。
func IsCloudflareWorkersProxy(proxy map[string]any) bool {
	for _, v := range collectProxyHostnames(proxy) {
		if containsWorkersDev(v) {
			return true
		}
	}
	return false
}

func containsWorkersDev(host string) bool {
	h := strings.ToLower(strings.TrimSpace(host))
	if h == "" {
		return false
	}
	// 去掉可能的端口
	if hostOnly, _, err := net.SplitHostPort(h); err == nil {
		h = hostOnly
	}
	h = strings.Trim(h, "[]")
	return h == "workers.dev" || strings.HasSuffix(h, ".workers.dev")
}

// collectProxyHostnames 收集节点中可能出现的域名/SNI 字段。
func collectProxyHostnames(proxy map[string]any) []string {
	var hosts []string
	appendStr := func(v any) {
		if s, ok := v.(string); ok && s != "" {
			hosts = append(hosts, s)
		}
	}

	for _, key := range []string{"server", "servername", "sni", "host"} {
		appendStr(proxy[key])
	}

	// ws-opts.headers.Host
	if wsOpts, ok := proxy["ws-opts"].(map[string]any); ok {
		if headers, ok := wsOpts["headers"].(map[string]any); ok {
			appendStr(headers["Host"])
			appendStr(headers["host"])
		}
	}

	// grpc-opts.grpc-service-name
	if grpcOpts, ok := proxy["grpc-opts"].(map[string]any); ok {
		appendStr(grpcOpts["grpc-service-name"])
	}

	// reality-opts.server-name
	if realityOpts, ok := proxy["reality-opts"].(map[string]any); ok {
		appendStr(realityOpts["server-name"])
	}

	return hosts
}

// IsCloudflareOrganization 判断出口组织名是否为 Cloudflare。
func IsCloudflareOrganization(org string) bool {
	return org != "" && strings.Contains(strings.ToLower(org), "cloudflare")
}
