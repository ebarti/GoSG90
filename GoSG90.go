package GoSG90

import (
	"errors"
	"github.com/stianeikeland/go-rpio/v4"
)

/*
type sg90servo interface {
	new(GPIOindex uint8) SG90
	init() (err error)
	getCurrentLocation() (currentAngle float64, err error)
	setTargetLocation(targetAngle float64) (err error)
	movePlus() (err error)
	moveMinus() (err error)
	angleRange() float64
	dutyRange() float64
	toAngle(dutyLength, cycleLength uint32) (angle float64)
	toDuty(targetAngle float64) (dutyLength, cycleLength uint32, err error)
}
*/

type mockablePin struct {
}

type SG90 struct {
	Pin                rpio.Pinner
	OperatingFreq      uint32
	MaxDutyCycleFactor float64
	MinDutyCycleFactor float64
	MinAngle           float64
	MaxAngle           float64
}

// Pin.Freq -> Param freq should be in range 4688Hz - 19.2MHz to prevent unexpected behavior
// Output frequency is computed as pwm clock frequency divided by cycle length.
// For a 50Hz pwm:
//	Min PWM clock: 4688Hz
//	OFreq: 50H
// 	CycleLength: 4688/50=94
// So, to set Pwm pin to freqency 50Hz with duty cycle 1/4, use this combination:
//
//  pin.Pwm()
//  pin.DutyCycle(1, 4)
//  pin.Freq(38000*4)

// Creates a new SG90 servo using the desired GPIO pin number.
// This module uses PWM. A raspberry pi has 2 PWM channels:
//   channel 1 (pwm0) for pins 12, 18, 40
//   channel 2 (pwm1) for pins 13, 19, 41, 45.
func New(GPIOindex uint8) SG90 {
	// raspberry pi 4
	validPins := map[uint8]bool{
		12: true,
		13: true,
		18: true,
		19: true,
		40: true,
		41: true,
		45: true,
	}
	if _, exists := validPins[GPIOindex]; false == exists {
		panic(errors.New("invalid GPIO number"))
	}
	return SG90{
		Pin:           rpio.Pin(GPIOindex),
		OperatingFreq: 50,
		MinAngle:      -90.0,
		MaxAngle:      90.0,
	}
}

func (servo *SG90) Init() (err error) {
	err = rpio.Open()
	return
}

func (servo *SG90) GetCurrentLocation() (currentAngle float64, err error) {
	currentAngle = 0.0
	dutyLength, cycleLength, err := servo.Pin.GetDutyCycle()
	if nil == err {
		currentAngle = servo.toAngle(dutyLength, cycleLength)
	}
	return
}

func (servo *SG90) SetTargetLocation(targetAngle float64) (err error) {
	if targetAngle < servo.MinAngle || targetAngle > servo.MaxAngle {
		err = errors.New("angle out of bounds")
		return
	}
	targetDutyLength, targetCycleLength, err := servo.toDuty(targetAngle)
	if nil == err {
		servo.Pin.SetDutyCycle(targetDutyLength, targetCycleLength)
	}
	return
}

func (servo *SG90) MovePlus() (err error) {
	currentAngle, err := servo.GetCurrentLocation()
	if nil == err {
		err = servo.SetTargetLocation(currentAngle + 1)
	}
	return
}

func (servo *SG90) MoveMinus() (err error) {
	currentAngle, err := servo.GetCurrentLocation()
	if nil == err {
		err = servo.SetTargetLocation(currentAngle - 1)
	}
	return
}

func (servo *SG90) angleRange() float64 {
	return servo.MaxAngle - servo.MinAngle
}

func (servo *SG90) dutyRange() float64 {
	return servo.MaxDutyCycleFactor - servo.MinDutyCycleFactor
}

func (servo *SG90) toAngle(dutyLength, cycleLength uint32) (angle float64) {
	currentDutyCycle := float64(dutyLength) / float64(cycleLength)
	angle = servo.MinAngle + (currentDutyCycle-servo.MinDutyCycleFactor)*servo.angleRange()/servo.dutyRange()
	return
}

func (servo *SG90) toDuty(targetAngle float64) (dutyLength, cycleLength uint32, err error) {
	dutyLength, cycleLength, err = servo.Pin.GetDutyCycle()
	if nil == err {
		targetDutyCycle := servo.MinDutyCycleFactor + (targetAngle-servo.MinAngle)*servo.dutyRange()/servo.angleRange()
		dutyLength = uint32(targetDutyCycle * float64(cycleLength))
	}
	return
}
