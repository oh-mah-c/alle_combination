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

func CalculateNondisjunctionCross(femalePhase, maleGenotype string, femaleErr, maleErr bool) map[string]float64 {
	results := make(map[string]float64)

	var eggs []string
	if femaleErr {
		cores := strings.Split(femalePhase, "/")
		if len(cores) == 2 {
			eggs = []string{fmt.Sprintf("X^%s X^%s", cores[0], cores[1]), "0"}
		}
	} else {
		cores := strings.Split(femalePhase, "/")
		if len(cores) == 2 {
			eggs = []string{fmt.Sprintf("X^%s", cores[0]), fmt.Sprintf("X^%s", cores[1])}
		}
	}

	var sperms []string
	if maleErr {
		sperms = []string{fmt.Sprintf("X^%s Y", maleGenotype), "0"}
	} else {
		sperms = []string{fmt.Sprintf("X^%s", maleGenotype), "Y"}
	}

	for _, e := range eggs {
		for _, s := range sperms {
			var parts []string
			if e != "0" {
				parts = append(parts, strings.Split(e, " ")...)
			}

			if s != "0" {
				parts = append(parts, strings.Split(s, " ")...)
			}

			if len(parts) == 0 {
				results["0"] += 0.25
				continue
			}

			offspring := strings.Join(parts, " ")
			results[offspring] += 0.25
		}
	}

	return results
}

func DecodeAdvancedSexLinkedPhenotype(results map[string]float64) map[string]float64 {
	phenoRatios := make(map[string]float64)

	for geno, prob := range results {
		if geno == "0" {
			phenoRatios["(Theory) 0 - Do not have sex chromosomes"] += prob
			continue
		}

		isMale := strings.Contains(geno, "Y")
		xCount := strings.Count(geno, "X^")

		gender := "Female"
		if isMale {
			gender = "Male"
		}

		if xCount == 0 && isMale {
			phenoRatios["(Theory) Y0 - Missing X"] += prob
			continue
		} else if xCount == 1 && !isMale {
			gender = "Female (Tunner X-0)"
		} else if xCount == 2 && isMale {
			gender = "Male (Klinefelter XXY)"
		} else if xCount == 3 && !isMale {
			gender = "Female (Triple X)"
		}

		var xCores []string
		parts := strings.Split(geno, " ")
		for _, p := range parts {
			if strings.HasPrefix(p, "X^") {
				xCores = append(xCores, p[2:])
			}
		}

		var phenotypes []string
		if len(xCores) > 0 {
			corelen := len(xCores[0])

			for i := 0; i < corelen; i++ {
				isDominant := false
				var locusChar rune

				for _, core := range xCores {
					allele := rune(core[i])
					locusChar = allele
					if unicode.IsUpper(allele) {
						isDominant = true
					}
				}

				if isDominant {
					phenotypes = append(phenotypes, fmt.Sprintf("%c-Dom", unicode.ToUpper(locusChar)))
				} else {
					phenotypes = append(phenotypes, fmt.Sprintf("%c-Rec", unicode.ToLower(locusChar)))
				}
			}
		}

		phenoStr := fmt.Sprintf("(%s) %s", gender, strings.Join(phenotypes, ", "))
		phenoRatios[phenoStr] += prob

	}

	return phenoRatios
}
