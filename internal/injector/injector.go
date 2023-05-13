// injector package is responsible for injecting dependencies on runtime.
package injector

import (
	configmod "github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/pkg/hledger"
)

func HledgerClient(config *configmod.Config) hledger.IClient {
	return hledger.NewClient(config.HLedgerExecutable, config.LedgerFile)
}
