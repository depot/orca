package resolvconf_test

import (
	"net/netip"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/depot/orca/resolvconf"
)

func TestParseFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "testresolvconf")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	conf := resolvconf.DefaultResolvConf()
	path, err := resolvconf.WriteResolvConf(dir, conf)
	if err != nil {
		t.Fatal(err)
	}

	conf2, err := resolvconf.ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, conf, conf2)
}

func TestParse(t *testing.T) {
	systemdResolvConf := `nameserver 127.0.0.53
	options edns0 trust-ad
	search .`

	actual, err := resolvconf.Parse(strings.NewReader(systemdResolvConf))
	if err != nil {
		t.Fatal(err)
	}
	expected := &resolvconf.Config{
		Nameservers: []netip.Addr{
			netip.AddrFrom4([4]byte{127, 0, 0, 53}),
		},
		Options: []string{"edns0", "trust-ad"},
		SearchDomains: []string{
			".",
		},
	}

	assertEqual(t, expected, actual)

	tailscaleResolvConf := `search taila12bc.ts.net taild34ef.ts.net
	nameserver 100.100.100.100`

	actual, err = resolvconf.Parse(strings.NewReader(tailscaleResolvConf))
	if err != nil {
		t.Fatal(err)
	}

	expected = &resolvconf.Config{
		Nameservers: []netip.Addr{
			netip.AddrFrom4([4]byte{100, 100, 100, 100}),
		},
		SearchDomains: []string{
			"taila12bc.ts.net",
			"taild34ef.ts.net",
		},
	}

	assertEqual(t, expected, actual)
}

func TestFromHost(t *testing.T) {
	dir, err := os.MkdirTemp("", "testresolvconf")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	systemdPath := path.Join(dir, "systemd-resolv.conf")
	f, err := os.Create(systemdPath)
	if err != nil {
		t.Fatal(err)
	}

	systemdResolvConf := &resolvconf.Config{
		Nameservers: []netip.Addr{
			netip.AddrFrom4([4]byte{155, 98, 64, 64}),
			netip.AddrFrom4([4]byte{155, 98, 111, 100}),
		},
		SearchDomains: []string{
			"cs.utah.edu",
			"eng.utah.edu",
		},
		Options: []string{
			"edns0",
			"trust-ad",
		},
	}

	err = systemdResolvConf.Write(f)
	_ = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	oldSystemdPath := resolvconf.SystemdPath
	defer func() { resolvconf.SystemdPath = oldSystemdPath }()
	resolvconf.SystemdPath = f.Name()

	confPath := path.Join(dir, "resolv.conf")
	f, err = os.Create(confPath)
	if err != nil {
		t.Fatal(err)
	}

	conf := &resolvconf.Config{
		Nameservers: []netip.Addr{
			netip.AddrFrom4([4]byte{127, 0, 0, 53}),
		},
	}

	err = conf.Write(f)
	_ = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	oldDefaultPath := resolvconf.DefaultPath
	defer func() { resolvconf.DefaultPath = oldDefaultPath }()
	resolvconf.DefaultPath = f.Name()

	actual, err := resolvconf.FromHost()
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, systemdResolvConf, actual)
}

func assertEqual(t *testing.T, conf *resolvconf.Config, conf2 *resolvconf.Config) {
	if len(conf2.Nameservers) != len(conf.Nameservers) {
		t.Fatalf("expected %d nameservers, got %d", len(conf.Nameservers), len(conf2.Nameservers))
	}

	for i, ns := range conf.Nameservers {
		if ns.String() != conf2.Nameservers[i].String() {
			t.Fatalf("expected nameserver %d to be %s, got %s", i, ns.String(), conf2.Nameservers[i].String())
		}
	}

	if len(conf2.SearchDomains) != len(conf.SearchDomains) {
		t.Fatalf("expected %d search domains, got %d", len(conf.SearchDomains), len(conf2.SearchDomains))
	}

	for i, sd := range conf.SearchDomains {
		if sd != conf2.SearchDomains[i] {
			t.Fatalf("expected search domain %d to be %s, got %s", i, sd, conf2.SearchDomains[i])
		}
	}

	if len(conf2.Options) != len(conf.Options) {
		t.Fatalf("expected %d options, got %d", len(conf.Options), len(conf2.Options))
	}

	for i, o := range conf.Options {
		if o != conf2.Options[i] {
			t.Fatalf("expected option %d to be %s, got %s", i, o, conf2.Options[i])
		}
	}
}
