package trauma

import (
	"io"

	"github.com/interrrp/trauma/internal/bytecode"
	"github.com/interrrp/trauma/internal/vm"
)

func Run(program string, reader io.Reader, writer io.Writer) (*vm.Result, error) {
	bc, err := bytecode.Compile(program)
	if err != nil {
		return nil, err
	}
	return vm.Run(bc, reader, writer)
}
