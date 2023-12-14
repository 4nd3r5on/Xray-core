package shadowsocks

import (
	"github.com/4nd3r5on/Xray-core/common/idsyncmap"
	"github.com/4nd3r5on/Xray-core/common/protocol"
	"github.com/4nd3r5on/Xray-core/common/session"
	"github.com/4nd3r5on/Xray-core/features/policy"
)

type (
	OnProcess struct {
		Exec func() error
	}
	OnHandleConn struct {
		Exec func(inbound *session.Inbound, sessionPolicy policy.Session) error
	}
	OnHandleUDP struct {
		Exec func(inbound *session.Inbound, sessionPolicy policy.Session) error
	}
)

type CallbackManager struct {
	CbsOnProcess    idsyncmap.IDSyncMap[OnProcess]
	CbsOnHandleConn idsyncmap.IDSyncMap[OnHandleConn]
	CbsOnHandleUDP  idsyncmap.IDSyncMap[OnHandleUDP]
}

func (cm *CallbackManager) ExecOnProcess() (id int32, err error) {
	for id, callback := range cm.CbsOnProcess.Get() {
		err = callback.Exec()
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func (cm *CallbackManager) ExecOnHandleConn(inbound *session.Inbound, sessionPolicy policy.Session) (id int32, err error) {
	for id, callback := range cm.CbsOnProcess.Get() {
		err = callback.Exec()
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

// TODO: Find where to use
func (cm *CallbackManager) ExecOnHandleUDP(request *protocol.RequestHeader, sessionPolicy policy.Session) (id int32, err error) {
	for id, callback := range cm.CbsOnProcess.Get() {
		err = callback.Exec()
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func NewCallbackManager() *CallbackManager {
	return &CallbackManager{
		CbsOnProcess:    idsyncmap.NewIDSyncMap[OnProcess](),
		CbsOnHandleConn: idsyncmap.NewIDSyncMap[OnHandleConn](),
		CbsOnHandleUDP:  idsyncmap.NewIDSyncMap[OnHandleUDP](),
	}
}
