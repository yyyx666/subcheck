package proxies

import "testing"

func TestIsCloudflareIP(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"", false},
		{"not-an-ip", false},
		{"104.16.0.1", true},
		{"162.158.0.1", true},
		{"172.64.0.1", true},
		{"1.1.1.1", false}, // Cloudflare DNS 但不在官方代理网段列表的部分前缀外需以列表为准；1.1.1.1 实际属 1.1.1.0/24 不在 ips-v4
		{"8.8.8.8", false},
		{"192.168.1.1", false},
		{"2606:4700::1", true},
		{"2400:cb00::1", true},
		{"2001:4860:4860::8888", false},
	}
	for _, tt := range tests {
		if got := IsCloudflareIP(tt.ip); got != tt.want {
			t.Errorf("IsCloudflareIP(%q) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestIsCloudflareWorkersProxy(t *testing.T) {
	tests := []struct {
		name  string
		proxy map[string]any
		want  bool
	}{
		{
			name:  "server workers.dev",
			proxy: map[string]any{"server": "example.workers.dev"},
			want:  true,
		},
		{
			name:  "sni workers.dev",
			proxy: map[string]any{"server": "1.2.3.4", "sni": "foo.workers.dev"},
			want:  true,
		},
		{
			name:  "servername workers.dev",
			proxy: map[string]any{"server": "1.2.3.4", "servername": "bar.workers.dev"},
			want:  true,
		},
		{
			name: "ws-opts host",
			proxy: map[string]any{
				"server": "1.2.3.4",
				"ws-opts": map[string]any{
					"headers": map[string]any{"Host": "ws.workers.dev"},
				},
			},
			want: true,
		},
		{
			name:  "normal domain",
			proxy: map[string]any{"server": "example.com", "sni": "www.google.com"},
			want:  false,
		},
		{
			name:  "empty",
			proxy: map[string]any{},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCloudflareWorkersProxy(tt.proxy); got != tt.want {
				t.Errorf("IsCloudflareWorkersProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCloudflareOrganization(t *testing.T) {
	if !IsCloudflareOrganization("CLOUDFLARENET") {
		t.Error("expected CLOUDFLARENET to match")
	}
	if !IsCloudflareOrganization("Cloudflare, Inc.") {
		t.Error("expected Cloudflare, Inc. to match")
	}
	if IsCloudflareOrganization("") {
		t.Error("empty org should not match")
	}
	if IsCloudflareOrganization("Google LLC") {
		t.Error("Google should not match")
	}
}
