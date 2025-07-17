package shadowsocks_callbacks

import (
	callbacksCommon "github.com/xtls/xray-core/common/callbacks"
	"github.com/xtls/xray-core/common/idsyncmap"
	"github.com/xtls/xray-core/common/session"
)

type ServerCallbackManager struct {
	CbsOnProcess idsyncmap.IDSyncMap[int, callbacksCommon.InboundOnProcess]
}

func (cm *ServerCallbackManager) ExecOnProcess(inbound *session.Inbound) (id int, err error) {
	for id, callback := range cm.CbsOnProcess.GetAll() {
		err = callback(inbound)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func NewServerCallbackManager() *ServerCallbackManager {
	return &ServerCallbackManager{
		CbsOnProcess: idsyncmap.NewIDSyncMap[int, callbacksCommon.InboundOnProcess](),
	}
}
