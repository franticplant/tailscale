// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

package netns

import "testing"

// Tailmesh: the app's switch must survive netmap updates, which set the node
// attribute value on every netmap.
func TestAppDisableAndroidBindToActiveNetwork(t *testing.T) {
	t.Cleanup(func() {
		SetDisableAndroidBindToActiveNetwork(t.Logf, false)
		SetAppDisableAndroidBindToActiveNetwork(t.Logf, false)
	})
	check := func(step string, wantDisabled, wantByAttr bool) {
		t.Helper()
		disabled, byAttr := AndroidBindToActiveNetworkDisabled()
		if disabled != wantDisabled || byAttr != wantByAttr {
			t.Errorf("%s: disabled=%v byNodeAttr=%v, want %v %v", step, disabled, byAttr, wantDisabled, wantByAttr)
		}
		if disableAndroidBindToActiveNetwork.Load() != wantDisabled {
			t.Errorf("%s: flag read by controlC = %v, want %v", step, !wantDisabled, wantDisabled)
		}
	}

	check("default", false, false)
	SetAppDisableAndroidBindToActiveNetwork(t.Logf, true)
	check("app on", true, false)
	SetDisableAndroidBindToActiveNetwork(t.Logf, false) // a netmap without the attribute
	check("app on, netmap without attr", true, false)
	SetDisableAndroidBindToActiveNetwork(t.Logf, true)
	check("app on, attr on", true, true)
	SetAppDisableAndroidBindToActiveNetwork(t.Logf, false)
	check("app off, attr on", true, true)
	SetDisableAndroidBindToActiveNetwork(t.Logf, false)
	check("both off", false, false)
}
