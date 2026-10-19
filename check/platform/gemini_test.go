package platform

import "testing"

func TestGeminiAvailable(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{"new format available", `false,\"45723056\":true,\"45631641\":true,\"45841483\":false`, true},
		{"new format unavailable", `false,\"45723056\":true,\"45631641\":false,\"45841483\":false`, false},
		{"unescaped json", `{"45631641" : true}`, true},
		{"legacy format", `45631641,null,true`, true},
		{"marker missing", `some other page`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := geminiAvailable([]byte(tt.body)); got != tt.want {
				t.Fatalf("geminiAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
