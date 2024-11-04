// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ptraceutil_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/open-telemetry/opentelemetry-collector-contrib/connector/routingconnector/internal/ptraceutil"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/ptracetest"
)

func TestMoveResourcesIf(t *testing.T) {
	testCases := []struct {
		name      string
		condition func(ptrace.ResourceSpans) bool
	}{
		{
			name: "move_none",
			condition: func(ptrace.ResourceSpans) bool {
				return false
			},
		},
		{
			name: "move_all",
			condition: func(ptrace.ResourceSpans) bool {
				return true
			},
		},
		{
			name: "move_one",
			condition: func(rl ptrace.ResourceSpans) bool {
				rname, ok := rl.Resource().Attributes().Get("resourceName")
				return ok && rname.AsString() == "resourceA"
			},
		},
		{
			name: "move_to_preexisting",
			condition: func(rl ptrace.ResourceSpans) bool {
				rname, ok := rl.Resource().Attributes().Get("resourceName")
				return ok && rname.AsString() == "resourceB"
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Load up a fresh copy of the input for each test, since it may be modified in place.
			from, err := golden.ReadTraces(filepath.Join("testdata", "resource", tt.name, "from.yaml"))
			require.NoError(t, err)

			to, err := golden.ReadTraces(filepath.Join("testdata", "resource", tt.name, "to.yaml"))
			require.NoError(t, err)

			fromModifed, err := golden.ReadTraces(filepath.Join("testdata", "resource", tt.name, "from_modified.yaml"))
			require.NoError(t, err)

			toModified, err := golden.ReadTraces(filepath.Join("testdata", "resource", tt.name, "to_modified.yaml"))
			require.NoError(t, err)

			ptraceutil.MoveResourcesIf(from, to, tt.condition)

			assert.NoError(t, ptracetest.CompareTraces(fromModifed, from), "from not modified as expected")
			assert.NoError(t, ptracetest.CompareTraces(toModified, to), "to not as expected")
		})
	}
}
