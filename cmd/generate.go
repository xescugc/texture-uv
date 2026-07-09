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

type overlayMapPair struct {
	Name        string
	OverlayPath string
	MapPath     string
	SourcePath  string
}

func scanPairs(dir string) ([]overlayMapPair, []string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read directory %q: %w", dir, err)
	}

	maps := make(map[string]string)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, "map.") && strings.HasSuffix(name, ".png") {
			key := strings.TrimSuffix(strings.TrimPrefix(name, "map."), ".png")
			if key != "" {
				maps[key] = filepath.Join(dir, name)
			}
		}
	}

	var pairs []overlayMapPair
	var warnings []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, "overlay.") && strings.HasSuffix(name, ".png") {
			key := strings.TrimSuffix(strings.TrimPrefix(name, "overlay."), ".png")
			if key == "" {
				continue
			}
			if mapPath, ok := maps[key]; ok {
				pairs = append(pairs, overlayMapPair{
					Name:        key,
					OverlayPath: filepath.Join(dir, name),
					MapPath:     mapPath,
					SourcePath:  filepath.Join(dir, "source."+key+".png"),
				})
			} else {
				warnings = append(warnings, key)
			}
		}
	}

	return pairs, warnings, nil
}

var generateCmd = &cli.Command{
	Name:      "generate",
	Usage:     "Scans a directory for overlay/map pairs and generates source files",
	ArgsUsage: "[directory]",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		dir := cmd.Args().Get(0)
		if dir == "" {
			dir = "."
		}

		pairs, warnings, err := scanPairs(dir)
		if err != nil {
			return err
		}

		for _, w := range warnings {
			fmt.Fprintf(os.Stderr, "warning: overlay.%s.png has no matching map.%s.png\n", w, w)
		}

		var errs []string
		for _, p := range pairs {
			overlayInfo, err := os.Stat(p.OverlayPath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", p.Name, err))
				continue
			}
			mapInfo, err := os.Stat(p.MapPath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", p.Name, err))
				continue
			}

			if sourceInfo, err := os.Stat(p.SourcePath); err == nil {
				sourceMod := sourceInfo.ModTime()
				if sourceMod.After(overlayInfo.ModTime()) && sourceMod.After(mapInfo.ModTime()) {
					continue
				}
			}

			fo, err := os.Open(p.OverlayPath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", p.Name, err))
				continue
			}

			fm, err := os.Open(p.MapPath)
			if err != nil {
				fo.Close()
				errs = append(errs, fmt.Sprintf("%s: %v", p.Name, err))
				continue
			}

			oimg, _, err := image.Decode(fo)
			fo.Close()
			if err != nil {
				fm.Close()
				errs = append(errs, fmt.Sprintf("%s: failed to decode overlay: %v", p.Name, err))
				continue
			}

			mimg, _, err := image.Decode(fm)
			fm.Close()
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: failed to decode map: %v", p.Name, err))
				continue
			}

			src, err := uv.NewSource(oimg, mimg)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", p.Name, err))
				continue
			}

			fout, err := os.Create(p.SourcePath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", p.Name, err))
				continue
			}

			if err = png.Encode(fout, src); err != nil {
				fout.Close()
				os.Remove(p.SourcePath)
				errs = append(errs, fmt.Sprintf("%s: %v", p.Name, err))
				continue
			}
			if err = fout.Close(); err != nil {
				os.Remove(p.SourcePath)
				errs = append(errs, fmt.Sprintf("%s: failed to write source: %v", p.Name, err))
				continue
			}

			fmt.Printf("generate: source.%s.png\n", p.Name)
		}

		if len(errs) > 0 {
			return fmt.Errorf("errors generating sources:\n%s", strings.Join(errs, "\n"))
		}

		return nil
	},
}
