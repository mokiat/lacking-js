module github.com/mokiat/lacking-js

go 1.17

require (
	github.com/mokiat/gomath v0.1.0
	github.com/mokiat/lacking v0.2.0
	github.com/mokiat/wasmgl v0.0.0-20181122210432-a78877b4f1c3
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
)

require golang.org/x/text v0.3.6 // indirect

replace github.com/mokiat/lacking => ../lacking

replace github.com/mokiat/wasmgl => ../../Libraries/wasmgl
