// Copyright (C) 2016 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package llvm90

import (
	"path/filepath"
	"strings"

	"android/soong/android"
	"android/soong/genrule"

	"github.com/google/blueprint"
)

func init() {
	android.RegisterModuleType("llvm90_tblgen", llvm90TblgenFactory)
}

var (
	pctx = android.NewPackageContext("android/soong/llvm90")

	llvm90Tblgen = pctx.HostBinToolVariable("llvm90Tblgen", "llvm90-tblgen")

	tblgenRule90 = pctx.StaticRule("tblgenRule90", blueprint.RuleParams{
		Depfile:     "${out}.d",
		Deps:        blueprint.DepsGCC,
		Command:     "${llvm90Tblgen} ${includes} ${generator} -d ${depfile} -o ${out} ${in}",
		CommandDeps: []string{"${llvm90Tblgen}"},
		Description: "LLVM90 TableGen $in => $out",
		Restat:      true,
	}, "includes", "depfile", "generator")
)

type tblgenProperties struct {
	In   string
	Outs []string
}

type tblgen struct {
	android.ModuleBase

	properties tblgenProperties

	exportedHeaderDirs android.Paths
	generatedHeaders   android.Paths
}

var _ genrule.SourceFileGenerator = (*tblgen)(nil)

func (t *tblgen) GenerateAndroidBuildActions(ctx android.ModuleContext) {
	in := android.PathForModuleSrc(ctx, t.properties.In)

	includes := []string{
		"-I " + ctx.ModuleDir(),
		"-I external/llvm90/include",
		"-I external/llvm90/lib/Target",
		"-I " + filepath.Dir(in.String()),
	}

	for _, o := range t.properties.Outs {
		out := android.PathForModuleGen(ctx, o)
		generator := outToGenerator(ctx, o)

		ctx.ModuleBuild(pctx, android.ModuleBuildParams{
			Rule:   tblgenRule90,
			Input:  in,
			Output: out,
			Args: map[string]string{
				"includes":  strings.Join(includes, " "),
				"generator": generator,
			},
		})
		t.generatedHeaders = append(t.generatedHeaders, out)
	}

	t.exportedHeaderDirs = append(t.exportedHeaderDirs, android.PathForModuleGen(ctx, ""))
}

func outToGenerator(ctx android.ModuleContext, out string) string {
	out = filepath.Base(out)
	switch {
	case strings.HasSuffix(out, "GenRegisterInfo.inc"):
		return "-gen-register-info"
	case strings.HasSuffix(out, "GenInstrInfo.inc"):
		return "-gen-instr-info"
	case strings.HasSuffix(out, "GenAsmWriter.inc"):
		return "-gen-asm-writer"
	case strings.HasSuffix(out, "GenAsmWriter1.inc"):
		return "-gen-asm-writer -asmwriternum=1"
	case strings.HasSuffix(out, "GenAsmMatcher.inc"):
		return "-gen-asm-matcher"
	case strings.HasSuffix(out, "GenCodeEmitter.inc"):
		return "-gen-emitter"
	case strings.HasSuffix(out, "GenMCCodeEmitter.inc"):
		return "-gen-emitter"
	case strings.HasSuffix(out, "GenMCPseudoLowering.inc"):
		return "-gen-pseudo-lowering"
	case strings.HasSuffix(out, "GenDAGISel.inc"):
		return "-gen-dag-isel"
	case strings.HasSuffix(out, "GenDisassemblerTables.inc"):
		return "-gen-disassembler"
	case strings.HasSuffix(out, "GenSearchableTables.inc"):
		return "-gen-searchable-tables"
	case strings.HasSuffix(out, "GenSystemOperands.inc"):
		return "-gen-searchable-tables"
	case strings.HasSuffix(out, "GenSystemRegister.inc"):
		return "-gen-searchable-tables"
	case strings.HasSuffix(out, "InstCombineTables.inc"):
		return "-gen-searchable-tables"
	case strings.HasSuffix(out, "GenEDInfo.inc"):
		return "-gen-enhanced-disassembly-info"
	case strings.HasSuffix(out, "GenFastISel.inc"):
		return "-gen-fast-isel"
	case strings.HasSuffix(out, "GenSubtargetInfo.inc"):
		return "-gen-subtarget"
	case strings.HasSuffix(out, "GenCallingConv.inc"):
		return "-gen-callingconv"
	case strings.HasSuffix(out, "GenIntrinsicEnums.inc"):
		return "-gen-tgt-intrinsic-enums"
	case strings.HasSuffix(out, "GenIntrinsicImpl.inc"):
		return "-gen-tgt-intrinsic-impl"
	case strings.HasSuffix(out, "GenDecoderTables.inc"):
		return "-gen-arm-decoder"
	case strings.HasSuffix(out, "Options.inc"):
		return "-gen-opt-parser-defs"
	case strings.HasSuffix(out, "GenDFAPacketizer.inc"):
		return "-gen-dfa-packetizer"
	case strings.HasSuffix(out, "GenRegisterBank.inc"):
		return "-gen-register-bank"
	case strings.HasSuffix(out, "GenGlobalISel.inc"):
		return "-gen-global-isel"
	case strings.HasSuffix(out, "X86GenEVEX2VEXTables.inc"):
		return "-gen-x86-EVEX2VEX-tables"
	case out == "Attributes.inc", out == "AttributesCompatFunc.inc":
		return "-gen-attrs"
	case out == "IntrinsicEnums.inc":
		return "-gen-intrinsic-enums"
	case out == "IntrinsicImpl.inc":
		return "-gen-intrinsic-impl"
	}

	ctx.ModuleErrorf("couldn't map output file %q to a generator", out)
	return ""
}

func (t *tblgen) DepsMutator(ctx android.BottomUpMutatorContext) {
}

func (t *tblgen) GeneratedHeaderDirs() android.Paths {
	return t.exportedHeaderDirs
}

func (t *tblgen) GeneratedSourceFiles() android.Paths {
	return nil
}

func (t *tblgen) GeneratedDeps() android.Paths {
	return t.generatedHeaders
}

func llvm90TblgenFactory() android.Module {
	t := &tblgen{}
	t.AddProperties(&t.properties)
	android.InitAndroidModule(t)
	return t
}
