module github.com/mokiat/lacking-js

go 1.17

require (
	github.com/mokiat/gomath v0.1.0
	github.com/mokiat/lacking v0.2.0
	github.com/mokiat/wasmgl v0.1.0
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
)

require golang.org/x/text v0.3.6 // indirect

replace github.com/mokiat/lacking => ../lacking
