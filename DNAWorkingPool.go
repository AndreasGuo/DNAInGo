package main

import "GoDNA/DNAAnalysis"

type seqMapSingle struct {
	index int
	seq   DNAAnalysis.Seq
}
type resultMapSingle struct {
	index int
	value float64
}
type seqMapPair struct {
	index1, index2 int
	seq1, seq2     DNAAnalysis.Seq
}
type resultMapPair struct {
	index1, index2 int
	value          float64
}

func continuityWorker(in <-chan seqMapSingle, out chan<- resultMapSingle) {
	for job := range in {
		continuity := DNAAnalysis.Continuity(job.seq, 3)
		out <- resultMapSingle{job.index, float64(continuity)}
	}
}

func hairpinWorker(in <-chan seqMapSingle, out chan<- resultMapSingle) {
	for job := range in {
		hairpin := DNAAnalysis.Hairpin(job.seq, 6, 6, 3)
		out <- resultMapSingle{job.index, float64(hairpin)}
	}
}

func hmeasureWorker(in <-chan seqMapPair, out chan<- resultMapPair) {
	for job := range in {
		hm := DNAAnalysis.HMeasure(job.seq1, job.seq2)
		out <- resultMapPair{job.index1, job.index2, float64(hm)}
	}
}

func similarityWorker(in <-chan seqMapPair, out chan<- resultMapPair) {
	for job := range in {
		sm := DNAAnalysis.Similarity(job.seq1, job.seq2)
		out <- resultMapPair{job.index1, job.index2, float64(sm)}
	}
}

func meltingTemperatureWorker(in <-chan seqMapSingle, out chan<- resultMapSingle) {
	for job := range in {
		mt := DNAAnalysis.MeltingTemperature(job.seq)
		out <- resultMapSingle{job.index, float64(mt)}
	}
}

type FitChan struct {
	ctIn chan seqMapSingle
	ctRe chan resultMapSingle
	hpIn chan seqMapSingle
	hpRe chan resultMapSingle
	hmIn chan seqMapPair
	hmRe chan resultMapPair
	smIn chan seqMapPair
	smRe chan resultMapPair
	mtIn chan seqMapSingle
	mtRe chan resultMapSingle
}

func (fitChan *FitChan) Close() {
	close(fitChan.ctIn)
	close(fitChan.ctRe)
	close(fitChan.hpIn)
	close(fitChan.hpRe)
	close(fitChan.hmIn)
	close(fitChan.hmRe)
	close(fitChan.smIn)
	close(fitChan.smRe)
	close(fitChan.mtIn)
	close(fitChan.mtRe)
}

func CreateWorker(numOfSingle, numOfPair int, bufferSize int) *FitChan {
	continuityChan := make(chan seqMapSingle, bufferSize)
	continuityResult := make(chan resultMapSingle, bufferSize)
	hairpinChan := make(chan seqMapSingle, bufferSize)
	hairpinResult := make(chan resultMapSingle, bufferSize)
	hmChan := make(chan seqMapPair, bufferSize)
	hmResult := make(chan resultMapPair, bufferSize)
	smChan := make(chan seqMapPair, bufferSize)
	smResult := make(chan resultMapPair, bufferSize)
	mtChan := make(chan seqMapSingle, bufferSize)
	mtResult := make(chan resultMapSingle, bufferSize)
	for i := 0; i < numOfSingle; i++ {
		go continuityWorker(continuityChan, continuityResult)
		go hairpinWorker(hairpinChan, hairpinResult)
		go meltingTemperatureWorker(mtChan, mtResult)
	}
	for i := 0; i < numOfPair; i++ {
		go hmeasureWorker(hmChan, hmResult)
		go similarityWorker(smChan, smResult)
	}
	return &FitChan{continuityChan,
		continuityResult,
		hairpinChan,
		hairpinResult,
		hmChan,
		hmResult,
		smChan,
		smResult,
		mtChan,
		mtResult}
}
