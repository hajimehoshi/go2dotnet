// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/sync/errgroup"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func writeInstCS(dir string, namespace string, importFuncs, funcs []*Func, exports []*Export, globals []*Global, types []*Type, tables [][]uint32) error {
	const groupSize = 64

	var g errgroup.Group
	for i := 0; i < (len(funcs)-1)/groupSize+1; i++ {
		i := i
		fs := funcs[groupSize*i : min(groupSize*(i+1), len(funcs))]
		g.Go(func() error {
			f, err := os.Create(filepath.Join(dir, fmt.Sprintf("Inst.Funcs%d.cs", i)))
			if err != nil {
				return err
			}
			defer f.Close()

			if err := instTmpl.Execute(f, struct {
				Namespace string
				Funcs     []*Func
			}{
				Namespace: namespace,
				Funcs:     fs,
			}); err != nil {
				return err
			}
			return nil
		})
	}
	g.Go(func() error {
		f, err := os.Create(filepath.Join(dir, "Inst.Exports.cs"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := instExportsTmpl.Execute(f, struct {
			Namespace string
			Exports   []*Export
		}{
			Namespace: namespace,
			Exports:   exports,
		}); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		f, err := os.Create(filepath.Join(dir, "Inst.Init.cs"))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := instInitTmpl.Execute(f, struct {
			Namespace   string
			ImportFuncs []*Func
			Funcs       []*Func
			Types       []*Type
			Tables      [][]uint32
			Globals     []*Global
		}{
			Namespace:   namespace,
			ImportFuncs: importFuncs,
			Funcs:       funcs,
			Types:       types,
			Tables:      tables,
			Globals:     globals,
		}); err != nil {
			return err
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

var instTmpl = template.Must(template.New("Inst.cs").Parse(`// Code generated by go2dotnet. DO NOT EDIT.

#pragma warning disable 162 // unreachable code
#pragma warning disable 164 // label
#pragma warning disable 219 // unused local variables

using System;
using System.Diagnostics;

namespace {{.Namespace}}
{
    sealed partial class Inst
    {
{{range $value := .Funcs}}{{$value.CSharp "        " false true}}
{{end}}    }
}
`))

var instExportsTmpl = template.Must(template.New("Inst.Exports.cs").Parse(`// Code generated by go2dotnet. DO NOT EDIT.

namespace {{.Namespace}}
{
    sealed partial class Inst
    {
{{range $value := .Exports}}{{$value.CSharp "        "}}
{{end}}    }
}
`))

var instInitTmpl = template.Must(template.New("Inst.Init.cs").Parse(`// Code generated by go2dotnet. DO NOT EDIT.

namespace {{.Namespace}}
{
    sealed partial class Inst
    {
        public Inst(Mem mem, IImport import)
        {
             initializeFuncs_();
             mem_ = mem;
             import_ = import;
        }

        private void initializeFuncs_()
        {
            funcs_ = new object[] {
{{range $value := .ImportFuncs}}                null,
{{end}}{{range $value := .Funcs}}                (Type{{.Type.Index}})({{.Identifier}}),
{{end}}            };
        }

{{range $value := .Types}}{{$value.CSharp "        "}}
{{end}}
        private static readonly uint[][] table_ = {
{{range $value := .Tables}}            new uint[] { {{- range $value2 := $value}}{{$value2}}, {{end}}},
{{end}}        };

{{range $value := .Globals}}{{$value.CSharp "        "}}
{{end}}
        private object[] funcs_;
        private Mem mem_;
        private IImport import_;
    }
}
`))
