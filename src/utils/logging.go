package utils

import "fmt"

func PrintlnOnDevMode(a ...any) {
        if Env().MODE == "dev" {
        fmt.Println(a...)
        }
}
