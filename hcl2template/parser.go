package hcl2template

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/dynblock"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/packer/packer"
	"github.com/zclconf/go-cty/cty"
)

const (
	packerLabel       = "packer"
	sourceLabel       = "source"
	variablesLabel    = "variables"
	variableLabel     = "variable"
	localsLabel       = "locals"
	buildLabel        = "build"
	communicatorLabel = "communicator"
)

var configSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: packerLabel},
		{Type: sourceLabel, LabelNames: []string{"type", "name"}},
		{Type: variablesLabel},
		{Type: variableLabel, LabelNames: []string{"name"}},
		{Type: localsLabel},
		{Type: buildLabel},
		{Type: communicatorLabel, LabelNames: []string{"type", "name"}},
	},
}

// packerBlockSchema is the schema for a top-level "packer" block in
// a configuration file.
var packerBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "required_version"},
	},
}

// Parser helps you parse HCL folders. It will parse an hcl file or directory
// and start builders, provisioners and post-processors to configure them with
// the parsed HCL and then return a []packer.Build. Packer will use that list
// of Builds to run everything in order.
type Parser struct {
	*hclparse.Parser

	BuilderSchemas packer.BuilderStore

	ProvisionersSchemas packer.ProvisionerStore

	PostProcessorsSchemas packer.PostProcessorStore
}

const (
	hcl2FileExt        = ".pkr.hcl"
	hcl2JsonFileExt    = ".pkr.json"
	hcl2VarFileExt     = ".auto.pkrvars.hcl"
	hcl2VarJsonFileExt = ".auto.pkrvars.json"
)

// Parse will Parse all HCL files in filename. Path can be a folder or a file.
//
// Parse will first Parse variables and then the rest; so that interpolation
// can happen.
//
// Parse returns a PackerConfig that contains configuration layout of a packer
// build; sources(builders)/provisioners/posts-processors will not be started
// and their contents wont be verified; Most syntax errors will cause an error.
func (p *Parser) Parse(filename string, varFiles []string, argVars map[string]string) (*PackerConfig, hcl.Diagnostics) {
	var files []*hcl.File
	var diags hcl.Diagnostics

	// parse config files
	if filename != "" {
		hclFiles, jsonFiles, moreDiags := GetHCL2Files(filename, hcl2FileExt, hcl2JsonFileExt)
		diags = append(diags, moreDiags...)
		if len(hclFiles)+len(jsonFiles) == 0 {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Could not find any config file in " + filename,
				Detail: "A config file must be suffixed with `.pkr.hcl` or " +
					"`.pkr.json`. A folder can be referenced.",
			})
		}
		for _, filename := range hclFiles {
			f, moreDiags := p.ParseHCLFile(filename)
			diags = append(diags, moreDiags...)
			files = append(files, f)
		}
		for _, filename := range jsonFiles {
			f, moreDiags := p.ParseJSONFile(filename)
			diags = append(diags, moreDiags...)
			files = append(files, f)
		}
		if diags.HasErrors() {
			return nil, diags
		}
	}

	basedir := filename
	if isDir, err := isDir(basedir); err == nil && !isDir {
		basedir = filepath.Dir(basedir)
	}
	wd, err := os.Getwd()
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Could not find current working directory",
			Detail:   err.Error(),
		})
	}
	cfg := &PackerConfig{
		Basedir:               basedir,
		Cwd:                   wd,
		builderSchemas:        p.BuilderSchemas,
		provisionersSchemas:   p.ProvisionersSchemas,
		postProcessorsSchemas: p.PostProcessorsSchemas,
		parser:                p,
		files:                 files,
	}

	for _, file := range files {
		coreVersionConstraints, moreDiags := sniffCoreVersionRequirements(file.Body)
		cfg.Packer.VersionConstraints = append(cfg.Packer.VersionConstraints, coreVersionConstraints...)
		diags = append(diags, moreDiags...)
	}

	// Before we go further, we'll check to make sure this version can read
	// that file, so we can produce a version-related error message rather than
	// potentially-confusing downstream errors.
	versionDiags := cfg.CheckCoreVersionRequirements()
	diags = append(diags, versionDiags...)
	if versionDiags.HasErrors() {
		return cfg, diags
	}

	// Decode variable blocks so that they are available later on. Here locals
	// can use input variables so we decode them firsthand.
	{
		for _, file := range files {
			diags = append(diags, cfg.decodeInputVariables(file)...)
		}

		for _, file := range files {
			moreLocals, morediags := cfg.parseLocalVariables(file)
			diags = append(diags, morediags...)
			cfg.LocalBlocks = append(cfg.LocalBlocks, moreLocals...)
		}
	}

	// parse var files
	{
		hclVarFiles, jsonVarFiles, moreDiags := GetHCL2Files(filename, hcl2VarFileExt, hcl2VarJsonFileExt)
		diags = append(diags, moreDiags...)
		for _, file := range varFiles {
			switch filepath.Ext(file) {
			case ".hcl":
				hclVarFiles = append(hclVarFiles, file)
			case ".json":
				jsonVarFiles = append(jsonVarFiles, file)
			default:
				diags = append(moreDiags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Could not guess format of " + file,
					Detail:   "A var file must be suffixed with `.hcl` or `.json`.",
				})
			}
		}
		var varFiles []*hcl.File
		for _, filename := range hclVarFiles {
			f, moreDiags := p.ParseHCLFile(filename)
			diags = append(diags, moreDiags...)
			if moreDiags.HasErrors() {
				continue
			}
			varFiles = append(varFiles, f)
		}
		for _, filename := range jsonVarFiles {
			f, moreDiags := p.ParseJSONFile(filename)
			diags = append(diags, moreDiags...)
			if moreDiags.HasErrors() {
				continue
			}
			varFiles = append(varFiles, f)
		}

		diags = append(diags, cfg.collectInputVariableValues(os.Environ(), varFiles, argVars)...)
	}
	return cfg, diags
}

// sniffCoreVersionRequirements does minimal parsing of the given body for
// "packer" blocks with "required_version" attributes, returning the
// requirements found.
//
// This is intended to maximize the chance that we'll be able to read the
// requirements (syntax errors notwithstanding) even if the config file contains
// constructs that might've been added in future versions
//
// This is a "best effort" sort of method which will return constraints it is
// able to find, but may return no constraints at all if the given body is
// so invalid that it cannot be decoded at all.
func sniffCoreVersionRequirements(body hcl.Body) ([]VersionConstraint, hcl.Diagnostics) {

	var sniffRootSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: packerLabel,
			},
		},
	}

	rootContent, _, diags := body.PartialContent(sniffRootSchema)

	var constraints []VersionConstraint

	for _, block := range rootContent.Blocks {
		content, blockDiags := block.Body.Content(packerBlockSchema)
		diags = append(diags, blockDiags...)

		attr, exists := content.Attributes["required_version"]
		if !exists {
			continue
		}

		constraint, constraintDiags := decodeVersionConstraint(attr)
		diags = append(diags, constraintDiags...)
		if !constraintDiags.HasErrors() {
			constraints = append(constraints, constraint)
		}
	}

	return constraints, diags
}

func (cfg *PackerConfig) Initialize() hcl.Diagnostics {
	var diags hcl.Diagnostics

	_, moreDiags := cfg.InputVariables.Values()
	diags = append(diags, moreDiags...)
	_, moreDiags = cfg.LocalVariables.Values()
	diags = append(diags, moreDiags...)
	diags = append(diags, cfg.evaluateLocalVariables(cfg.LocalBlocks)...)

	for _, variable := range cfg.InputVariables {
		if !variable.Sensitive {
			continue
		}
		value, _ := variable.Value()
		_ = cty.Walk(value, func(_ cty.Path, nested cty.Value) (bool, error) {
			if nested.IsWhollyKnown() && !nested.IsNull() && nested.Type().Equals(cty.String) {
				packer.LogSecretFilter.Set(nested.AsString())
			}
			return true, nil
		})
	}

	// decode the actual content
	for _, file := range cfg.files {
		diags = append(diags, cfg.parser.decodeConfig(file, cfg)...)
	}

	return diags
}

// decodeConfig looks in the found blocks for everything that is not a variable
// block. It should be called after parsing input variables and locals so that
// they can be referenced.
func (p *Parser) decodeConfig(f *hcl.File, cfg *PackerConfig) hcl.Diagnostics {
	var diags hcl.Diagnostics

	body := dynblock.Expand(f.Body, cfg.EvalContext(nil))
	content, moreDiags := body.Content(configSchema)
	diags = append(diags, moreDiags...)

	for _, block := range content.Blocks {
		switch block.Type {
		case sourceLabel:
			source, moreDiags := p.decodeSource(block)
			diags = append(diags, moreDiags...)
			if moreDiags.HasErrors() {
				continue
			}

			ref := source.Ref()
			if existing, found := cfg.Sources[ref]; found {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate " + sourceLabel + " block",
					Detail: fmt.Sprintf("This "+sourceLabel+" block has the "+
						"same builder type and name as a previous block declared "+
						"at %s. Each "+sourceLabel+" must have a unique name per builder type.",
						existing.block.DefRange.Ptr()),
					Subject: source.block.DefRange.Ptr(),
				})
				continue
			}

			if cfg.Sources == nil {
				cfg.Sources = map[SourceRef]SourceBlock{}
			}
			cfg.Sources[ref] = source

		case buildLabel:
			build, moreDiags := p.decodeBuildConfig(block, cfg)
			diags = append(diags, moreDiags...)
			if moreDiags.HasErrors() {
				continue
			}
			cfg.Builds = append(cfg.Builds, build)

		}
	}

	return diags
}
