package gamebro

type JoyPad struct {
	selectReg   byte
	buttonState byte // low-active
	dirState    byte // low-active
}

func NewJoyPad() *JoyPad {
	return &JoyPad{0xF0, 0x0F, 0x0F}
}

func (j *JoyPad) SetStates(buttonState, dirState byte) {
	j.buttonState = dirState
	j.dirState = dirState
}

func (j *JoyPad) Read() byte {
	var inputs byte = 0x0F
	if j.selectReg&0x20 == 0 {
		inputs &= j.buttonState
	}
	if j.selectReg&0x10 == 0 {
		inputs &= j.dirState
	}
	return j.selectReg | inputs
}

func (j *JoyPad) Write(data byte) {
	j.selectReg = 0xC0 | data&0x30
}
