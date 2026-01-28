//go:build darwin

package integration

import tea "github.com/charmbracelet/bubbletea"

type Instance struct{}

func Init(p *tea.Program) *Instance {

	return nil
}

func (ins *Instance) Close() {}

func (ins *Instance) UpdateStatus(status string)   {}
func (ins *Instance) UpdateMetadata(meta Metadata) {}
func (ins *Instance) ClearMetadata()               {}

func (ins *Instance) GetStatus() string {
	return "Stopped"
}
