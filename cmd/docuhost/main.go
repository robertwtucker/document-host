//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package main

import (
	_ "github.com/robertwtucker/document-host/cmd" // makes other commands visible
	"github.com/robertwtucker/document-host/cmd/root"
)

func main() {
	root.Execute()
}
