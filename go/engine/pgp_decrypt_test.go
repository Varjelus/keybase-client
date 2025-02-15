// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package engine

import (
	"bytes"
	"crypto/rand"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol"
)

func decengctx(fu *FakeUser, tc libkb.TestContext) *Context {
	return &Context{
		IdentifyUI: &FakeIdentifyUI{},
		SecretUI:   fu.NewSecretUI(),
		LogUI:      tc.G.UI.GetLogUI(),
		PgpUI:      &TestPgpUI{},
	}
}

func TestPGPDecrypt(t *testing.T) {
	tc := SetupEngineTest(t, "PGPDecrypt")
	defer tc.Cleanup()
	fu := createFakeUserWithPGPOnly(t, tc)

	// encrypt a message
	msg := "10 days in Japan"
	sink := libkb.NewBufferCloser()
	ctx := decengctx(fu, tc)
	arg := &PGPEncryptArg{
		Source:       strings.NewReader(msg),
		Sink:         sink,
		NoSign:       true,
		BinaryOutput: true,
	}
	enc := NewPGPEncrypt(arg, tc.G)
	if err := RunEngine(enc, ctx); err != nil {
		t.Fatal(err)
	}
	out := sink.Bytes()

	t.Logf("encrypted data: %x", out)

	// decrypt it
	decoded := libkb.NewBufferCloser()
	decarg := &PGPDecryptArg{
		Source: bytes.NewReader(out),
		Sink:   decoded,
	}
	dec := NewPGPDecrypt(decarg, tc.G)
	if err := RunEngine(dec, ctx); err != nil {
		t.Fatal(err)
	}
	decmsg := string(decoded.Bytes())
	if decmsg != msg {
		t.Errorf("decoded: %q, expected: %q", decmsg, msg)
	}

	owner := dec.Owner()
	if owner == nil {
		t.Errorf("owner is nil")
	}
}

func TestPGPDecryptArmored(t *testing.T) {
	tc := SetupEngineTest(t, "PGPDecrypt")
	defer tc.Cleanup()
	fu := createFakeUserWithPGPOnly(t, tc)

	// encrypt a message
	msg := "10 days in Japan"
	ctx := decengctx(fu, tc)
	sink := libkb.NewBufferCloser()
	arg := &PGPEncryptArg{
		Source: strings.NewReader(msg),
		Sink:   sink,
		NoSign: true,
	}
	enc := NewPGPEncrypt(arg, tc.G)
	if err := RunEngine(enc, ctx); err != nil {
		t.Fatal(err)
	}
	out := sink.Bytes()

	t.Logf("encrypted data: %x", out)

	// decrypt it
	decoded := libkb.NewBufferCloser()
	decarg := &PGPDecryptArg{
		Source: bytes.NewReader(out),
		Sink:   decoded,
	}
	dec := NewPGPDecrypt(decarg, tc.G)
	if err := RunEngine(dec, ctx); err != nil {
		t.Fatal(err)
	}
	decmsg := string(decoded.Bytes())
	if decmsg != msg {
		t.Errorf("decoded: %q, expected: %q", decmsg, msg)
	}

	// A saltpack message
	saltpack := `BEGIN KEYBASE SALTPACK ENCRYPTED MESSAGE.
ZUHRHckf9VJ6ich bKthcf8yieGqdWj H4gqRvnQEBYi4Lw JzZL5NnUd2ssiZ3 hReptn2rUJuGcna uAkp7yQGAaxVQyc oM2o7JykpRwsGbM fmF4r2Mj1aCAKd6 U316QFRCxkC4Uik JszTub0Mt1a1IRi wQNIl3ru791WSMK JRXkD0biElgMSIf Qk9B0j1vRHeLB11 Ig9qX7TBMzg9InW tHfZp0GjYEJWCvZ i0h3SG99wTPL3Ov PSOjc2oGRqI2bya XaGdao8vmfYk8jd auqe5QBgA2RaO7I WYZGyHwsWNg3PSX eEcUtccZeD7ryBP MMDvFZUDUia1Njf ELaaKo932JWKdLg zONmMEinFizgGEg 7MrB4kIjeY3O4Od GEjKV6okfzdmYOP rZCsNzFILna2K1h D53RmTDkoddt81Y mwPFWwaeCA9YR9S EZcu8kWlaFjjDRX SG7MCuwtS3PuUXr EaCfe1ib65Nrq3s IUxGJ7H6JRidavn Ql7Bo9qbZGqkglf 0g1caODxGe4mQfZ ixH7my73DAKw5aL idukKTwb0qPVlvD 2vqskVSukRNra5i 1t1PadrLZSgaJqb WtSVmbgiWm40P1S WlR8nq2I95RIZSx LP3JEvHqBNAb9Ci rzkkSaOBk5FawpE yCVbUDk556V0e7F z8YSYPFqzwf14Yo JEztr4noRMxvME2 OBTt16BAJF4K1NA pTtFpwRbubTggxb 74abisSK1DgPd30 UgdY75zTUKd57pk CHTu2BwPHFGjwgf XxJMTNpYRQiz6uZ lIAmcTKhGoKBy97 7S82DvT1tB0cGjn M5JnyMJkzj2WaVf MFkAonXbkYNOk2n olE7RldZhcyy4Xd edu1Mke8AevVaAc lat6mso7hS2XAuW ZCcrWakNFGPPqkW 40YMZMHqL1mbIIS oooC5lSP2dd2c04 j9yjhpsYjs2izOA HRhQslExQDU6Uio WEJBnQDpMbhQG1y E9jPqsSgmUrfRgD nEuLPYpRLm6UJgW TZFga6U0khClKRX DdTWaLVGrmY8xdB GNy2Dd4HStxQ4PF vwTIyDUee5loag1 ePdbLqPMdkh19zR PcBFK4gtKCdhzIu ea2Ncg69SqseUBi wp45MaNEGlIh9Y7 sU85K0nKEMK5if3 7HOGDDVCL7pwmNj 9A6DKk8MrwpPTeB 2B3uqQpHsKIdPJx qD7S2IJahokxoiY 5UPvodqZY8JRljz In7rV6I5LUyh5SM tSL4t0Z2VfrMDZH En5QqKvVJykQU9S ELH1U2Hxh7ANzCK v2R3xrf102D1zaG 3sO6yLwzpZH2Rq0 q3h9GbWIEndPZHA IsPJC2MFfN0sOwe e6nEuqR2NlsDwBk hMWLszNY7iOICmE RvwFZXFqTwbRlxA qBoAPZYJyPFdKpV MsSRHrFUIwYTE6S ZPQ9bRmb2B1hAu9 rkECh80CmFO04Fc rn1KX392TtgIYu1 PN1LGcrAYdD6UC1 O9Vx2fuBiTIqmo9 XbPsWdRxmX57BjS EZPTZ9wGaxLZVqB cmsDn1mU1Uzbesx 2pXT0mVO7A72mkw wBdlD7QDUeN8Na7 j8W9tUIWUAbAePO 2Z9OSU0M1KIRSuE ePOZBlNonU1dUCz KQWlw.
END KEYBASE SALTPACK ENCRYPTED MESSAGE.`

	decoded = libkb.NewBufferCloser()
	decarg = &PGPDecryptArg{
		Source: strings.NewReader(saltpack),
		Sink:   decoded,
	}
	dec = NewPGPDecrypt(decarg, tc.G)
	err := RunEngine(dec, ctx)
	if wse, ok := err.(libkb.WrongCryptoFormatError); !ok {
		t.Fatalf("Wanted a WrongCryptoFormat error, but got %T (%v)", err, err)
	} else if wse.Wanted != libkb.CryptoMessageFormatPGP ||
		wse.Received != libkb.CryptoMessageFormatSaltpack ||
		wse.Operation != "decrypt" {
		t.Fatalf("Bad error: %v", wse)
	}
}

// TestPGPDecryptSignedSelf tests that the user who signed the
// message can decrypt it.
func TestPGPDecryptSignedSelf(t *testing.T) {
	tc := SetupEngineTest(t, "PGPDecrypt")
	defer tc.Cleanup()
	fu := createFakeUserWithPGPOnly(t, tc)

	// encrypt a message
	msg := "We pride ourselves on being meticulous; no issue is too small."
	ctx := decengctx(fu, tc)
	sink := libkb.NewBufferCloser()
	arg := &PGPEncryptArg{
		Source:       strings.NewReader(msg),
		Sink:         sink,
		BinaryOutput: true,
	}
	enc := NewPGPEncrypt(arg, tc.G)
	if err := RunEngine(enc, ctx); err != nil {
		t.Fatal(err)
	}
	out := sink.Bytes()

	t.Logf("encrypted data: %x", out)

	// decrypt it
	decoded := libkb.NewBufferCloser()
	decarg := &PGPDecryptArg{
		Source:       bytes.NewReader(out),
		Sink:         decoded,
		AssertSigned: true,
	}
	dec := NewPGPDecrypt(decarg, tc.G)
	if err := RunEngine(dec, ctx); err != nil {
		t.Fatal(err)
	}
	decmsg := string(decoded.Bytes())
	if decmsg != msg {
		t.Errorf("decoded: %q, expected: %q", decmsg, msg)
	}
}

// TestPGPDecryptSignedOther tests that a user who didn't sign the
// message can verify the signature.
func TestPGPDecryptSignedOther(t *testing.T) {
	tcRecipient := SetupEngineTest(t, "PGPDecrypt - Recipient")
	defer tcRecipient.Cleanup()
	recipient := createFakeUserWithPGPSibkey(tcRecipient)
	Logout(tcRecipient)

	tcSigner := SetupEngineTest(t, "PGPDecrypt - Signer")
	defer tcSigner.Cleanup()
	signer := createFakeUserWithPGPSibkey(tcSigner)

	// encrypt a message
	msg := "We pride ourselves on being meticulous; no issue is too small."
	ctx := decengctx(signer, tcSigner)
	sink := libkb.NewBufferCloser()
	arg := &PGPEncryptArg{
		Recips:       []string{recipient.Username},
		Source:       strings.NewReader(msg),
		Sink:         sink,
		BinaryOutput: true,
		TrackOptions: keybase1.TrackOptions{BypassConfirm: true},
	}
	enc := NewPGPEncrypt(arg, tcSigner.G)
	if err := RunEngine(enc, ctx); err != nil {
		t.Fatal(err)
	}
	out := sink.Bytes()

	t.Logf("encrypted data: %x", out)

	// signer logs out, recipient logs in:
	t.Logf("signer (%q) logging out", signer.Username)
	Logout(tcSigner)
	libkb.G = tcRecipient.G
	// G = libkb.G
	t.Logf("recipient (%q) logging in", recipient.Username)
	recipient.LoginOrBust(tcRecipient)

	rtrackUI := &FakeIdentifyUI{}
	ctx = &Context{
		IdentifyUI: rtrackUI,
		SecretUI:   recipient.NewSecretUI(),
		LogUI:      tcRecipient.G.UI.GetLogUI(),
		PgpUI:      &TestPgpUI{},
	}

	// decrypt it
	decoded := libkb.NewBufferCloser()
	decarg := &PGPDecryptArg{
		Source:       bytes.NewReader(out),
		Sink:         decoded,
		AssertSigned: true,
	}
	dec := NewPGPDecrypt(decarg, tcRecipient.G)
	if err := RunEngine(dec, ctx); err != nil {
		t.Fatal(err)
	}
	decmsg := string(decoded.Bytes())
	if decmsg != msg {
		t.Errorf("decoded: %q, expected: %q", decmsg, msg)
	}
}

// TestPGPDecryptSignedIdentify tests that the signer is
// identified regardless of AssertSigned, SignedBy args.
func TestPGPDecryptSignedIdentify(t *testing.T) {
	tcRecipient := SetupEngineTest(t, "PGPDecrypt - Recipient")
	defer tcRecipient.Cleanup()
	recipient := createFakeUserWithPGPSibkey(tcRecipient)
	Logout(tcRecipient)

	tcSigner := SetupEngineTest(t, "PGPDecrypt - Signer")
	defer tcSigner.Cleanup()
	signer := createFakeUserWithPGPSibkey(tcSigner)

	// encrypt a message
	msg := "We pride ourselves on being meticulous; no issue is too small."
	ctx := decengctx(signer, tcSigner)
	sink := libkb.NewBufferCloser()
	arg := &PGPEncryptArg{
		Recips:       []string{recipient.Username},
		Source:       strings.NewReader(msg),
		Sink:         sink,
		BinaryOutput: true,
		TrackOptions: keybase1.TrackOptions{BypassConfirm: true},
	}
	enc := NewPGPEncrypt(arg, tcSigner.G)
	if err := RunEngine(enc, ctx); err != nil {
		t.Fatal(err)
	}
	out := sink.Bytes()

	t.Logf("encrypted data: %x", out)

	// signer logs out, recipient logs in:
	t.Logf("signer (%q) logging out", signer.Username)
	Logout(tcSigner)
	libkb.G = tcRecipient.G
	t.Logf("recipient (%q) logging in", recipient.Username)
	recipient.LoginOrBust(tcRecipient)

	idUI := &FakeIdentifyUI{}
	pgpUI := &TestPgpUI{}
	ctx = &Context{
		IdentifyUI: idUI,
		SecretUI:   recipient.NewSecretUI(),
		LogUI:      tcRecipient.G.UI.GetLogUI(),
		PgpUI:      pgpUI,
	}

	// decrypt it
	decoded := libkb.NewBufferCloser()
	decarg := &PGPDecryptArg{
		Source:       bytes.NewReader(out),
		Sink:         decoded,
		AssertSigned: false,
	}
	dec := NewPGPDecrypt(decarg, tcRecipient.G)
	if err := RunEngine(dec, ctx); err != nil {
		t.Fatal(err)
	}

	if idUI.User == nil {
		t.Fatal("identify ui user is nil")
	}
	if idUI.User.Username != signer.Username {
		t.Errorf("idUI username: %q, expected %q", idUI.User.Username, signer.Username)
	}
	if pgpUI.OutputCount != 1 {
		t.Errorf("PgpUI output called %d times, expected 1", pgpUI.OutputCount)
	}
}

func TestPGPDecryptLong(t *testing.T) {
	tc := SetupEngineTest(t, "PGPDecrypt")
	defer tc.Cleanup()
	fu := createFakeUserWithPGPSibkey(tc)

	// encrypt a message
	msg := make([]byte, 1024*1024)

	if _, err := rand.Read(msg); err != nil {
		t.Fatal(err)
	}

	sink := libkb.NewBufferCloser()
	ctx := decengctx(fu, tc)
	arg := &PGPEncryptArg{
		Source:       bytes.NewReader(msg),
		Sink:         sink,
		NoSign:       true,
		BinaryOutput: true,
	}
	enc := NewPGPEncrypt(arg, tc.G)
	if err := RunEngine(enc, ctx); err != nil {
		t.Fatal(err)
	}
	out := sink.Bytes()

	// decrypt it
	decoded := libkb.NewBufferCloser()
	decarg := &PGPDecryptArg{
		Source: bytes.NewReader(out),
		Sink:   decoded,
	}
	dec := NewPGPDecrypt(decarg, tc.G)
	if err := RunEngine(dec, ctx); err != nil {
		t.Fatal(err)
	}
	decmsg := decoded.Bytes()
	if len(decmsg) != len(msg) {
		t.Fatalf("decoded msg size: %d, expected %d", len(decmsg), len(msg))
	}

	for i, b := range msg {
		if decmsg[i] != b {
			t.Errorf("decode msg differs at byte %d: %x, expected %x", i, decmsg[i], b)
		}
	}

	owner := dec.Owner()
	if owner == nil {
		t.Errorf("owner is nil")
	}
}

type cstest struct {
	name string
	msg  string
}

var cstests = []cstest{
	{name: "ascii", msg: "hello"},
	{name: "emoji", msg: "😓😕😙"},
}

func TestPGPDecryptClearsign(t *testing.T) {
	tc := SetupEngineTest(t, "PGPDecrypt")
	defer tc.Cleanup()

	fu := createFakeUserWithPGPSibkey(tc)
	ctx := decengctx(fu, tc)

	for _, test := range cstests {
		signedMsg := sign(ctx, tc, test.msg, keybase1.SignMode_CLEAR)
		t.Logf("%s: signed message:\n\n%s\n", test.name, signedMsg)

		decoded := libkb.NewBufferCloser()
		arg := &PGPDecryptArg{
			Source: strings.NewReader(signedMsg),
			Sink:   decoded,
		}
		eng := NewPGPDecrypt(arg, tc.G)
		if err := RunEngine(eng, ctx); err != nil {
			t.Errorf("%s: decrypt error: %q", test.name, err)
			continue
		}
		msg := decoded.Bytes()
		trimmed := strings.TrimSpace(string(msg))
		t.Logf("clearsign test %q decoded message: %s\n", test.name, trimmed)
		if trimmed != test.msg {
			t.Errorf("%s: expected msg %q, got %q", test.name, test.msg, trimmed)
		}

		status := eng.SignatureStatus()
		if !status.IsSigned {
			t.Errorf("%s: expected IsSigned", test.name)
		}
		if !status.Verified {
			t.Errorf("%s: expected Verified", test.name)
		}
		if status.Entity == nil {
			t.Errorf("%s: signature status entity is nil", test.name)
		}
	}
}

type TestPgpUI struct {
	OutputCount int
}

func (t *TestPgpUI) OutputSignatureSuccess(context.Context, keybase1.OutputSignatureSuccessArg) error {
	t.OutputCount++
	return nil
}
