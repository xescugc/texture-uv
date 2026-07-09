package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/xescugc/texture-uv/uv"
)

var newOverlayCmd = &cli.Command{
	Name:      "new-overlay",
	Usage:     "Generates an overlay from a sprite, color mapping, and map image",
	ArgsUsage: "[sprite, color-mapping, map]",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Overlay output", Required: true},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		spritePath := cmd.Args().Get(0)
		if spritePath == "" {
			return fmt.Errorf("Sprite is required")
		}

		mappingPath := cmd.Args().Get(1)
		if mappingPath == "" {
			return fmt.Errorf("Color mapping is required")
		}

		mapPath := cmd.Args().Get(2)
		if mapPath == "" {
			return fmt.Errorf("Map is required")
		}

		fs, err := os.Open(spritePath)
		if err != nil {
			return fmt.Errorf("failed to open Sprite at %q: %w", spritePath, err)
		}
		defer fs.Close()

		fm, err := os.Open(mapPath)
		if err != nil {
			return fmt.Errorf("failed to open Map at %q: %w", mapPath, err)
		}
		defer fm.Close()

		fc, err := os.Open(mappingPath)
		if err != nil {
			return fmt.Errorf("failed to open color mapping at %q: %w", mappingPath, err)
		}
		defer fc.Close()

		spriteImg, _, err := image.Decode(fs)
		if err != nil {
			return fmt.Errorf("failed to decode Sprite: %w", err)
		}

		mapImg, _, err := image.Decode(fm)
		if err != nil {
			return fmt.Errorf("failed to decode Map: %w", err)
		}

		var mapping uv.ColorMapping
		if err := json.NewDecoder(fc).Decode(&mapping); err != nil {
			return fmt.Errorf("failed to decode color mapping: %w", err)
		}

		overlay, err := uv.NewOverlay(spriteImg, mapImg, mapping)
		if err != nil {
			return fmt.Errorf("failed to create overlay: %w", err)
		}

		fout, err := os.Create(cmd.String("output"))
		if err != nil {
			return fmt.Errorf("failed to create output file at %q: %w", cmd.String("output"), err)
		}
		defer fout.Close()

		if err = png.Encode(fout, overlay); err != nil {
			return fmt.Errorf("failed to encode PNG: %w", err)
		}

		return nil
	},
}
