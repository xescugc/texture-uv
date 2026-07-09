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

var newMapCmd = &cli.Command{
	Name:  "new-map",
	Usage: "Generates a Map image from a template",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "from-template", Required: true, Usage: "Path to template/mask image"},
		&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Required: true, Usage: "Map output path"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		templatePath := cmd.String("from-template")

		ft, err := os.Open(templatePath)
		if err != nil {
			return fmt.Errorf("failed to open template at %q: %w", templatePath, err)
		}
		defer ft.Close()

		tmplImg, _, err := image.Decode(ft)
		if err != nil {
			return fmt.Errorf("failed to decode template: %w", err)
		}

		mapImg, err := uv.NewMap(tmplImg)
		if err != nil {
			return fmt.Errorf("failed to create map: %w", err)
		}

		fout, err := os.Create(cmd.String("output"))
		if err != nil {
			return fmt.Errorf("failed to create output file at %q: %w", cmd.String("output"), err)
		}
		defer fout.Close()

		if err = png.Encode(fout, mapImg); err != nil {
			return fmt.Errorf("failed to encode PNG: %w", err)
		}

		return nil
	},
}
