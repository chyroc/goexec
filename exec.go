package goexec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Exec struct {
	Name       string          // command name
	Args       []string        // command args
	dir        string          // run dir
	ignoreEnvs map[string]bool // ignore env to pass to command
	setEnvs    map[string]string
	isAssert   bool
	isLog      bool
}

// New create new Exec instance
func New(args ...string) *Exec {
	cmd := new(Exec)
	for idx, v := range args {
		if idx == 0 {
			cmd.Name = v
		} else {
			cmd.Args = append(cmd.Args, v)
		}
	}

	cmd.ignoreEnvs = make(map[string]bool)
	cmd.setEnvs = make(map[string]string)

	return cmd
}

// SetDir set command run dir
func (r *Exec) SetDir(s string) *Exec {
	r.dir = s
	return r
}

// IgnoreEnv ignore envs
func (r *Exec) IgnoreEnv(envs ...string) *Exec {
	for _, v := range envs {
		r.ignoreEnvs[v] = true
	}
	return r
}

// WithAssert must run success
func (r *Exec) WithAssert() *Exec {
	r.isAssert = true
	return r
}

// WithLog log command
func (r *Exec) WithLog() *Exec {
	r.isLog = true
	return r
}

// SetEnv set env
func (r *Exec) SetEnv(k, v string) *Exec {
	r.setEnvs[k] = v
	return r
}

func (r *Exec) getEnvs() []string {
	result := []string{}
	for k, v := range getEnvs() {
		if !r.ignoreEnvs[k] {
			result = append(result, fmt.Sprintf("%s=%s", k, v))
		}
	}
	for k, v := range r.setEnvs {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}

	return result
}

// Run run command
func (r *Exec) Run() (stdout string, stderr string, err error) {
	return r.run(false)
}

// RunInStream print env to os.stdout and os.stderr
func (r *Exec) RunInStream() error {
	_, _, err := r.run(true)
	return err
}

func (r *Exec) run(isStream bool) (sout string, serr string, err error) {
	if r.isLog {
		_, _ = fmt.Fprintf(os.Stdout, r.formatCommand())
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.Command(r.Name, r.Args...)
	if isStream {
		cmd.Stdout = newDoubleWriter(os.Stdout, stdout)
		cmd.Stderr = newDoubleWriter(os.Stderr, stderr)
	} else {
		cmd.Stdout = stdout
		cmd.Stderr = stderr
	}
	cmd.Dir = r.dir
	cmd.Env = r.getEnvs()
	err = cmd.Run()

	if err != nil && r.isAssert {
		os.Exit(1)
	}

	return stdout.String(), stderr.String(), err
}

func (r *Exec) formatCommand() string {
	s := strings.Builder{}
	s.WriteString(r.Name)
	for _, v := range r.Args {
		s.WriteString(" " + v)
	}
	s.WriteString(fmt.Sprintf(", dir=%q", r.dir))

	return s.String()
}

func getEnvs() map[string]string {
	envs := make(map[string]string)
	for _, v := range os.Environ() {
		pair := strings.SplitN(v, "=", 2)
		switch len(pair) {
		case 0, 1:
			continue
		default:
			key, value := pair[0], pair[1]
			envs[key] = value
		}
	}

	return envs
}

type doubleWriter struct {
	main   io.Writer
	second io.Writer
}

func newDoubleWriter(main, second io.Writer) io.Writer {
	return &doubleWriter{main: main, second: second}
}

func (r *doubleWriter) Write(p []byte) (n int, err error) {
	_, _ = r.second.Write(p)
	return r.main.Write(p)
}
