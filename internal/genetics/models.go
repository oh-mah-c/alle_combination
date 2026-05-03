package genetics

import (
	"strings"
	"unicode"
)

type Allele rune

func (a Allele) IsDominant() bool {
	return unicode.IsUpper(rune(a))
}

func (a Allele) String() string {
	return string(a)
}

type Trait struct {
	Allele1 Allele `json:"allele_1"`
	Allele2 Allele `json:"allele_2"`
}

func (t Trait) Normalize() string {
	if t.Allele1 < t.Allele2 {
		return t.Allele1.String() + t.Allele2.String()
	}
	return t.Allele2.String() + t.Allele1.String()
}

type Genotype struct {
	Traits []Trait `json:"traits"`
}

func (g Genotype) String() string {
	var builder strings.Builder
	for _, t := range g.Traits {
		builder.WriteString(t.Normalize())
	}
	return builder.String()
}

type PunnetResult struct {
	Parent1         Genotype           `json:"parent_1"`
	Parent2         Genotype           `json:"parent_2"`
	GenotypeRatios  map[string]float64 `json:"genotype_ratios"`
	PhenotypeRatios map[string]float64 `json:"phenotype_ratios"`
}
