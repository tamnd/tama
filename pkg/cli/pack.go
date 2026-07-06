package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/klauspost/compress/zstd"
	"github.com/spf13/cobra"
)

// packHeader is the top of every pack file; the full pack format is an M8
// concern, M1 only reads these fields.
type packHeader struct {
	Format  int    `json:"format"`
	Course  string `json:"course"`
	Version int    `json:"version"`
	Title   string `json:"title"`
}

func newPackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pack",
		Short: "Inspect course pack files",
		Long: "Reads pack files, zstd-compressed or plain JSON. M1 parses only the\n" +
			"header (format, course, version); deep validation lands with the pack\n" +
			"format in M8.",
		Example: "  tama pack inspect es-from-en-v1.pack",
	}
	cmd.AddCommand(newPackInspectCmd(), newPackValidateCmd())
	return cmd
}

func newPackInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "inspect <file>",
		Short:   "Print a pack's header",
		Long:    "Prints the header fields of a pack file.",
		Example: "  tama pack inspect es-from-en-v1.pack",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := readPackHeader(args[0])
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "format:  %d\n", h.Format)
			fmt.Fprintf(out, "course:  %s\n", h.Course)
			fmt.Fprintf(out, "version: %d\n", h.Version)
			if h.Title != "" {
				fmt.Fprintf(out, "title:   %s\n", h.Title)
			}
			return nil
		},
	}
}

func newPackValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "validate <file>",
		Short:   "Check a pack's header",
		Long:    "Fails unless the file parses and the header carries a format, course, and version.",
		Example: "  tama pack validate es-from-en-v1.pack",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := readPackHeader(args[0])
			if err != nil {
				return err
			}
			if h.Format < 1 {
				return fmt.Errorf("%s: missing or invalid format", args[0])
			}
			if h.Course == "" {
				return fmt.Errorf("%s: missing course", args[0])
			}
			if h.Version < 1 {
				return fmt.Errorf("%s: missing or invalid version", args[0])
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s: header ok (format %d, %s v%d)\n", args[0], h.Format, h.Course, h.Version)
			return nil
		},
	}
}

// zstdMagic starts every zstd frame.
var zstdMagic = []byte{0x28, 0xb5, 0x2f, 0xfd}

func readPackHeader(path string) (packHeader, error) {
	var h packHeader
	raw, err := os.ReadFile(path)
	if err != nil {
		return h, err
	}
	if bytes.HasPrefix(raw, zstdMagic) {
		dec, err := zstd.NewReader(nil)
		if err != nil {
			return h, err
		}
		defer dec.Close()
		if raw, err = dec.DecodeAll(raw, nil); err != nil {
			return h, fmt.Errorf("%s: decompress: %w", path, err)
		}
	}
	if err := json.Unmarshal(raw, &h); err != nil {
		return h, fmt.Errorf("%s: parse header: %w", path, err)
	}
	return h, nil
}
