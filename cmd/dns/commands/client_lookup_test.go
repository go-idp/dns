package commands

import (
	"testing"
)

func TestParseLookupArgv(t *testing.T) {
	t.Parallel()
	o, err := parseLookupArgv([]string{"baidu.com", "--server", "223.5.5.5"})
	if err != nil {
		t.Fatal(err)
	}
	if o.domain != "baidu.com" || len(o.servers) != 1 || o.servers[0] != "223.5.5.5" {
		t.Fatalf("got %+v", o)
	}

	o2, err := parseLookupArgv([]string{"--server", "8.8.8.8", "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if o2.domain != "example.com" || len(o2.servers) != 1 || o2.servers[0] != "8.8.8.8" {
		t.Fatalf("got %+v", o2)
	}
}
