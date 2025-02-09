// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package engine

import (
	"github.com/jonboulle/clockwork"
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol"
	"strings"
	"testing"
	"time"
)

type flakeyRooterAPI struct {
	orig     libkb.ExternalAPI
	flakeOut bool
}

func newFlakeyRooterAPI(x libkb.ExternalAPI) *flakeyRooterAPI {
	return &flakeyRooterAPI{
		orig: x,
	}
}

func (e *flakeyRooterAPI) GetText(arg libkb.APIArg) (*libkb.ExternalTextRes, error) {
	return e.orig.GetText(arg)
}

func (e *flakeyRooterAPI) Get(arg libkb.APIArg) (res *libkb.ExternalAPIRes, err error) {

	// Show an error if we're in flakey mode
	if e.flakeOut && strings.Contains(arg.Endpoint, "rooter") {
		return &libkb.ExternalAPIRes{HTTPStatus: 429}, &libkb.APIError{Msg: "Ratelimited", Code: 429}
	}

	return e.orig.Get(arg)
}

func (e *flakeyRooterAPI) GetHTML(arg libkb.APIArg) (res *libkb.ExternalHTMLRes, err error) {
	return e.orig.GetHTML(arg)
}

func (e *flakeyRooterAPI) Post(arg libkb.APIArg) (res *libkb.ExternalAPIRes, err error) {
	return e.orig.Post(arg)
}

func (e *flakeyRooterAPI) PostHTML(arg libkb.APIArg) (res *libkb.ExternalHTMLRes, err error) {
	return e.orig.PostHTML(arg)
}

func TestSoftSnooze(t *testing.T) {
	tc := SetupEngineTest(t, "track")
	defer tc.Cleanup()
	fu := CreateAndSignupFakeUser(tc, "track")

	fakeClock := clockwork.NewFakeClock()
	tc.G.Clock = fakeClock
	// to pick up the new clock...
	tc.G.ResetLoginState()

	flakeyAPI := flakeyRooterAPI{orig: tc.G.XAPI, flakeOut: false}
	tc.G.XAPI = &flakeyAPI

	idUI := &FakeIdentifyUI{}
	username := "t_tracy"
	arg := &keybase1.Identify2Arg{
		UserAssertion: username,
		NeedProofSet:  true,
	}
	ctx := &Context{
		LogUI:      tc.G.UI.GetLogUI(),
		IdentifyUI: idUI,
		SecretUI:   fu.NewSecretUI(),
	}

	// Identify tracy; all proofs should work
	eng := NewResolveThenIdentify2(tc.G, arg)
	if err := RunEngine(eng, ctx); err != nil {
		t.Fatal(err)
	}
	targ := TrackTokenArg{
		Token:   idUI.Token,
		Options: keybase1.TrackOptions{BypassConfirm: true},
	}

	// Track tracy
	teng := NewTrackToken(&targ, tc.G)
	if err := RunEngine(teng, ctx); err != nil {
		t.Fatal(err)
	}

	// Now make her Rooter proof flakey / fail with a 429
	flakeyAPI.flakeOut = true
	idUI = &FakeIdentifyUI{}
	ctx.IdentifyUI = idUI

	// Advance so that our previous cached success is out of
	// cache on its own, but still can override a 429-like soft failure.
	fakeClock.Advance(tc.G.Env.GetProofCacheMediumDur() + time.Minute)

	eng = NewResolveThenIdentify2(tc.G, arg)
	eng.testArgs = &Identify2WithUIDTestArgs{noCache: true}
	// Should not get an error
	if err := RunEngine(eng, ctx); err != nil {
		t.Fatal(err)
	}
	result, found := idUI.ProofResults["rooter"]
	if !found {
		t.Fatal("Failed to find a rooter proof")
	}
	if pe := libkb.ImportProofError(result.SnoozedResult); pe == nil {
		t.Fatal("expected a snoozed error result")
	}

	// Now time out the success that allowed us to circumvent
	// the soft failure.
	fakeClock.Advance(tc.G.Env.GetProofCacheLongDur())
	eng = NewResolveThenIdentify2(tc.G, arg)
	eng.testArgs = &Identify2WithUIDTestArgs{noCache: true}
	if err := RunEngine(eng, ctx); err == nil {
		t.Fatal("Expected a failure in our proof")
	}

	result, found = idUI.ProofResults["rooter"]
	if !found {
		t.Fatal("Failed to find a rooter proof")
	}
	if pe := libkb.ImportProofError(result.ProofResult); pe == nil {
		t.Fatal("expected a Rooter error result")
	}

	defer runUntrack(tc.G, fu, username)
	assertTracking(tc, username)
}
