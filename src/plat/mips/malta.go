package mips

import (
//"fmt"
)

import (
	//"cpu/mips"
	//	"dev/mem"
	//	"dev/char"
	//	"dev/char/ser8250"
	"plat"
)

type PlatMalta struct {
	plat.Plat
}

func maltaInit() error {
	//mips.Mips
	return nil
}

func init() {
	println("Init plat/mips/malta")
	//plat.Register("malta", maltaInit())
}

func (pm *PlatMalta) Init() error {
	return nil
}

func (pm *PlatMalta) Start() error {
	return nil
}

func (pm *PlatMalta) Stop() error {
	return nil
}

func NewPlatMalta() *plat.Plat {
	p := plat.NewPlat()

	p.SetInfo("malta", "MIPS", "1.0")

	return p
}
