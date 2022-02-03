//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package healthcheck

import (
	"context"
)

// UseCase defines the operations supported by the document resource
type UseCase interface {
	Get(ctx context.Context) error
}
