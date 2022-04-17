module github.com/mokiat/lacking-js

go 1.18

require (
	github.com/mokiat/lacking v0.3.0
	github.com/mokiat/wasmgl v0.1.0
)

require (
	github.com/mokiat/gomath v0.2.0 // indirect
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/mokiat/gomath => ../gomath

replace github.com/mokiat/lacking => ../lacking

replace github.com/mokiat/wasmgl => ../wasmgl
