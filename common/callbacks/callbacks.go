package callbacks

import "github.com/xtls/xray-core/common/session"

type InboundOnProcess func(inbound *session.Inbound) error
