package hosts

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

// Hosts is a parsed /etc/hosts file.
type Hosts struct {
	Entries []HostEntry
}

// HostEntry associates IP addresses with host names.
type HostEntry struct {
	Addr  netip.Addr
	Hosts []string
}

// WriteEtcHosts writes the Hosts to a file in outputDir and returns the path to the file.
// The intent is that this is written to a containerd snapshot directory to be mounted into a container.
func WriteEtcHosts(outputDir string, hosts *Hosts) (string, error) {
	path := filepath.Join(outputDir, "hosts")

	f, err := os.CreateTemp(outputDir, "hosts")
	if err != nil {
		return "", err
	}
	defer func() { _ = os.Remove(f.Name()) }()

	err = hosts.Write(f)
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

// Local parses and returns the local /etc/hosts file.
func Local() (*Hosts, error) {
	return ParseFile(DefaultPath)
}

// Default returns a Hosts with the default localhost entries.
func Default() *Hosts {
	return &Hosts{
		Entries: []HostEntry{
			{
				Addr:  netip.MustParseAddr("127.0.0.1"),
				Hosts: []string{"localhost"},
			},
			{
				Addr: netip.MustParseAddr("::1"),
				// RedHat: localhost6 for ::1.
				// Debian: ip6-loopback & ip6-localhost for ::1.
				Hosts: []string{"localhost", "localhost6", "ip6-localhost", "ip6-loopback"},
			},
		},
	}
}

// ParseFile parses the named hosts file path.
func ParseFile(name string) (*Hosts, error) {
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

// Parse parses a hosts file from r.
func Parse(r io.Reader) (*Hosts, error) {
	etcHosts := new(Hosts)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line, _, _ = strings.Cut(line, "#") // remove any comments
		line = strings.TrimSpace(line)
		i := strings.IndexAny(line, " \t")
		if i == -1 {
			continue
		}

		addr, err := netip.ParseAddr(line[:i])
		if err != nil {
			continue
		}

		hosts := strings.TrimSpace(line[i+1:])
		hostEntry := HostEntry{
			Addr: addr,
		}

		for len(hosts) > 0 {
			host := hosts
			i := strings.IndexAny(hosts, " \t")
			if i != -1 {
				host = host[:i]
				hosts = strings.TrimSpace(hosts[i+1:])
			} else {
				hosts = ""
			}

			hostEntry.Hosts = append(hostEntry.Hosts, host)
		}

		if len(hostEntry.Hosts) > 0 {
			etcHosts.Entries = append(etcHosts.Entries, hostEntry)
		}
	}

	return etcHosts, scanner.Err()
}

// Write writes hosts  to w.
func (h *Hosts) Write(w io.Writer) error {
	buf := new(bytes.Buffer)
	for _, entry := range h.Entries {
		io.WriteString(buf, entry.Addr.String())
		io.WriteString(buf, " ")
		io.WriteString(buf, strings.Join(entry.Hosts, " "))
		io.WriteString(buf, "\n")
	}

	_, err := w.Write(buf.Bytes())
	return err
}
