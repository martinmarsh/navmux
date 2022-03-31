/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package pid

import (
	"math"
)

type Pid struct {
    integral float64
    last float64
	max_output float64
	kp float64
	ki float64
	kd float64
	gain float64
	Scale_kp float64
	Scale_ki float64
	Scale_kd float64
	Scale_gain float64
}

// the computation factors are predefind after running some PID tuning and use to create a PID object
// During running the parameters can be trimmed or adjusted using the scaling factors
// which are by default set to 1
// gain is used be used to adjust the overall output to within the range used by the actuator.
// Scale_gain can be used to trim the runtime sensitivity
func MakePid(kp, ki, kd, gain, max_output float64) *Pid {
	p := Pid{
		integral: 0,
		last: 0,
		max_output: max_output,
		kp: kp,
		ki: ki,
		kd: kd,
		gain: gain,
		Scale_kp: 1.0,
		Scale_ki: 1.0,
		Scale_kd: 1.0,
		Scale_gain: 1.0,
	}
	return &p
} 

func (p *Pid) Reset() {
	p.integral = 0
}

// Compute the Actuating Signal from the error term
// The error term is supplied externally and is typically the command signal (or Set Point )
// minus the Process Variable (feed back value).  
// This is done by the application since Sp-Pv but in some cases may need conditioning for example
// compass substrations should be based +/- 180 after substraction. 
// The result is the actuator value.
// To assist with scaling the paramenters scaling variables may be set. This makes it easier to use
// user friendly values for settings
// -pv is used instead of sp-pv as avoids spiking if set point changed
// The assumption is a constant calculation rate
func (p *Pid) Compute(sp_pv, pv float64) float64 {

	proportional := sp_pv * p.kp * p.Scale_kp
	i_in :=  sp_pv * p.ki * p.Scale_ki
	d_inc := -pv * p.kd * p.Scale_kd

	p.integral += i_in
	differential := d_inc - p.last
	p.last = d_inc

	as := (p.integral + proportional + differential) * p.gain * p.Scale_gain

	// Integral latch up protection
	// Prevent further integration if max output is achieved and addition to integral is same sign
	if math.Abs(as) > p.max_output{
		if same_sign(i_in, as){
			p.integral -= i_in
		}
	}

	return as	
}

func same_sign(x, y float64) bool {
	if x >= 0.0 && y >= 0.0 {
		return true
	} else if x < 0.0 && y <0.0 {
		return true
	}
	return false
}
