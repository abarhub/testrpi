module testrpi

go 1.20

//replace machine => C:/tinygo/src/machine

require (
	github.com/shanghuiyang/rpi-devices v0.0.0-20220821024425-835ee611bf61
	periph.io/x/conn/v3 v3.7.0
	periph.io/x/devices/v3 v3.7.1
	periph.io/x/host/v3 v3.8.2
	tinygo.org/x/drivers v0.26.0
)

require (
	github.com/mdp/monochromeoled v0.0.0-20171027213216-a0c6b5c996cf // indirect
	github.com/stianeikeland/go-rpio/v4 v4.6.0 // indirect
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07 // indirect
	golang.org/x/exp v0.0.0-20210526181343-b47a03e3048a // indirect
	golang.org/x/sys v0.1.0 // indirect
)
