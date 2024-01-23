// Package resolvconf provides a parser for resolv.conf files.
//
// Inspired from nerdctl's resolvconf and tailscale's resolvconffile.
package resolvconf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
)

var (
	// DefaultPath is the default path to the resolv.conf that contains information to resolve DNS.
	DefaultPath = "/etc/resolv.conf"
	// SystemdPath is systemd's default path to the resolv.conf that contains information to resolve DNS.
	// This path is used when we detect systemd's systemd-resolved ip address of 127.0.0.53 in /etc/resolv.conf
	SystemdPath = "/run/systemd/resolve/resolv.conf"
)

// Note: the default IPv4 & IPv6 resolvers are set to Google's Public DNS
// This follows Docker's and nerdctl's default behavior.
//
// These are only used if the host's resolv.conf is not available or the name server is localhost.
var DefaultNameservers = []string{
	"8.8.8.8",
	"8.8.4.4",
	"2001:4860:4860::8888",
	"2001:4860:4860::8844",
}

// Config is a parsed resolv.conf file.
type Config struct {
	Nameservers   []netip.Addr
	SearchDomains []string
	Options       []string
}

// DefaultResolvConf returns a default resolv.conf using `DefaultNameservers`.
func DefaultResolvConf() *Config {
	return &Config{
		Nameservers: defaultNameserverAddrs(),
	}
}

// WriteResolvConf writes the Config to a file in outputDir and returns the path to the file.
// The intent is that this is written to a containerd snapshot directory to be mounted into a container.
func WriteResolvConf(outputDir string, conf *Config) (string, error) {
	path := filepath.Join(outputDir, "resolv.conf")

	f, err := os.CreateTemp(outputDir, "resolv.conf")
	if err != nil {
		return "", err
	}
	defer func() { _ = os.Remove(f.Name()) }()

	err = conf.Write(f)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	if err != nil {
		return "", err
	}

	err = os.Chmod(f.Name(), 0644)
	if err != nil {
		return "", err
	}

	err = os.Rename(f.Name(), path)
	if err != nil {
		return "", err
	}

	return path, nil
}

// FromHost parses and returns the host's resolv.conf.
// If it detects systemd's systemd-resolved ip address
// it will use systemd's resolv.conf instead.
//
// Additionally, it filters out any loopback addresses as it is assumed
// that this Config will be used in a container and containers may not be
// able to use loopback addresses.
func FromHost() (*Config, error) {
	conf, err := ParseFile(DefaultPath)
	if err != nil {
		return nil, err
	}

	systemdNameserver := netip.AddrFrom4([4]byte{127, 0, 0, 53})

	if len(conf.Nameservers) == 1 && conf.Nameservers[0] == systemdNameserver {
		conf, err = ParseFile(SystemdPath)
		if err != nil {
			if os.IsNotExist(err) {
				return DefaultResolvConf(), nil
			}
			return nil, err
		}
	}

	conf.FilterLoopback()

	return conf, nil
}

// ParseFile parses the named resolv.conf file.
func ParseFile(name string) (*Config, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if n := fi.Size(); n > 10<<10 {
		return nil, fmt.Errorf("unexpectedly large %q file: %d bytes", name, n)
	}
	all, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return Parse(bytes.NewReader(all))
}

// Parse parses a resolv.conf file from r.
func Parse(r io.Reader) (*Config, error) {
	config := new(Config)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line, _, _ = strings.Cut(line, "#") // remove any comments
		line = strings.TrimSpace(line)

		if s, ok := strings.CutPrefix(line, "nameserver"); ok {
			nameserver := strings.TrimSpace(s)
			if len(nameserver) == len(s) {
				return nil, fmt.Errorf("missing space after \"nameserver\" in %q", line)
			}
			ip, err := netip.ParseAddr(nameserver)
			if err != nil {
				return nil, err
			}
			config.Nameservers = append(config.Nameservers, ip)
			continue
		}

		if s, ok := strings.CutPrefix(line, "search"); ok {
			domains := strings.TrimSpace(s)
			if len(domains) == len(s) {
				return nil, fmt.Errorf("missing space after search in %q", line)
			}
			for len(domains) > 0 {
				domain := domains
				i := strings.IndexAny(domain, " \t")
				if i != -1 {
					domain = domain[:i]
					domains = strings.TrimSpace(domains[i+1:])
				} else {
					domains = ""
				}
				config.SearchDomains = append(config.SearchDomains, domain)
			}
		}

		if o, ok := strings.CutPrefix(line, "options"); ok {
			options := strings.TrimSpace(o)
			if len(options) == len(o) {
				return nil, fmt.Errorf("missing space after options in %q", line)
			}
			for len(options) > 0 {
				option := options
				i := strings.IndexAny(options, " \t")
				if i != -1 {
					option = option[:i]
					options = strings.TrimSpace(options[i+1:])
				} else {
					options = ""
				}
				config.Options = append(config.Options, option)
			}
		}
	}
	return config, nil
}

// FilterLoopback filters out any loopback addresses because
// it is assumed that this Config will be used in a container
// and containers may not be able to use loopback addresses.
func (rc *Config) FilterLoopback() {
	nameservers := []netip.Addr{}
	for _, ns := range rc.Nameservers {
		if ns.IsLoopback() {
			continue
		}
		nameservers = append(nameservers, ns)
	}

	if len(nameservers) == 0 {
		nameservers = defaultNameserverAddrs()
	}

	rc.Nameservers = nameservers
}

// Write writes the Config to w.
func (rc *Config) Write(w io.Writer) error {
	buf := new(bytes.Buffer)

	for _, ns := range rc.Nameservers {
		io.WriteString(buf, "nameserver ")
		io.WriteString(buf, ns.String())
		io.WriteString(buf, "\n")
	}

	if len(rc.SearchDomains) > 0 {
		io.WriteString(buf, "search")
		for _, domain := range rc.SearchDomains {
			if strings.Trim(domain, " ") == "." {
				continue
			}
			io.WriteString(buf, " ")
			io.WriteString(buf, domain)
		}
		io.WriteString(buf, "\n")
	}

	if len(rc.Options) > 0 {
		io.WriteString(buf, "options")
		for _, option := range rc.Options {
			if strings.Trim(option, " ") == "" {
				continue
			}
			io.WriteString(buf, " ")
			io.WriteString(buf, option)
		}
	}

	_, err := w.Write(buf.Bytes())
	return err
}

func defaultNameserverAddrs() []netip.Addr {
	nameservers := []netip.Addr{}
	for _, ip := range DefaultNameservers {
		addr, _ := netip.ParseAddr(ip)
		nameservers = append(nameservers, addr)
	}
	return nameservers
}
