package xray_common_callbacks

import "github.com/xtls/xray-core/common/session"

type OnProcess func(inbound *session.Inbound) error
