package internal

func NewCommandBuffer() *CommandBuffer {
	return &CommandBuffer{}
}

type CommandBuffer struct{}

// Following is an idea for a command queue that can hold various types.

// package main

// import (
// 	"fmt"
// 	"unsafe"
// )

// type A struct {
// 	First  uint8
// 	Second uint8
// }

// type B struct {
// 	Value [4]float32
// }

// type Header struct {
// 	Kind uint8
// }

// func write[T any](buffer []byte, v T, offset uintptr) uintptr {
// 	target := (*T)(unsafe.Pointer(&buffer[offset]))
// 	*target = v
// 	return offset + unsafe.Sizeof(v)
// }

// func read[T any](buffer []byte, offset uintptr) (T, uintptr) {
// 	target := (*T)(unsafe.Pointer(&buffer[offset]))
// 	return *target, offset + unsafe.Sizeof(*target)
// }

// func main() {
// 	buffer := make([]byte, 64)
// 	offset := write(buffer, Header{
// 		Kind: 1,
// 	}, 0)
// 	offset = write(buffer, A{
// 		First:  7,
// 		Second: 13,
// 	}, offset)
// 	offset = write(buffer, Header{
// 		Kind: 2,
// 	}, offset)
// 	offset = write(buffer, B{
// 		Value: [4]float32{1.0, 2.0, 3.0, 4.0},
// 	}, offset)
// 	offset = write(buffer, Header{
// 		Kind: 2,
// 	}, offset)
// 	offset = write(buffer, B{
// 		Value: [4]float32{8.0, 7.0, 6.0, 5.0},
// 	}, offset)
// 	offset = write(buffer, Header{
// 		Kind: 0,
// 	}, offset)

// 	offset = 0
// 	for {
// 		var header Header
// 		header, offset = read[Header](buffer, offset)
// 		switch header.Kind {
// 		case 1:
// 			var a A
// 			a, offset = read[A](buffer, offset)
// 			fmt.Printf("A: %+v\n", a)
// 		case 2:
// 			var b B
// 			b, offset = read[B](buffer, offset)
// 			fmt.Printf("B: %+v\n", b)
// 		default:
// 			return
// 		}
// 	}
// }
