package core

import (
	"encoding/json"
	"strings"

	"github.com/currycan/supkube/pkg/logger"
)

//Print is
func (s *Installer) Print(process ...string) {
	if len(process) == 0 {
		configJSON, _ := json.Marshal(s)
		logger.Info("\n[globals] config is: ", string(configJSON))
	} else {
		var sb strings.Builder
		for _, v := range process {
			sb.Write([]byte("==>"))
			sb.Write([]byte(v))
		}
		logger.Debug(sb.String())
	}
}
func (s *Installer) PrintFinish() {
	logger.Info(" install success.")
}
