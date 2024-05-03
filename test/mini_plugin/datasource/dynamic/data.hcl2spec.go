// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package dynamic

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	Nesteds []FlatNestedFirst `mapstructure:"extra" required:"false" cty:"extra" hcl:"extra"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"extra": &hcldec.BlockListSpec{TypeName: "extra", Nested: hcldec.ObjectSpec((*FlatNestedFirst)(nil).HCL2Spec())},
	}
	return s
}

// FlatDatasourceOutput is an auto-generated flat version of DatasourceOutput.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatDatasourceOutput struct {
	Status *string `mapstructure:"data" cty:"data" hcl:"data"`
}

// FlatMapstructure returns a new FlatDatasourceOutput.
// FlatDatasourceOutput is an auto-generated flat version of DatasourceOutput.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*DatasourceOutput) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatDatasourceOutput)
}

// HCL2Spec returns the hcl spec of a DatasourceOutput.
// This spec is used by HCL to read the fields of DatasourceOutput.
// The decoded values from this spec will then be applied to a FlatDatasourceOutput.
func (*FlatDatasourceOutput) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"data": &hcldec.AttrSpec{Name: "data", Type: cty.String, Required: false},
	}
	return s
}

// FlatNestedFirst is an auto-generated flat version of NestedFirst.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatNestedFirst struct {
	Name    *string            `mapstructure:"name" required:"true" cty:"name" hcl:"name"`
	Nesteds []FlatNestedSecond `mapstructure:"extra" required:"false" cty:"extra" hcl:"extra"`
}

// FlatMapstructure returns a new FlatNestedFirst.
// FlatNestedFirst is an auto-generated flat version of NestedFirst.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*NestedFirst) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatNestedFirst)
}

// HCL2Spec returns the hcl spec of a NestedFirst.
// This spec is used by HCL to read the fields of NestedFirst.
// The decoded values from this spec will then be applied to a FlatNestedFirst.
func (*FlatNestedFirst) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"name":  &hcldec.AttrSpec{Name: "name", Type: cty.String, Required: false},
		"extra": &hcldec.BlockListSpec{TypeName: "extra", Nested: hcldec.ObjectSpec((*FlatNestedSecond)(nil).HCL2Spec())},
	}
	return s
}

// FlatNestedSecond is an auto-generated flat version of NestedSecond.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatNestedSecond struct {
	Name *string `mapstructure:"name" required:"true" cty:"name" hcl:"name"`
}

// FlatMapstructure returns a new FlatNestedSecond.
// FlatNestedSecond is an auto-generated flat version of NestedSecond.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*NestedSecond) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatNestedSecond)
}

// HCL2Spec returns the hcl spec of a NestedSecond.
// This spec is used by HCL to read the fields of NestedSecond.
// The decoded values from this spec will then be applied to a FlatNestedSecond.
func (*FlatNestedSecond) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"name": &hcldec.AttrSpec{Name: "name", Type: cty.String, Required: false},
	}
	return s
}
