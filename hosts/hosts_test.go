package hosts_test

import (
	"net/netip"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/depot/orca/hosts"
)

func TestParse(t *testing.T) {
	example := `
255.255.255.255	broadcasthost
127.0.0.2	hydrogen
127.0.0.3	hydrogen  # inline comment
::2             hydrogen
127.1.1.1	helium
# aliases
127.1.1.2	lithium lithiumhost
fe80::1%lo0	localhost
# Bogus entries that must be ignored.
127.10.10.10 # no hostname
123.123.123	copper
321.321.321.321`

	actual, err := hosts.Parse(strings.NewReader(example))
	if err != nil {
		t.Fatal(err)
	}

	expected := &hosts.Hosts{
		Entries: []hosts.HostEntry{
			{
				Addr:  netip.MustParseAddr("255.255.255.255"),
				Hosts: []string{"broadcasthost"},
			},
			{
				Addr:  netip.MustParseAddr("127.0.0.2"),
				Hosts: []string{"hydrogen"},
			},
			{
				Addr:  netip.MustParseAddr("127.0.0.3"),
				Hosts: []string{"hydrogen"},
			},
			{
				Addr:  netip.MustParseAddr("::2"),
				Hosts: []string{"hydrogen"},
			},
			{
				Addr:  netip.MustParseAddr("127.1.1.1"),
				Hosts: []string{"helium"},
			},
			{
				Addr:  netip.MustParseAddr("127.1.1.2"),
				Hosts: []string{"lithium", "lithiumhost"},
			},
			{
				Addr:  netip.MustParseAddr("fe80::1%lo0"),
				Hosts: []string{"localhost"},
			},
		},
	}
	assertEqual(t, expected, actual)
}

func TestWrite(t *testing.T) {
	hs := hosts.Hosts{
		Entries: []hosts.HostEntry{
			{
				Addr:  netip.MustParseAddr("255.255.255.255"),
				Hosts: []string{"broadcasthost"},
			},
			{
				Addr:  netip.MustParseAddr("127.0.0.2"),
				Hosts: []string{"hydrogen"},
			},
			{
				Addr:  netip.MustParseAddr("127.0.0.3"),
				Hosts: []string{"hydrogen"},
			},
			{
				Addr:  netip.MustParseAddr("::2"),
				Hosts: []string{"hydrogen"},
			},
			{
				Addr:  netip.MustParseAddr("127.1.1.1"),
				Hosts: []string{"helium"},
			},
			{
				Addr:  netip.MustParseAddr("127.1.1.2"),
				Hosts: []string{"lithium", "lithiumhost"},
			},
			{
				Addr:  netip.MustParseAddr("fe80::1%lo0"),
				Hosts: []string{"localhost"},
			},
		},
	}

	w := new(strings.Builder)
	err := hs.Write(w)
	if err != nil {
		t.Fatal(err)
	}

	actual := w.String()
	expected := `255.255.255.255 broadcasthost
127.0.0.2 hydrogen
127.0.0.3 hydrogen
::2 hydrogen
127.1.1.1 helium
127.1.1.2 lithium lithiumhost
fe80::1%lo0 localhost
`

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestFromHost(t *testing.T) {
	dir, err := os.MkdirTemp("", "testdefaulthosts")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	f, err := os.Create(path.Join(dir, "hosts"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString("127.0.0.2 hydrogen")
	if err != nil {
		t.Fatal(err)
	}

	oldPath := hosts.DefaultPath
	defer func() { hosts.DefaultPath = oldPath }()
	hosts.DefaultPath = f.Name()

	actual, err := hosts.Local()
	if err != nil {
		t.Fatal(err)
	}

	expected := &hosts.Hosts{
		Entries: []hosts.HostEntry{
			{
				Addr:  netip.MustParseAddr("127.0.0.2"),
				Hosts: []string{"hydrogen"},
			},
		},
	}

	assertEqual(t, expected, actual)
}

func TestWriteEtcHosts(t *testing.T) {
	dir, err := os.MkdirTemp("", "testwriteetchost")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	hs, err := hosts.Local()
	if err != nil {
		t.Fatal(err)
	}

	path, err := hosts.WriteEtcHosts(dir, hs)
	if err != nil {
		t.Fatal(err)
	}

	hs2, err := hosts.ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, hs, hs2)
}

func TestDefault(t *testing.T) {
	hs := hosts.Default()
	w := new(strings.Builder)
	err := hs.Write(w)
	if err != nil {
		t.Fatal(err)
	}

	actual := w.String()
	expected := `127.0.0.1 localhost
::1 localhost localhost6 ip6-localhost ip6-loopback
`

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func assertEqual(t *testing.T, expected *hosts.Hosts, actual *hosts.Hosts) {
	if len(actual.Entries) != len(expected.Entries) {
		t.Fatalf("expected %d entries, got %d", len(expected.Entries), len(actual.Entries))
	}
	for i := range expected.Entries {
		if actual.Entries[i].Addr != expected.Entries[i].Addr {
			t.Errorf("expected addr %s, got %s", expected.Entries[i].Addr, actual.Entries[i].Addr)
		}
		if len(actual.Entries[i].Hosts) != len(expected.Entries[i].Hosts) {
			t.Errorf("expected %d hosts, got %d", len(expected.Entries[i].Hosts), len(actual.Entries[i].Hosts))
		}
		for j := range expected.Entries[i].Hosts {
			if actual.Entries[i].Hosts[j] != expected.Entries[i].Hosts[j] {
				t.Errorf("expected host %s, got %s", expected.Entries[i].Hosts[j], actual.Entries[i].Hosts[j])
			}
		}
	}
}
