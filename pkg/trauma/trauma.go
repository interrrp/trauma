package trauma

import (
	"io"

	"github.com/interrrp/trauma/internal/bytecode"
	"github.com/interrrp/trauma/internal/vm"
)

func Run(program string, reader io.Reader, writer io.Writer) (*vm.VM, error) {
	code, err := bytecode.Compile(program)
	if err != nil {
		return nil, err
	}
	return vm.Run(code, reader, writer)
}
