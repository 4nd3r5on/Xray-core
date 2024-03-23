package httpupgrade

import (
	"context"

	"github.com/4nd3r5on/Xray-core/common"
)

//go:generate go run github.com/4nd3r5on/Xray-core/common/errors/errorgen

const protocolName = "httpupgrade"

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return nil, newError("httpupgrade is a transport protocol.")
	}))
}
