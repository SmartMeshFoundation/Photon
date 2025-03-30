// SPDX-FileCopyrightText: 2021 The Go-SSB Authors
//
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"go.cryptoscope.co/muxrpc/v2"
	cli "gopkg.in/urfave/cli.v2"

	refs "go.mindeco.de/ssb-refs"
)

var publishCmd = &cli.Command{
	Name:  "publish",
	Usage: "p",
	Subcommands: []*cli.Command{
		publishRawCmd,
		publishPostCmd,
		publishAboutCmd,
		publishContactCmd,
		publishVoteCmd,
	},
}

var publishRawCmd = &cli.Command{
	Name:      "raw",
	UsageText: "reads JSON from stdin and publishes that as content",
	// TODO: add private

	Action: func(ctx *cli.Context) error {
		var content interface{}
		err := json.NewDecoder(os.Stdin).Decode(&content)
		if err != nil {
			return fmt.Errorf("publish/raw: invalid json input from stdin: %w", err)
		}

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var v string
		err = client.Async(longctx, &v, muxrpc.TypeString, muxrpc.Method{"publish"}, content)
		if err != nil {
			return fmt.Errorf("publish call failed: %w", err)
		}
		newMsg, err := refs.ParseMessageRef(v)
		if err != nil {
			return err
		}
		log.Log("event", "published", "type", "raw", "ref", newMsg.String())
		fmt.Fprintln(os.Stdout, newMsg.String())
		return nil
	},
}

var publishPostCmd = &cli.Command{
	Name:      "post",
	ArgsUsage: "text of the post",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "root", Value: "", Usage: "the ID of the first message of the thread"},
		// TODO: Slice of branches
		&cli.StringFlag{Name: "branch", Value: "", Usage: "the post ID that is beeing replied to"},

		&cli.StringSliceFlag{Name: "recps", Usage: "as a PM to these feeds"},
	},
	Action: func(ctx *cli.Context) error {
		arg := map[string]interface{}{
			"text": ctx.Args().First(),
			"type": "post",
		}
		if r := ctx.String("root"); r != "" {
			arg["root"] = r
			if b := ctx.String("branch"); b != "" {
				arg["branch"] = b
			} else {
				arg["branch"] = r
			}
		}

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var v string
		if recps := ctx.StringSlice("recps"); len(recps) > 0 {
			err = client.Async(longctx, &v,
				muxrpc.TypeString,
				muxrpc.Method{"private", "publish"}, arg, recps)
		} else {
			err = client.Async(longctx, &v,
				muxrpc.TypeString,
				muxrpc.Method{"publish"}, arg)
		}
		if err != nil {
			return fmt.Errorf("publish call failed: %w", err)
		}

		newMsg, err := refs.ParseMessageRef(v)
		if err != nil {
			return err
		}
		log.Log("event", "published", "type", "post", "ref", newMsg.String())
		fmt.Fprintln(os.Stdout, newMsg.String())
		return nil
	},
}

var publishVoteCmd = &cli.Command{
	Name:      "vote",
	ArgsUsage: "%linkedMessage.sha256",
	Flags: []cli.Flag{
		&cli.IntFlag{Name: "value", Usage: "usually 1 (like) or 0 (unlike)"},
		&cli.StringFlag{Name: "expression", Usage: "Dig/Yup/Heart"},

		&cli.StringFlag{Name: "root", Value: "", Usage: "the ID of the first message of the thread"},
		// TODO: Slice of branches
		&cli.StringFlag{Name: "branch", Value: "", Usage: "the post ID that is beeing replied to"},

		&cli.StringSliceFlag{Name: "recps", Usage: "as a PM to these feeds"},
	},
	Action: func(ctx *cli.Context) error {
		mref, err := refs.ParseMessageRef(ctx.Args().First())
		if err != nil {
			return fmt.Errorf("publish/vote: invalid msg ref: %w", err)
		}

		arg := map[string]interface{}{
			"vote": map[string]interface{}{
				"link":       mref.String(),
				"value":      ctx.Int("value"),
				"expression": ctx.String("expression"),
			},
			"type": "vote",
		}

		if r := ctx.String("root"); r != "" {
			arg["root"] = r
			if b := ctx.String("branch"); b != "" {
				arg["branch"] = b
			} else {
				arg["branch"] = r
			}
		}

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var v string
		if recps := ctx.StringSlice("recps"); len(recps) > 0 {
			err = client.Async(longctx, &v,
				muxrpc.TypeString,
				muxrpc.Method{"private", "publish"}, arg, recps)
		} else {
			err = client.Async(longctx, &v,
				muxrpc.TypeString,
				muxrpc.Method{"publish"}, arg)
		}
		if err != nil {
			return fmt.Errorf("publish call failed: %w", err)
		}

		newMsg, err := refs.ParseMessageRef(v)
		if err != nil {
			return err
		}

		log.Log("event", "published", "type", "vote", "ref", newMsg.String())
		fmt.Fprintln(os.Stdout, newMsg.String())
		return nil
	},
}

var publishAboutCmd = &cli.Command{
	Name:      "about",
	ArgsUsage: "@aboutkeypair.ed25519",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Usage: "what name to give"},
		&cli.StringFlag{Name: "image", Usage: "image blob ref"},
	},
	Action: func(ctx *cli.Context) error {
		aboutRef, err := refs.ParseFeedRef(ctx.Args().First())
		if err != nil {
			return fmt.Errorf("publish/about: invalid feed ref: %w", err)
		}
		arg := map[string]interface{}{
			"about": aboutRef.String(),
			"type":  "about",
		}
		if n := ctx.String("name"); n != "" {
			arg["name"] = n
		}
		if img := ctx.String("image"); img != "" {
			blobRef, err := refs.ParseBlobRef(img)
			if err != nil {
				return fmt.Errorf("publish/about: invalid blob ref: %w", err)
			}
			arg["image"] = blobRef
		}

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var v string
		err = client.Async(longctx, &v, muxrpc.TypeString, muxrpc.Method{"publish"}, arg)
		if err != nil {
			return fmt.Errorf("publish call failed: %w", err)
		}
		newMsg, err := refs.ParseMessageRef(v)
		if err != nil {
			return err
		}
		log.Log("event", "published", "type", "about", "ref", newMsg.String())
		fmt.Fprintln(os.Stdout, newMsg.String())
		return nil
	},
}

var publishContactCmd = &cli.Command{
	Name:      "contact",
	ArgsUsage: "@contactKeypair.ed25519",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "following"},
		&cli.BoolFlag{Name: "blocking"},

		&cli.StringSliceFlag{Name: "recps", Usage: "as a PM to these feeds"},
	},
	Action: func(ctx *cli.Context) error {
		cref, err := refs.ParseFeedRef(ctx.Args().First())
		if err != nil {
			return fmt.Errorf("publish/contact: invalid feed ref: %w", err)
		}
		if ctx.Bool("following") && ctx.Bool("blocking") {
			return fmt.Errorf("publish/contact: can't be both true")
		}
		arg := map[string]interface{}{
			"contact":   cref.String(),
			"type":      "contact",
			"following": ctx.Bool("following"),
			"blocking":  ctx.Bool("blocking"),
		}

		client, err := newClient(ctx)
		if err != nil {
			return err
		}

		var v string
		err = client.Async(longctx, &v, muxrpc.TypeString, muxrpc.Method{"publish"}, arg)
		if err != nil {
			return fmt.Errorf("publish call failed: %w", err)
		}

		newMsg, err := refs.ParseMessageRef(v)
		if err != nil {
			return err
		}
		log.Log("event", "published", "type", "contact", "ref", newMsg.String())
		fmt.Fprintln(os.Stdout, newMsg.String())
		return nil
	},
}
