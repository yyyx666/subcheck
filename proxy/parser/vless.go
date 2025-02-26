package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// 将vless格式的节点转换为clash的节点
func ParseVless(data string) (map[string]any, error) {
	parsedURL, err := url.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("解析失败: %v", err)
	}

	if parsedURL.Scheme != "vless" {
		return nil, fmt.Errorf("不是vless格式")
	}

	hostPort := strings.Split(parsedURL.Host, ":")
	if len(hostPort) != 2 {
		return nil, nil
	}

	host := hostPort[0]
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return nil, fmt.Errorf("格式错误: 端口格式不正确")
	}

	// 解析参数
	params, err := url.ParseQuery(parsedURL.RawQuery)
	if err != nil {
		return nil, nil
	}

	// 构建 clash 格式的代理配置
	proxy := map[string]any{
		"name":        parsedURL.Fragment,
		"type":        "vless",
		"server":      host,
		"port":        port,
		"uuid":        parsedURL.User.String(),
		"network":     params.Get("type"),
		"tls":         params.Get("security") == "tls",
		"servername":  params.Get("sni"),
		"flow":        params.Get("flow"),
		"fingerprint": params.Get("fp"),
	}

	// 添加 ws 特定配置
	if params.Get("type") == "ws" {
		wsOpts := map[string]any{
			"path": params.Get("path"),
			"headers": map[string]any{
				"Host": params.Get("host"),
			},
		}
		proxy["ws-opts"] = wsOpts
	}
	realityOpts := map[string]any{
		"public-key": params.Get("pbk"),
		"short-id":   params.Get("sid"),
	}
	proxy["reality-opts"] = realityOpts

	return proxy, nil
}
