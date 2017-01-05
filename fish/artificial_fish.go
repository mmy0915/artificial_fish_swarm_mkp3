package fish

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var l sync.Mutex

type ArtificialFish struct {
	Object          []int
	Capacity        []int     //各背包所剩容量
	FoodConsistence float64   //食物浓度（背包中总价值）
	NeighborNum     int       //邻居数
	NeighborList    []bool    //邻居列表
	NeighborFitness []float64 //邻居适应度

	Visual int
}

func NewArtificialFish() *ArtificialFish {
	af := &ArtificialFish{
		Object:          make([]int, OBJECT_NUM),
		Capacity:        make([]int, BAG_NUM),
		NeighborFitness: make([]float64, POPSIZE),
		NeighborList:    make([]bool, POPSIZE),
	}

	//给背包容量赋值
	for i := 0; i < BAG_NUM; i++ {
		af.Capacity[i] = Capacity_Limit[i]
	}

	for i := 0; i < OBJECT_NUM; i++ {

		l.Lock()
		BagTag := r.Intn(BAG_NUM) - 1 //-1表示不放入任何包
		l.Unlock()
		if BagTag == UNSELECTED {
			af.Object[i] = UNSELECTED //背包的第flag项赋成未选择
		} else {
			if Weight[i] > af.Capacity[BagTag] { //重量超过，不放入
				af.Object[i] = UNSELECTED
			} else {
				af.Object[i] = BagTag                     //放入
				af.Capacity[BagTag] -= Weight[i] //背包容量减去放入物品重量
			}
		}
	}

	af.Visual = OBJECT_NUM * 3 / 4

	af.EstimateFoodConsistence()
	return af
}

//原始放置算子
//在背包bag中放入物品objectFlag，放入成功返回true，失败返回false
func (af *ArtificialFish) PutInto(bag int, objectID int) {
	//背包总容量充足
	//fmt.Printf("bag is %d, object is %d\n", bag, objectID)
	if objectID < 0 || bag < 0 {
		return
	}
	if af.Object[objectID] == bag {
		return
	}
	if Weight[objectID] < Capacity_Limit[bag] {

		if af.Object[objectID] != UNSELECTED {
			af.TakeOut(objectID) //物品objectID从包中取出
		}

		//包剩余容量充足的情况
		if af.Capacity[bag] >= Weight[objectID] {
			af.Object[objectID] = bag                     //放入背包
			af.Capacity[bag] -= Weight[objectID] //背包剩余容量减去物品重量
		} else { //包容量不足的情况
			//可能取出的物品的位置编号
			endFlag := OBJECT_NUM - 1
			for af.Capacity[bag] < Weight[objectID] { //取出部分物品补充容量
				//l.Lock()
				//j := r.Intn(endFlag) //随机产生j为取出的物品号
				//l.Unlock()
				//if af.Object[tag[j]] != UNSELECTED {
				if af.Object[endFlag] == bag {
					af.TakeOut(endFlag) //取出bag里面第endFlag个物品
				}
				//}
				endFlag--
			}

			if af.Capacity[bag] >= Weight[objectID] {
				af.Object[objectID] = bag
				af.Capacity[bag] -= Weight[objectID] //把物品ObjectFlag放入bag里
			}
		}
	}

	af.Adjust(bag, objectID)
}

func (af *ArtificialFish) Adjust(bag int, objectID int) {
	for i := objectID; i < OBJECT_NUM; i++ {
		if af.Object[i] == UNSELECTED && af.Capacity[bag] >= Weight[objectID] {
			//fmt.Printf("puting object %d to bag %d\n", objectID, bag)
			af.Object[i] = bag
			af.Capacity[bag] -= Weight[objectID]  //把物品ObjectFlag放入bag里
		}
	}
}

//visual内的随机游动
func (af *ArtificialFish) RandomlyMove() {

	limit := float64(OBJECT_NUM) - math.Ceil(af.FoodConsistence/Bulletin.FoodConsistence*float64(OBJECT_NUM)) + 1

	//fmt.Printf("limit is %f\n", limit)
	//随机步长
	l.Lock()
	step := r.Intn(int(math.Abs(limit)) + 1) + 1
	l.Unlock()

	//当前人工鱼的副本，用来计算移动的距离
	meCopy := NewArtificialFish()
	DeepCopy(meCopy, af)

	for {
		l.Lock()
		bagID := r.Intn(BAG_NUM)
		objectID := r.Intn(OBJECT_NUM)
		l.Unlock()

		//随机把nObjectFlag放入bag中
		af.PutInto(bagID, objectID)

		if meCopy.Distance(af) >= step {
			break
		}
	}

	af.EstimateFoodConsistence()
}

func (af *ArtificialFish) Prey() {

	nAttempt := 0
	for {
		meCopy := NewArtificialFish()
		DeepCopy(meCopy, af)

		//visual内试探
		meCopy.RandomlyMove()
		nAttempt++

		//下一个位置适应值比较大
		if meCopy.FoodConsistence > af.FoodConsistence {
			//fmt.Printf("fitness is better after randomly move. before is %f, after is %f\n", af.FoodConsistence, meCopy.FoodConsistence)
			DeepCopy(af, meCopy)
			return
		}
		if nAttempt >= ATTEMPT {
			break
		}
	}

	af.RandomlyMove()
}

func (af *ArtificialFish) Follow(allFish []*ArtificialFish, selfID int) {

	//随机步长
	step := OBJECT_NUM / 6

	cur := 0

	for i := 0; i < OBJECT_NUM; i++ {
		if af.Object[i] != Bulletin.Object[i] && Bulletin.Object[i] != UNSELECTED {
			cur++
			af.PutInto(Bulletin.Object[i], i)
		}
		if cur == step {
			break
		}
	}

	af.EstimateFoodConsistence()
}

func (af *ArtificialFish) Swarm(allFish []*ArtificialFish, selfID int) {

	//1. 计算鱼的中心
	centerID := 0
	mostNeighborCnt := 0
	var neighborNumArray = make([]int, POPSIZE)

	//1. 统计邻居个数
	for i := 0; i < POPSIZE; i++ {
		allFish[i].UpdateNeighborList(allFish, i)
		neighborNumArray[i] = allFish[i].NeighborNum
	}

	//计算中心
	mostNeighborCnt = allFish[0].NeighborNum

	for i := 0; i < POPSIZE; i++ {
		if i != selfID {
			if allFish[i].NeighborNum > mostNeighborCnt {
				mostNeighborCnt = allFish[i].NeighborNum
				centerID = i
			}
		}
	}

	//2. 统计自己的邻居个数
	af.UpdateNeighborList(allFish, selfID)

	//3. 根据条件执行相应动作
	for i := 0; i < OBJECT_NUM; i++ {
		if af.Object[i] != allFish[centerID].Object[i] {
			af.PutInto(allFish[centerID].Object[i], i)
		}
	}

	af.EstimateFoodConsistence()
}

//Get FoodConsistence
func (af *ArtificialFish) EstimateFoodConsistence() float64 {
	af.FoodConsistence = 0.0
	for i := 0; i < OBJECT_NUM; i++ {
		if af.Object[i] != UNSELECTED {
			af.FoodConsistence += float64(Dim[af.Object[i]][i])
		}
	}

	return af.FoodConsistence
}

//Update Neighborhood information
func (af *ArtificialFish) UpdateNeighborList(allFish []*ArtificialFish, selfID int) {
	af.NeighborNum = 0
	for i := 0; i < POPSIZE; i++ {
		if selfID == i {
			continue
		}
		//fmt.Printf("all fish distance is %d, visual is %d\n", af.Distance(allFish[i]), af.Visual)
		if af.Distance(allFish[i]) < af.Visual {
			af.NeighborList[i] = true
			af.NeighborNum++
		} else {
			af.NeighborList[i] = false
		}
	}
}

//Calculate the distance between the fish.
func (af *ArtificialFish) Distance(other *ArtificialFish) int {
	var distance = 0
	for i := 0; i < OBJECT_NUM; i++ {
		if af.Object[i] != other.Object[i] {
			distance++
		}
	}
	return distance
}

//Update the bulletin data.
func (af *ArtificialFish) UpdateSpecimen() {
	if af.FoodConsistence > Bulletin.FoodConsistence {
		Bulletin = *af
	}
}

//Show fish's self information.
func (af *ArtificialFish) ShowInfo() {
	for i := 0; i < BAG_NUM; i++ {
		fmt.Printf("\n背包 %v 所剩容量: %v", i, af.Capacity[i])
	}

	fmt.Printf("\n总价值: %v\n", af.FoodConsistence)
}

func (af *ArtificialFish) Escape() {
	var meCopy = NewArtificialFish()
	for {
		DeepCopy(meCopy, af)
		meCopy.RandomlyMove()
		if meCopy.Distance(af) > af.Visual {
			break
		}
	}
	DeepCopy(af, meCopy)
}

func (af *ArtificialFish) BehaviorSelection(allFish []*ArtificialFish, selfID int) {

	//更新邻居信息
	af.UpdateNeighborList(allFish, selfID)
/*
	if float64(af.NeighborNum) <= LINE_Neighbor*float64(POPSIZE) {
		af.Follow(allFish, selfID)
	} else { //邻居多
		fmt.Println("fish is preying. NeighborNum is too big")
		af.Prey()
	}*/

	if af.NeighborNum != 0 {
		af.Follow(allFish, selfID)
	}

	af.Prey()
}

//从背包中取出物品
func (af *ArtificialFish) TakeOut(objectID int) {

	//record bag's id
	bagTag := af.Object[objectID]

	//make object empty
	af.Object[objectID] = UNSELECTED

	//increase bag's capacity
	af.Capacity[bagTag] += Weight[objectID]
}
