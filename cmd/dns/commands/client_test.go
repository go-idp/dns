package commands

import "testing"

func TestPlainDNSAddressForStress(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"127.0.0.1", "127.0.0.1:53", false},
		{"127.0.0.1:5353", "127.0.0.1:5353", false},
		{"tls://1.1.1.1", "", true},
		{"https://cloudflare-dns.com/dns-query", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			got, err := plainDNSAddressForStress(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}
