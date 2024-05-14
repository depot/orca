package buildspec_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/platforms"
	buildspec "github.com/depot/orca/buildspec"
)

func ExampleWithBuildSpec() {
	opts := &buildspec.Options{
		Platform: &platforms.Platform{
			OS: "linux",
		},
		RootFSPath:     "/tmp",
		Args:           []string{"sh", "-c", "echo hello world"},
		Userstr:        "1",
		ReadonlyRootFS: false,
		RootlessRunC:   true,
	}
	buildSpecOpts := buildspec.WithBuildSpec(opts)

	ctx := namespaces.WithNamespace(context.Background(), "test1")
	spec, err := oci.GenerateSpecWithPlatform(
		ctx,
		nil,
		"linux/amd64",
		&containers.Container{},
		buildSpecOpts,
	)
	if err != nil {
		panic(err)
	}

	buf, _ := json.MarshalIndent(spec, "", "  ")
	output := string(buf)
	fmt.Printf("%s\n", output)
	//Output:
	//{
	//   "ociVersion": "1.1.0",
	//   "process": {
	//     "user": {
	//       "uid": 0,
	//       "gid": 0,
	//       "additionalGids": [
	//         0
	//       ],
	//       "username": "1"
	//     },
	//     "args": [
	//       "sh",
	//       "-c",
	//       "echo hello world"
	//     ],
	//     "env": [
	//       "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
	//     ],
	//     "cwd": "",
	//     "capabilities": {
	//       "bounding": [
	//         "CAP_CHOWN",
	//         "CAP_DAC_OVERRIDE",
	//         "CAP_FSETID",
	//         "CAP_FOWNER",
	//         "CAP_MKNOD",
	//         "CAP_NET_RAW",
	//         "CAP_SETGID",
	//         "CAP_SETUID",
	//         "CAP_SETFCAP",
	//         "CAP_SETPCAP",
	//         "CAP_NET_BIND_SERVICE",
	//         "CAP_SYS_CHROOT",
	//         "CAP_KILL",
	//         "CAP_AUDIT_WRITE"
	//       ],
	//       "effective": [
	//         "CAP_CHOWN",
	//         "CAP_DAC_OVERRIDE",
	//         "CAP_FSETID",
	//         "CAP_FOWNER",
	//         "CAP_MKNOD",
	//         "CAP_NET_RAW",
	//         "CAP_SETGID",
	//         "CAP_SETUID",
	//         "CAP_SETFCAP",
	//         "CAP_SETPCAP",
	//         "CAP_NET_BIND_SERVICE",
	//         "CAP_SYS_CHROOT",
	//         "CAP_KILL",
	//         "CAP_AUDIT_WRITE"
	//       ],
	//       "permitted": [
	//         "CAP_CHOWN",
	//         "CAP_DAC_OVERRIDE",
	//         "CAP_FSETID",
	//         "CAP_FOWNER",
	//         "CAP_MKNOD",
	//         "CAP_NET_RAW",
	//         "CAP_SETGID",
	//         "CAP_SETUID",
	//         "CAP_SETFCAP",
	//         "CAP_SETPCAP",
	//         "CAP_NET_BIND_SERVICE",
	//         "CAP_SYS_CHROOT",
	//         "CAP_KILL",
	//         "CAP_AUDIT_WRITE"
	//       ]
	//     },
	//     "noNewPrivileges": true
	//   },
	//   "root": {
	//     "path": "/tmp"
	//   },
	//   "mounts": [
	//     {
	//       "destination": "/proc",
	//       "type": "proc",
	//       "source": "proc",
	//       "options": [
	//         "nosuid",
	//         "noexec",
	//         "nodev"
	//       ]
	//     },
	//     {
	//       "destination": "/dev",
	//       "type": "tmpfs",
	//       "source": "tmpfs",
	//       "options": [
	//         "nosuid",
	//         "strictatime",
	//         "mode=755",
	//         "size=65536k"
	//       ]
	//     },
	//     {
	//       "destination": "/dev/pts",
	//       "type": "devpts",
	//       "source": "devpts",
	//       "options": [
	//         "nosuid",
	//         "noexec",
	//         "newinstance",
	//         "ptmxmode=0666",
	//         "mode=0620",
	//         "gid=5"
	//       ]
	//     },
	//     {
	//       "destination": "/dev/shm",
	//       "type": "tmpfs",
	//       "source": "shm",
	//       "options": [
	//         "nosuid",
	//         "noexec",
	//         "nodev",
	//         "mode=1777",
	//         "size=65536k"
	//       ]
	//     },
	//     {
	//       "destination": "/dev/mqueue",
	//       "type": "mqueue",
	//       "source": "mqueue",
	//       "options": [
	//         "nosuid",
	//         "noexec",
	//         "nodev"
	//       ]
	//     },
	//     {
	//       "destination": "/run",
	//       "type": "tmpfs",
	//       "source": "tmpfs",
	//       "options": [
	//         "nosuid",
	//         "strictatime",
	//         "mode=755",
	//         "size=65536k"
	//       ]
	//     },
	//     {
	//       "destination": "/etc/hosts",
	//       "type": "bind",
	//       "source": "/etc/hosts",
	//       "options": [
	//         "rbind",
	//         "ro"
	//       ]
	//     }
	//   ],
	//   "linux": {
	//     "namespaces": [
	//       {
	//         "type": "pid"
	//       },
	//       {
	//         "type": "ipc"
	//       },
	//       {
	//         "type": "uts"
	//       },
	//       {
	//         "type": "mount"
	//       }
	//     ],
	//     "seccomp": {
	//       "defaultAction": ""
	//     },
	//     "maskedPaths": [
	//       "/proc/acpi",
	//       "/proc/asound",
	//       "/proc/kcore",
	//       "/proc/keys",
	//       "/proc/latency_stats",
	//       "/proc/timer_list",
	//       "/proc/timer_stats",
	//       "/proc/sched_debug",
	//       "/sys/firmware",
	//       "/sys/devices/virtual/powercap",
	//       "/proc/scsi"
	//     ],
	//     "readonlyPaths": [
	//       "/proc/bus",
	//       "/proc/fs",
	//       "/proc/irq",
	//       "/proc/sys",
	//       "/proc/sysrq-trigger"
	//     ]
	//   }
	//}
}
