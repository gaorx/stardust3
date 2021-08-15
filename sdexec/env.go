package sdexec

import (
	"github.com/gaorx/stardust3/sderr"
)

type Env struct {
	Dir  string
	Vars map[string]string
}

var (
	NoEnv = Env{}
)

func (p *Env) ensure() *Env {
	if p.Vars == nil {
		p.Vars = map[string]string{}
	}
	return p
}

func (p *Env) applyCmd(cmd *Cmd) *Cmd {
	return cmd.SetDir(p.Dir).SetVars(p.Vars)
}

func (p Env) Run(line string) error {
	cmd, err := Parse(line)
	if err != nil {
		return sderr.WithStack(err)
	}
	return sderr.WithStack(p.applyCmd(cmd).Run())
}

func (p Env) RunOutput(line string, combine bool) ([]byte, error) {
	cmd, err := Parse(line)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	r, err := p.applyCmd(cmd).RunOutput(combine)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func (p Env) RunOutputString(line string, combine bool) (string, error) {
	cmd, err := Parse(line)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	r, err := p.applyCmd(cmd).RunOutputString(combine)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	return r, nil
}

func (p Env) RunConsole(line string) error {
	cmd, err := Parse(line)
	if err != nil {
		return sderr.WithStack(err)
	}
	err = p.applyCmd(cmd).RunConsole()
	return sderr.WithStack(err)
}

func (p Env) RunResult(line string) (*Result, error) {
	cmd, err := Parse(line)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return p.applyCmd(cmd).RunResult(), nil
}
