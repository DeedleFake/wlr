package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
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
		const prefix = "XKB_KEY_"

		parts := strings.Fields(s.Text())
		if len(parts) < 2 {
			continue
		}

		name := strings.TrimPrefix(parts[1], prefix)
		if name == parts[1] {
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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: format: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Create(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: create output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.Write(formatted)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: write output file: %v\n", err)
		os.Exit(1)
	}
}
