package genetics

import (
	"fmt"
	"strings"
	"unicode"
)

func CalculateAdvancedSexLinkedCross(femalePhase, maleGenotype string, f float64) map[string]float64 {
	results := make(map[string]float64)

	femaleRawGemetes := GenerateLinkedGametes(femalePhase, f)
	eggGemetes := make(map[string]float64)

	for g, prob := range femaleRawGemetes {
		eggGemetes[g] = prob
	}

	spermGemetes := map[string]float64{
		maleGenotype: 0.5,
		"Y":          0.5,
	}

	for eCore, eProb := range eggGemetes {
		for sCore, sProb := range spermGemetes {
			var offspring string
			if sCore == "Y" {
				offspring = fmt.Sprintf("X^%s Y", eCore)
			} else {
				if eCore < sCore {
					offspring = fmt.Sprintf("X^%s X^%s", eCore, sCore)
				} else {
					offspring = fmt.Sprintf("X^%s X^%s", sCore, eCore)
				}
			}

			results[offspring] += eProb * sProb
		}
	}

	return results
}

func DecodeAdvancedSexLinkedPhenotype(results map[string]float64) map[string]float64 {
	phenoRatios := make(map[string]float64)

	for geno, prob := range results {
		isMale := strings.Contains(geno, "Y")
		gender := "Female"
		if isMale {
			gender = "Male"
		}

		var phenotypes []string

		if isMale {
			core := strings.Split(geno, " ")[0][2:]
			for _, allele := range core {
				if unicode.IsUpper(allele) {
					phenotypes = append(phenotypes, string(allele)+"-Dom")
				} else {
					phenotypes = append(phenotypes, string(allele)+"-Rec")
				}
			}
		} else {
			parts := strings.Split(geno, " ")
			core1 := parts[0][2:]
			core2 := parts[0][2:]

			for i := 0; i < len(core1); i++ {
				a1 := core1[i]
				a2 := core2[i]

				if unicode.IsUpper(rune(a1)) || unicode.IsUpper(rune(a2)) {
					dom := unicode.ToUpper(rune(a1))
					phenotypes = append(phenotypes, string(dom)+"-Dom")
				} else {
					rec := unicode.ToLower(rune(a1))
					phenotypes = append(phenotypes, string(rec)+"-Rec")
				}
			}
		}

		phenoStr := fmt.Sprintf("(%s) %s", gender, strings.Join(phenotypes, ", "))
		phenoRatios[phenoStr] += prob
	}

	return phenoRatios
}
