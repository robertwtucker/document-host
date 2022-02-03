//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package healthcheck

import "context"

// DatabaseHelper defines the operations supported by the healthcheck resource
type DatabaseHelper interface {
	CheckDB(ctx context.Context) error
}
