package genetics

import (
	"strings"
)

func GenerateLinkedGametes(phase string, f float64) map[string]float64 {
	gametes := make(map[string]float64)

	chromosomes := strings.Split(phase, "/")
	if len(chromosomes) != 2 || len(chromosomes[0]) != 2 || len(chromosomes[1]) != 2 {
		return gametes
	}

	c1 := chromosomes[0]
	c2 := chromosomes[1]

	parentalProb := (1.0 - f) / 2.0
	gametes[c1] += parentalProb
	gametes[c2] += parentalProb

	recombinant1 := string(c1[0]) + string(c2[1])
	recombinant2 := string(c2[0]) + string(c1[1])

	recombinantProb := f / 2.0
	gametes[recombinant1] += recombinantProb
	gametes[recombinant2] += recombinantProb

	return gametes
}

func CombineGemetes(gamete1, gamete2 string) string {
	var builder strings.Builder
	for i := 0; i < len(gamete1); i++ {
		t := Trait{Allele1: Allele(gamete1[i]), Allele2: Allele(gamete2[i])}
		builder.WriteString(t.Normalize())
	}
	return builder.String()
}

func CalculateLinkedCross(p1Phase, p2Phase string, f float64) PunnetResult {
	gametes1 := GenerateLinkedGametes(p1Phase, f)
	gametes2 := GenerateLinkedGametes(p2Phase, f)

	genoTypeRatios := make(map[string]float64)

	for g1, prob1 := range gametes1 {
		for g2, prob2 := range gametes2 {
			combineGeno := CombineGemetes(g1, g2)
			genoTypeRatios[combineGeno] += prob1 * prob2
		}
	}

	phenoTypeRatios := make(map[string]float64)
	for geno, prob := range genoTypeRatios {
		pheno := determinePhenotype(geno)
		phenoTypeRatios[pheno] += prob
	}

	p1Flat := strings.ReplaceAll(p1Phase, "/", "")
	p2Flat := strings.ReplaceAll(p2Phase, "/", "")

	return PunnetResult{
		Parent1:         ParseGenotype(CombineGemetes(p1Flat[0:2], p1Flat[2:4])),
		Parent2:         ParseGenotype(CombineGemetes(p2Flat[0:2], p2Flat[2:4])),
		GenotypeRatios:  genoTypeRatios,
		PhenotypeRatios: phenoTypeRatios,
	}
}
