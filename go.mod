module GoSG90

go 1.14

require (
	github.com/golang/mock v1.4.3 // indirect
	github.com/stianeikeland/go-rpio/v4 v4.4.0
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/tools v0.0.0-20200702044944-0cc1aa72b347 // indirect
)

replace github.com/stianeikeland/go-rpio/v4 v4.4.0 => github.com/ebarti/go-rpio/v4 v4.4.1
