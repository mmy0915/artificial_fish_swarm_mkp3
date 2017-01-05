package main

import (
	. "artificial_fish_swarm_mkp4/fish"
	"fmt"
	"os"
	"sync"
	"time"
)

var l sync.Mutex

const RUN_TIMES = 1
const MAX_GENERATION = 10000000 //最大进化代数

const SHOW = true

var RunTime = make([]float64, RUN_TIMES)
var FirstTimes = make([]int, RUN_TIMES)
var BestResults = make([]float64, RUN_TIMES)
var Individual = make([]*ArtificialFish, POPSIZE)

var startTime time.Time

func main() {

	startTime = time.Now()

	err := InitGoods("G:/source/git/artificial_fish_swarm_mkp/generated_MKP_instances/online/test.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := 0; i < RUN_TIMES; i++ {

		fmt.Printf("runtimes = %d\n", i)

		Bulletin.FoodConsistence = 0

		for j := 0; j < POPSIZE; j++ {
			Individual[j] = NewArtificialFish()
			Individual[j].UpdateSpecimen()
		}

		worker(i)

		BestResults[i] = Bulletin.FoodConsistence
	}

	end := time.Now()
	fmt.Printf("\n总运行时间：%f\n", end.Sub(startTime).Seconds())

	//========================================================
	fmt.Printf("\n每次找到最优解的时间\n")
	for i := 0; i < RUN_TIMES; i++ {
		fmt.Printf("RUN_TIMES[%d], TIME=%f\n", i, RunTime[i])
	}

	//========================================================
	fmt.Printf("\n每次运行的最优值\n")

	for i := 0; i < RUN_TIMES; i++ {
		fmt.Printf("RUN_TIMES[%d], BEST=%f \n", i, BestResults[i])
	}

	//===============================
	fmt.Printf("\n第一次找到最优值时的次数\n")
	for i := 0; i < RUN_TIMES; i++ {
		fmt.Printf("RUN_TIMES[%d], FirstTime=%d\n", i, FirstTimes[i])
	}
}

func worker(runTime int) {
	for g := 0; g < MAX_GENERATION; g++ {
		for pop := 0; pop < POPSIZE; pop++ {
			//Fish's Choice
			Individual[pop].BehaviorSelection(Individual, pop)
			if Individual[pop].FoodConsistence > Bulletin.FoodConsistence {
				distance := Individual[pop].Distance(&Bulletin)
				fmt.Printf("\n第 %d 条鱼 的 %d 代 和公告板的距离：%d\n", pop, g, distance)
				Bulletin = *Individual[pop]
				end := time.Now()
				RunTime[runTime] = end.Sub(startTime).Seconds()
				FirstTimes[runTime] = g
				if SHOW {
					//Display information
					fmt.Printf("第 %d 条鱼 的 %d 代最优解: %f, 耗时%f\n", pop, g, Bulletin.FoodConsistence, RunTime[runTime])
				}

				checkResult()
			}
		}

	}
}

func checkResult() {
	var value float64
	var capacities = make([]int, BAG_NUM)
	var error = 0

	for i := 0; i < OBJECT_NUM; i++ {
		if Bulletin.Object[i] != UNSELECTED {
			capacities[Bulletin.Object[i]] += Weight[i]
			value += float64(Dim[Bulletin.Object[i]][i])
		}
	}

	for i := 0; i < BAG_NUM; i++ {
		if capacities[i] > Capacity_Limit[i] {
			fmt.Printf("bag is %v\n", i)
			fmt.Printf("limit is %v, current is %v\n", Capacity_Limit[i], capacities[i])
			fmt.Printf("it is invalid result.\n")
			error++
		}
	}

	if error != 0 {
		fmt.Printf("it is a invalid result, there are %d bags\n", error)
	} else {
		fmt.Printf("it is valid result. fitness is %f, generated fitness is %f\n", value, Bulletin.FoodConsistence)
	}
}
