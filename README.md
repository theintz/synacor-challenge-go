# synacor-challenge-go
Solving the Synacor challenge in Go. This is a naive implementation of the VM that runs the binary distributed with the challenge.

## build + run
Run the program in the most simple way thinkable: `go run *.go`

## TODO (spoilers)
The binary is actually a text-mode RPG that can be played in the terminal. There are 8 codes hidden in the binary that can be collected. The first 4 are discovered with a working VM implementation (such as this one). The remaining 4 can only be found by playing the game. Others have automated this, maybe in the future I will find time to do that as well.
