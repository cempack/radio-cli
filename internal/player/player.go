package player

import (
	"os/exec"
	"sync"
)

type PlayerState int

const (
	StateStopped PlayerState = iota
	StateConnecting
	StatePlaying
	StateError
)

type Player struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	state  PlayerState
	url    string
	errMsg string
}

func New() *Player {
	return &Player{state: StateStopped}
}

func Available() bool {
	_, err := exec.LookPath("mpv")
	return err == nil
}

func (p *Player) Play(streamURL string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
		_ = p.cmd.Wait()
	}

	if !Available() {
		p.state = StateError
		p.errMsg = "mpv not found in PATH"
		return nil
	}

	p.url = streamURL
	p.state = StateConnecting
	p.errMsg = ""
	p.cmd = exec.Command("mpv", "--no-video", "--quiet", streamURL)

	if err := p.cmd.Start(); err != nil {
		p.state = StateError
		p.errMsg = err.Error()
		return err
	}

	p.state = StatePlaying

	go func() {
		err := p.cmd.Wait()
		p.mu.Lock()
		defer p.mu.Unlock()
		if p.state != StateStopped {
			if err != nil {
				p.state = StateError
				p.errMsg = err.Error()
			} else {
				p.state = StateStopped
			}
		}
	}()

	return nil
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
		_ = p.cmd.Wait()
		p.cmd = nil
	}
	p.state = StateStopped
	p.url = ""
}

func (p *Player) State() PlayerState {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.state
}

func (p *Player) URL() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.url
}

func (p *Player) Error() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.errMsg
}
