package main

import (
	"GoDNA/DNAAlgorithm"
	"GoDNA/DNAAnalysis"
	"GoDNA/algorithm"
	"fmt"
	"os"
	"strconv"
)

func main() {
	fitChan := DNAAnalysis.CreateWorker(100, 100, 10)
	defer fitChan.Close()

	var config = DNAAlgorithm.DefaultConfig()
	var dnaSet = RandomDNASet(config.DIM, config.DNASIZE)
	result := ""
	for it := 0; it < config.DNASETITERATION; it++ {
		fmt.Println("DNASet iteration ", it+1, " ...")

		for index := range dnaSet {
			fitFunc := DNAAlgorithm.FitnessCall(dnaSet, index, fitChan, config)
			alg := algorithm.PO{Pop: nil, MaxIteration: config.MAXIT}
			pop := new(DNAAlgorithm.DNAPopulation)
			pop.SetConfig(config)
			pop.SetFitFunc(fitFunc)
			alg.Initialize(pop, dnaSet[index])
			inv := alg.Iteration()
			dnaSet[index] = inv
		}
		fmt.Println("\rDone")
		result = ""
		for ind, inv := range dnaSet {
			DNAString, err := inv.String()
			if err != nil {
				panic("error while decoding")
			}
			result += DNAString
			//fmt.Print(DNAString, " ")
			// continuity
			fitChan.CtIn <- DNAAnalysis.SeqMapSingle{ind, inv.Represent()}
			result += " " + strconv.FormatFloat((<-fitChan.CtRe).Value, 'f', 4, 64)
			//fmt.Print((<-fitChan.CtRe).Value, " ")
			// hairpin
			fitChan.HpIn <- DNAAnalysis.SeqMapSingle{ind, inv.Represent()}
			result += " " + strconv.FormatFloat((<-fitChan.HpRe).Value, 'f', 4, 64)
			//fmt.Print((<-fitChan.HpRe).Value, " ")
			// hm
			hmList := make([]float64, len(dnaSet))
			for j, o := range dnaSet {
				fitChan.HmIn <- DNAAnalysis.SeqMapPair{ind, j, inv.Represent(), o.Represent()}
				hmList[j] = (<-fitChan.HmRe).Value
			}
			//fmt.Print(sum(hmList), " ")
			result += " " + strconv.FormatFloat(sum(hmList), 'f', 4, 64)
			// sm
			smList := make([]float64, len(dnaSet))
			for j, o := range dnaSet {
				if j != ind {
					fitChan.SmIn <- DNAAnalysis.SeqMapPair{ind, j, inv.Represent(), o.Represent()}
					smList[j] = (<-fitChan.SmRe).Value
				}
			}
			//fmt.Print(sum(smList), " ")
			result += " " + strconv.FormatFloat(sum(smList), 'f', 4, 64)
			//mt
			fitChan.MtIn <- DNAAnalysis.SeqMapSingle{ind, inv.Represent()}
			//fmt.Printf("%.4f\n", (<-fitChan.MtRe).Value)
			result += " " + strconv.FormatFloat((<-fitChan.MtRe).Value, 'f', 4, 64) + "\n"
		}
		fmt.Println()
		fmt.Println(result)
	}
	fileName := "result.txt"
	os.WriteFile(fileName, []byte(result), 0644)
}

func RandomDNASet(dim, size int) []algorithm.Individual {
	dnaSet := []algorithm.Individual{}
	for i := 0; i < size; i++ {
		dnaSet = append(dnaSet, DNAAlgorithm.CreateDNAAgent(dim, 0, 3))
	}
	return dnaSet
}

func sum[T int | float64](lt []T) T {
	var s T = 0
	for i := range lt {
		s += lt[i]
	}
	return s
}
