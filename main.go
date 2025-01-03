// main.go
package main

import (
    "log"

    "github.com/zamorofthat/cribl-storage-tool/cmd"
)

func main() {
    rootCmd := cmd.NewRootCmd()
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}
