# Trauma

Trauma is a very fast Brainfuck interpreter written in Go.

## Installation

```sh
go get github.com/interrrp/trauma/cmd/trauma
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

- Compiles programs into fast bytecode with specialized instructions:

  | Bytecode    | Brainfuck    |
  | ----------- | ------------ |
  | `Inc 1`     | `+++--`      |
  | `IncPtr -2` | `>><<<<`     |
  | `Clear`     | `[-]`, `[+]` |
  | `Move -1`   | `[-<+>]`     |

- Infinite tape size

## Examples

```sh
go run ./cmd/trauma ./programs/mandelbrot.b
```

## API

```sh
go get github.com/interrrp/trauma/pkg/trauma
```

```go
import (
  "os"

  "github.com/interrrp/trauma/pkg/trauma"
)

func main() {
  res, err := trauma.Run("+++", os.Stdin, os.Stdout)
  fmt.Println(res)
}
```

## Development

Run tests:

```sh
go test ./...
```

## License

MIT
