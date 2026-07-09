package cmd

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/xescugc/texture-uv/uv"
)

var batchApplyCmd = &cli.Command{
	Name:      "batch-apply",
	Usage:     "Applies multiple lookups to the same source",
	ArgsUsage: "[source, lookups-dir]",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Required: true, Usage: "Output directory"},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		source := cmd.Args().Get(0)
		if source == "" {
			return fmt.Errorf("Source is required")
		}

		lookupsDir := cmd.Args().Get(1)
		if lookupsDir == "" {
			return fmt.Errorf("Lookups directory is required")
		}

		fs, err := os.Open(source)
		if err != nil {
			return fmt.Errorf("failed to open Source at %q: %w", source, err)
		}
		defer fs.Close()

		sourceImg, _, err := image.Decode(fs)
		if err != nil {
			return fmt.Errorf("failed to decode Source: %w", err)
		}

		entries, err := os.ReadDir(lookupsDir)
		if err != nil {
			return fmt.Errorf("failed to read lookups directory %q: %w", lookupsDir, err)
		}

		outputDir := cmd.String("output")
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %q: %w", outputDir, err)
		}

		var errs []string
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if !strings.HasSuffix(strings.ToLower(name), ".png") {
				continue
			}

			lookupPath := filepath.Join(lookupsDir, name)
			fl, err := os.Open(lookupPath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", name, err))
				continue
			}

			lookupImg, _, err := image.Decode(fl)
			fl.Close()
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: failed to decode lookup: %v", name, err))
				continue
			}

			result := uv.Apply(sourceImg, lookupImg)

			outPath := filepath.Join(outputDir, name)
			fout, err := os.Create(outPath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", name, err))
				continue
			}

			if err = png.Encode(fout, result); err != nil {
				fout.Close()
				os.Remove(outPath)
				errs = append(errs, fmt.Sprintf("%s: %v", name, err))
				continue
			}
			if err = fout.Close(); err != nil {
				os.Remove(outPath)
				errs = append(errs, fmt.Sprintf("%s: failed to write output: %v", name, err))
				continue
			}

			fmt.Printf("batch-apply: %s\n", name)
		}

		if len(errs) > 0 {
			return fmt.Errorf("errors applying lookups:\n%s", strings.Join(errs, "\n"))
		}

		return nil
	},
}
