package genetics

import (
	"strings"
	"unicode"
)

func ParseGenotype(s string) Genotype {
	var traits []Trait

	for i := 0; i < len(s); i += 2 {
		t := Trait{
			Allele1: Allele(s[i]),
			Allele2: Allele(s[i+1]),
		}
		traits = append(traits, t)
	}

	return Genotype{Traits: traits}
}

func CrossSingleTrait(t1, t2 Trait) map[string]float64 {
	results := make(map[string]float64)
	allele1 := []Allele{t1.Allele1, t1.Allele2}
	allele2 := []Allele{t2.Allele1, t2.Allele2}

	for _, a1 := range allele1 {
		for _, a2 := range allele2 {
			t := Trait{Allele1: a1, Allele2: a2}
			results[t.Normalize()] += 0.25
		}
	}

	return results
}

func combineTraits(traits []map[string]float64) map[string]float64 {
	if len(traits) == 1 {
		return traits[0]
	}

	currentTrait := traits[0]
	ramainingTraits := combineTraits(traits[1:])
	combined := make(map[string]float64)

	for g1, p1 := range currentTrait {
		for g2, p2 := range ramainingTraits {
			combined[g1+g2] += p1 * p2
		}
	}

	return combined
}

func determinePhenotype(genotypeStr string) string {
	var phenotypes []string
	for i := 0; i < len(genotypeStr); i += 2 {
		a1 := Allele(genotypeStr[i])
		a2 := Allele(genotypeStr[i+1])

		if a1.IsDominant() || a2.IsDominant() {
			dom := unicode.ToUpper(rune(a1))
			phenotypes = append(phenotypes, string(dom)+"-Dom")
		} else {
			rec := unicode.ToLower(rune(a1))
			phenotypes = append(phenotypes, string(rec)+"-Rec")
		}
	}

	return strings.Join(phenotypes, ", ")
}

func CalculatePolyHybrid(p1Str, p2Str string) PunnetResult {
	p1 := ParseGenotype(p1Str)
	p2 := ParseGenotype(p2Str)

	var singleCrosses []map[string]float64

	for i := 0; i < len(p1.Traits); i++ {
		crossResult := CrossSingleTrait(p1.Traits[i], p2.Traits[i])
		singleCrosses = append(singleCrosses, crossResult)
	}

	genotypeRatios := combineTraits(singleCrosses)
	phenotypeRatios := make(map[string]float64)

	for geno, prob := range genotypeRatios {
		phenotype := determinePhenotype(geno)
		phenotypeRatios[phenotype] += prob
	}

	return PunnetResult{
		Parent1:         p1,
		Parent2:         p2,
		GenotypeRatios:  genotypeRatios,
		PhenotypeRatios: phenotypeRatios,
	}
}
