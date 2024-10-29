package memory

import (
	"github.com/gera-cl/xk6-working-memory/memory"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register(memory.ImportPath, memory.New())
}
