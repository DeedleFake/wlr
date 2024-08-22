package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log/slog"
	"os"
	"strings"
	"text/template"
)

func gen(w io.Writer) error {
	file, err := os.Open("/usr/include/xkbcommon/xkbcommon-keysyms.h")
	if err != nil {
		return err
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		parts := strings.Fields(s.Text())
		if len(parts) < 3 {
			continue
		}
		if parts[0] != "#define" {
			continue
		}

		name, ok := strings.CutPrefix(parts[1], "XKB_KEY_")
		if !ok {
			continue
		}

		var comment string
		if len(parts) > 3 {
			comment = " " + strings.Join(parts[3:], " ")
		}

		_, err := fmt.Fprintf(w, "\tKeySym%v KeySym = C.%v%v\n", name, parts[1], comment)
		if err != nil {
			return fmt.Errorf("generate %v: %w", name, err)
		}
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

func genString() (string, error) {
	var str strings.Builder
	err := gen(&str)
	return str.String(), err
}

func main() {
	tmpl := template.New("syms.go").Funcs(map[string]any{
		"gen": genString,
	})
	template.Must(tmpl.Parse(`package xkb

/*
#include <xkbcommon/xkbcommon-keysyms.h>
*/
import "C"

type KeySym int32

const (
	{{gen}}
)`))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, nil)
	if err != nil {
		slog.Error("execute template", "err", err)
		os.Exit(1)
	}

	rawoutput := buf.Bytes()
	formatted, err := format.Source(rawoutput)
	if err != nil {
		slog.Error("format output", "err", err)
		formatted = rawoutput
	}

	file, err := os.Create(os.Args[1])
	if err != nil {
		slog.Error("create output file", "file", os.Args[1], "err", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.Write(formatted)
	if err != nil {
		slog.Error("write output", "err", err)
		os.Exit(1)
	}
}
