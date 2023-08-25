// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type MaxBlocks struct{}

func (v MaxBlocks) Visit(ctx context.Context, node hclsyntax.Node, nodeSchema schema.Schema) (diags hcl.Diagnostics) {
	body, ok := node.(*hclsyntax.Body)
	if !ok {
		return
	}

	if nodeSchema == nil {
		return
	}

	foundBlocks := make(map[string]uint64)
	for _, block := range body.Blocks {
		if _, ok := foundBlocks[block.Type]; !ok {
			foundBlocks[block.Type] = 0
		}

		foundBlocks[block.Type]++
	}

	bodySchema := nodeSchema.(*schema.BodySchema)
	for name, blockSchema := range bodySchema.Blocks {
		if blockSchema.MaxItems != 0 {
			foundBlocks, ok := foundBlocks[name]
			if ok && foundBlocks > blockSchema.MaxItems {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf("Too many blocks specified for %q", name),
					Detail:   fmt.Sprintf("Only %d block(s) are expected for %q", blockSchema.MaxItems, name),
					Subject:  node.Range().Ptr(),
				})
			}
		}
	}

	return
}