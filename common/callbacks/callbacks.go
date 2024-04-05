package xray_common_callbacks

import "github.com/4nd3r5on/Xray-core/common/session"

type OnProcess func(inbound *session.Inbound) error
