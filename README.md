# Trauma

Trauma is a Brainfuck interpreter written in Go.

## Installation

```sh
go get github.com/interrrp/trauma
```

## Usage

```sh
trauma path
```

## Features

- Full Brainfuck language support:

  - `+` Increment value at current cell
  - `-` Decrement value at current cell
  - `>` Move pointer right
  - `<` Move pointer left
  - `[` Jump past matching `]` if current cell is 0
  - `]` Jump back to matching `[` if current cell is not 0
  - `.` Output byte at current cell
  - `,` Read byte into current cell

- Compiles programs into fast bytecode
- Infinite tape size

## Examples

```sh
go run . programs/mandelbrot.b
```

## Development

Run tests:

```sh
go test ./...
```

## License

MIT
