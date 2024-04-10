//go:generate sqlc generate

package main

import "github.com/cs5224virgo/virgo/cmd"

func main() {
	cmd.Execute()
}
