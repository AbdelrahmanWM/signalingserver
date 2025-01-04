//go:build js && wasm

// +build: js,wasm
package utils

import (
	"fmt"
	"log"
	"syscall/js"
)
/// Logging into html div
func Log(value string) {
	log.Println(value)
	el := GetElementByID("logArea")
	el.Set("innerHTML", el.Get("innerHTML").String()+"* "+value+"<br>")
}

func GetElementByID(id string) js.Value {
	el := js.Global().Get("document").Call("getElementById", id)
	if el.IsNull() {
		Log(fmt.Sprintf("Element with id '%s' not found", id))
	}
	return el
}
func ClearLog(v js.Value, p []js.Value) any {
	el := GetElementByID("logArea")
	el.Set("innerHTML", "")
	return nil
}
