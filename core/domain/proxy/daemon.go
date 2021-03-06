package proxy

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"regexp"
)

const ServiceFile = "/etc/systemd/system/vite.service"

var EnableCmds = [][]string{
	{"systemctl", "daemon-reload"},
	{"systemctl", "enable", "vite.service"},
	{"systemctl", "start", "vite.service"},
}

var DisableCmds = [][]string{
	{"systemctl", "stop", "vite.service"},
	{"systemctl", "disable", "vite.service"},
	{"systemctl", "daemon-reload"},
}

type DaemonStatus int

const (
	Absent DaemonStatus = iota
	Running
	Errored
)

//go:embed vite.service
var daemonStatus string

func (s DaemonStatus) String() string {
	switch s {
	case Running:
		return "running"
	case Errored:
		return "errored"
	case Absent:
		return "absent"
	default:
		panic("invalid service status")
	}
}

const UptimeRegex = `(?m); (.+) ago\n`

func State() (DaemonStatus, string, error) {
	out, err := exec.Command("systemctl", "status", "vite.service").Output()

	re := regexp.MustCompile(UptimeRegex)
	matches := re.FindSubmatch(out)

	uptime := ""
	if len(matches) > 1 {
		uptime = string(matches[1])
	} else {
		uptime = "errored"
	}

	if err != nil {
		code := err.(*exec.ExitError).ExitCode()

		if code == 3 || code == 4 {
			return Absent, uptime, nil
		}

		return Absent, uptime, fmt.Errorf("%w: %s", err, out)
	}

	if bytes.Contains(out, []byte("code=exited")) {
		return Errored, uptime, nil
	}

	return Running, uptime, nil
}

func DaemonConfig(user string) ([]byte, error) {
	e, err := os.Executable()
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("vite.service").Parse(daemonStatus)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, struct {
		User string
		Cmd  string
	}{
		User: user,
		Cmd:  e + " proxy run",
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
