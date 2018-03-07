package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"

	structurebidon "gitlab.in.weborama.fr/go/pprofPOC/structureBidon"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var traceprofile = flag.String("trace", "", "Trace the program during its lifespan")

var solus = flag.Bool("solus", false, "Execture the program with once enhanced version of the program to show the profiling utility")

func initCPUProfiling() {
	f, err := os.Create(*cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
}

func memoryProfile() {
	f, err := os.OpenFile(*memprofile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal(err)
	}
	runtime.GC() // get up-to-date statistics
	pprof.WriteHeapProfile(f)
	f.Close()
}

func sampleProfilingWrongCPU() {
	collectionOfBidon := map[*structurebidon.Bidon]int{}
	// Create 10 000 object
	// and save one in the middle to get ask for its order apparition later
	var savedBidon *structurebidon.Bidon
	for creationOrder := 0; creationOrder < 10000; creationOrder++ {
		anOtherBidon := structurebidon.New()
		if 5000 == creationOrder {
			savedBidon = anOtherBidon
		}
		collectionOfBidon[anOtherBidon] = creationOrder
	}

	if order, ok := collectionOfBidon[savedBidon]; ok {
		_ = order
	} else {
		fmt.Println("WTF it was not planned in the progam ! Damned Murphy law, I hate this man")
	}
}

func main() {
	flag.Parse()
	runXTimes := 2000

	var run func()
	if *solus {
		run = solutionForCPUExample
	} else {
		run = sampleProfilingWrongCPU
	}

	if *memprofile != "" {
		run()
		fmt.Println("SnapShot Of the Heap")
		memoryProfile()

	} else if *cpuprofile != "" {
		fmt.Println("Initializing CPU profiling")
		initCPUProfiling()
		defer pprof.StopCPUProfile()

		for i := 0; i < runXTimes; i++ {
			run()
		}
	} else if *traceprofile != "" {
		fmt.Println("Begining tracing of program")

		f, err := os.Create(*traceprofile)
		if err != nil {
			panic("Error while creating trace file")
		}
		defer f.Close()

		if err = trace.Start(f); err != nil {
			panic(err)
		}
		defer trace.Stop()

		for i := 0; i < runXTimes; i++ {
			go run()
		}
	}

	fmt.Println("I finished !\n")
}

func solutionForCPUExample() {
	collectionOfBidon := []*structurebidon.Bidon{}
	// Create 10 000 object
	// and save one in the middle to get ask for its order apparition later
	var savedBidon *structurebidon.Bidon
	for creationOrder := 0; creationOrder < 10000; creationOrder++ {
		anOtherBidon := structurebidon.New()
		if 5000 == creationOrder {
			savedBidon = anOtherBidon
		}
		collectionOfBidon = append(collectionOfBidon, anOtherBidon)
	}

	for creationOrder, bidon := range collectionOfBidon {
		if savedBidon == bidon {
			_ = creationOrder
			return
		}
	}

	fmt.Println("WTF it was not planned in the progam ! Damned Murphy law, I hate this man")
}
