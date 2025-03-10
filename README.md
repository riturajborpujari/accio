# Accio
An HTTP(s) downloader written in Go
(_name inspiration from the Harry potter spell `accio` which fetches the requested object_)

## Features
- Allows download via a terminal command
- shows progress as the download happens
- uses only Go standard library to keep things really simple

## Dependencies
- Golang

## Usage
Clone the repo and build with `go build`. the executable named `accio` will be
generated

Install it into a directory mentioned in `PATH` so that it can be invoked from
anywhere.

On Linux / Unix like systems this would be
`install -m755 ./accio /usr/local/bin` 

which would copy the file to the directory `/usr/local/bin` and set the
necessary permissions to allow execution of the program from any user

Now to download url the command would look like `accio URL`

for example running the following will download the 64 bit Alpine linux
installer and place the file as `alpine-standard-3.21.3-x86_64.iso` in the
current directory

`https://dl-cdn.alpinelinux.org/alpine/v3.21/releases/x86_64/alpine-standard-3.21.3-x86_64.iso`

## Implementation overview
Uses `http.Get` to get the file response and a custom copying function which
copies from the response to a file in batches of 4096 bytes. This is done in a
separate goroutine

After every batch the copy function posts status / progress of the download on a
channel which is read from the main goroutine and writes formatted progress /
error onto the terminal

## References
1. https://pkg.go.dev/std
2. https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview
3. https://www2.ccs.neu.edu/research/gpc/VonaUtils/vona/terminal/vtansi.htm
