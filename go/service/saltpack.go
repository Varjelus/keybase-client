// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package service

import (
	"github.com/keybase/client/go/engine"
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol"
	rpc "github.com/keybase/go-framed-msgpack-rpc"
	"golang.org/x/net/context"
)

type SaltpackHandler struct {
	*BaseHandler
	libkb.Contextified
}

type RemoteSaltpackUI struct {
	sessionID int
	cli       keybase1.SaltpackUiClient
}

func NewRemoteSaltpackUI(sessionID int, c *rpc.Client) *RemoteSaltpackUI {
	return &RemoteSaltpackUI{
		sessionID: sessionID,
		cli:       keybase1.SaltpackUiClient{Cli: c},
	}
}

func (r *RemoteSaltpackUI) SaltpackPromptForDecrypt(ctx context.Context, arg keybase1.SaltpackPromptForDecryptArg) (err error) {
	arg.SessionID = r.sessionID
	return r.cli.SaltpackPromptForDecrypt(ctx, arg)
}

func (r *RemoteSaltpackUI) SaltpackVerifySuccess(ctx context.Context, arg keybase1.SaltpackVerifySuccessArg) (err error) {
	arg.SessionID = r.sessionID
	return r.cli.SaltpackVerifySuccess(ctx, arg)
}

func NewSaltpackHandler(xp rpc.Transporter, g *libkb.GlobalContext) *SaltpackHandler {
	return &SaltpackHandler{
		BaseHandler:  NewBaseHandler(xp),
		Contextified: libkb.NewContextified(g),
	}
}

func (h *SaltpackHandler) SaltpackDecrypt(_ context.Context, arg keybase1.SaltpackDecryptArg) (info keybase1.SaltpackEncryptedMessageInfo, err error) {
	cli := h.getStreamUICli()
	src := libkb.NewRemoteStreamBuffered(arg.Source, cli, arg.SessionID)
	snk := libkb.NewRemoteStreamBuffered(arg.Sink, cli, arg.SessionID)
	earg := &engine.SaltpackDecryptArg{
		Sink:   snk,
		Source: src,
		Opts:   arg.Opts,
	}

	ctx := &engine.Context{
		IdentifyUI: h.NewRemoteIdentifyUI(arg.SessionID, h.G()),
		SecretUI:   h.getSecretUI(arg.SessionID),
		SaltpackUI: h.getSaltpackUI(arg.SessionID),
	}
	eng := engine.NewSaltpackDecrypt(earg, h.G())
	err = engine.RunEngine(eng, ctx)
	info = eng.MessageInfo()
	return info, err
}

func (h *SaltpackHandler) SaltpackEncrypt(_ context.Context, arg keybase1.SaltpackEncryptArg) error {
	cli := h.getStreamUICli()
	src := libkb.NewRemoteStreamBuffered(arg.Source, cli, arg.SessionID)
	snk := libkb.NewRemoteStreamBuffered(arg.Sink, cli, arg.SessionID)
	earg := &engine.SaltpackEncryptArg{
		Opts:   arg.Opts,
		Sink:   snk,
		Source: src,
	}

	ctx := &engine.Context{
		IdentifyUI: h.NewRemoteIdentifyUI(arg.SessionID, h.G()),
		SecretUI:   h.getSecretUI(arg.SessionID),
	}
	eng := engine.NewSaltpackEncrypt(earg, h.G())
	return engine.RunEngine(eng, ctx)
}

func (h *SaltpackHandler) SaltpackSign(_ context.Context, arg keybase1.SaltpackSignArg) error {
	cli := h.getStreamUICli()
	src := libkb.NewRemoteStreamBuffered(arg.Source, cli, arg.SessionID)
	snk := libkb.NewRemoteStreamBuffered(arg.Sink, cli, arg.SessionID)
	earg := &engine.SaltpackSignArg{
		Opts:   arg.Opts,
		Sink:   snk,
		Source: src,
	}

	ctx := &engine.Context{
		IdentifyUI: h.NewRemoteIdentifyUI(arg.SessionID, h.G()),
		SecretUI:   h.getSecretUI(arg.SessionID),
	}
	eng := engine.NewSaltpackSign(earg, h.G())
	return engine.RunEngine(eng, ctx)
}

func (h *SaltpackHandler) SaltpackVerify(_ context.Context, arg keybase1.SaltpackVerifyArg) error {
	cli := h.getStreamUICli()
	src := libkb.NewRemoteStreamBuffered(arg.Source, cli, arg.SessionID)
	snk := libkb.NewRemoteStreamBuffered(arg.Sink, cli, arg.SessionID)
	earg := &engine.SaltpackVerifyArg{
		Opts:   arg.Opts,
		Sink:   snk,
		Source: src,
	}

	ctx := &engine.Context{
		IdentifyUI: h.NewRemoteIdentifyUI(arg.SessionID, h.G()),
		SecretUI:   h.getSecretUI(arg.SessionID),
		SaltpackUI: h.getSaltpackUI(arg.SessionID),
	}
	eng := engine.NewSaltpackVerify(earg, h.G())
	return engine.RunEngine(eng, ctx)
}
