// Copyright 2019 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gons

import (
	"io"
	"os"
	"testing"

	"github.com/moby/moby/pkg/reexec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGonsSuite(t *testing.T) {
	// If there was a failure in switching namespaces during inital startup,
	// then report this and end the process with a non-zero status. We do this
	// regardless of whether we're the original test or a reexecuted child.
	if err := Status(); err != nil {
		_, _ = io.WriteString(os.Stderr, err.Error())
		_, _ = io.WriteString(os.Stderr, "\n")
		os.Exit(1)
	}
	// There were no namespace switching errors, so we next register this
	// generic reexecution handler that helps our test procedures. It simply
	// puts the reexecuted child to sleep, waiting to be killed. This allows
	// the parent test to examine the child's namespaces taking all the time
	// it needs in order to figure out if all went well.
	reexec.Register("sleepingunbeauty", func() {
		// Just keep this reexecuted child sleeping; we will be killed by our
		// parent when the test is done. What a lovely family.
		select {}
	})
	// Ensure that the registered handler is run in the reexecuted child. This
	// won't trigger the handler while we're in the parent, because the
	// parent's Arg[0] won't match the name of our handler.
	if !reexec.Init() {
		// Okay, we're a real test suite, and there was no reexecuted child
		// handler triggering... :)
		RegisterFailHandler(Fail)
		RunSpecs(t, "gons suite")
	}
}
