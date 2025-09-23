package cmd

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/xescugc/golang-texture-uv/uv"
)

var (
	Cmd = &cli.Command{
		Name:  "Texture UV",
		Usage: "Tool to convert assets to Texture UV maps and apply a lookup to that UV",
		//Flags: []cli.Flag{
		//&cli.BoolFlag{Name: "verbose", Value: false, Usage: "Activate verbose mode to display logs and info"},
		//},
		Commands: []*cli.Command{
			{
				Name:      "new-source",
				Usage:     "Creates a new Source based on the Overlay and the Map",
				ArgsUsage: "[overlay, map]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Source output", Required: true},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					overlay := cmd.Args().Get(0)
					if overlay == "" {
						return fmt.Errorf("Overlay is required")
					}

					omap := cmd.Args().Get(1)
					if omap == "" {
						return fmt.Errorf("Map is required")
					}

					fo, err := os.Open(overlay)
					if err != nil {
						return fmt.Errorf("failed to open Overlay at %q: %w", overlay, err)
					}
					defer fo.Close()

					fm, err := os.Open(omap)
					if err != nil {
						return fmt.Errorf("failed to open Map at %q: %w", omap, err)
					}
					defer fm.Close()

					oimg, _, err := image.Decode(fo)
					if err != nil {
						return fmt.Errorf("failed to decode Overlay: %w", err)
					}
					mimg, _, err := image.Decode(fm)
					if err != nil {
						return fmt.Errorf("failed to decode Map: %w", err)
					}

					src, err := uv.NewSource(oimg, mimg)
					if err != nil {
						return fmt.Errorf("failed to create Source: %w", err)
					}

					fout, err := os.Create(cmd.String("output"))
					if err != nil {
						return fmt.Errorf("failed to create output file at %q: %w", cmd.String("output"), err)
					}
					defer fout.Close()

					if err = png.Encode(fout, src); err != nil {
						return fmt.Errorf("failed to encode PNG: %w", err)
					}

					return nil
				},
			},
			{
				Name:      "apply",
				Usage:     "Applies to the source the lookup",
				ArgsUsage: "[source, lookup]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Image output", Required: true},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					source := cmd.Args().Get(0)
					if source == "" {
						return fmt.Errorf("Source is required")
					}

					lookup := cmd.Args().Get(1)
					if lookup == "" {
						return fmt.Errorf("Lookup is required")
					}

					fs, err := os.Open(source)
					if err != nil {
						return fmt.Errorf("failed to open Source at %q: %w", source, err)
					}
					defer fs.Close()

					fl, err := os.Open(lookup)
					if err != nil {
						return fmt.Errorf("failed to open Lookup at %q: %w", lookup, err)
					}
					defer fl.Close()

					simg, _, err := image.Decode(fs)
					if err != nil {
						return fmt.Errorf("failed to decode Source: %w", err)
					}
					limg, _, err := image.Decode(fl)
					if err != nil {
						return fmt.Errorf("failed to decode Lookup: %w", err)
					}

					img := uv.Apply(simg, limg)

					fout, err := os.Create(cmd.String("output"))
					if err != nil {
						return fmt.Errorf("failed to create output file at %q: %w", cmd.String("output"), err)
					}
					defer fout.Close()

					if err = png.Encode(fout, img); err != nil {
						return fmt.Errorf("failed to encode PNG: %w", err)
					}

					return nil
				},
			},
		},
	}
)
