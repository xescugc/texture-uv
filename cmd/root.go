package cmd

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/xescugc/texture-uv/uv"
)

var (
	Cmd = &cli.Command{
		Name:  "Texture UV",
		Usage: "Tool to convert assets to Texture UV maps and apply a lookup to that UV",
		//Flags: []cli.Flag{
		//&cli.BoolFlag{Name: "verbose", Value: false, Usage: "Activate verbose mode to display logs and info"},
		//},
		Commands: []*cli.Command{
			generateCmd,
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
			{
				Name:      "validate",
				Usage:     "Validates that all overlay colors exist in the map",
				ArgsUsage: "[overlay, map]",
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

					mismatches, err := uv.Validate(oimg, mimg)
					if err != nil {
						return fmt.Errorf("validation error: %w", err)
					}

					if len(mismatches) == 0 {
						fmt.Println("OK: all overlay colors found in map")
						return nil
					}

					for _, m := range mismatches {
						fmt.Printf("mismatch at (%d, %d): RGBA(%d, %d, %d, %d)\n", m.X, m.Y, m.R, m.G, m.B, m.A)
					}
					return fmt.Errorf("found %d color mismatch(es)", len(mismatches))
				},
			},
			{
				Name:      "diff",
				Usage:     "Compares two overlays pixel-by-pixel and highlights differences in magenta",
				ArgsUsage: "[overlay1, overlay2]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Diff output", Required: true},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					overlay1 := cmd.Args().Get(0)
					if overlay1 == "" {
						return fmt.Errorf("Overlay1 is required")
					}

					overlay2 := cmd.Args().Get(1)
					if overlay2 == "" {
						return fmt.Errorf("Overlay2 is required")
					}

					f1, err := os.Open(overlay1)
					if err != nil {
						return fmt.Errorf("failed to open Overlay1 at %q: %w", overlay1, err)
					}
					defer f1.Close()

					f2, err := os.Open(overlay2)
					if err != nil {
						return fmt.Errorf("failed to open Overlay2 at %q: %w", overlay2, err)
					}
					defer f2.Close()

					img1, _, err := image.Decode(f1)
					if err != nil {
						return fmt.Errorf("failed to decode Overlay1: %w", err)
					}
					img2, _, err := image.Decode(f2)
					if err != nil {
						return fmt.Errorf("failed to decode Overlay2: %w", err)
					}

					diff := uv.Diff(img1, img2)

					fout, err := os.Create(cmd.String("output"))
					if err != nil {
						return fmt.Errorf("failed to create output file at %q: %w", cmd.String("output"), err)
					}
					defer fout.Close()

					if err = png.Encode(fout, diff); err != nil {
						return fmt.Errorf("failed to encode PNG: %w", err)
					}

					return nil
				},
			},
		},
	}
)
