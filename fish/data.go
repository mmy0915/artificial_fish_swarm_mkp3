package fish

const (
	UNSELECTED    = -1    //物品未被选中
	POPSIZE       = 100   //鱼群规模
	VISUAL        = 6     //感知距离
	ATTEMPT       = 200   //尝试次数
	DELTA         = 0.938 //因子
	LINE_Neighbor = 0.3
	LINE_Fitness  = 0.7
)

var OBJECT_NUM int //物品数量
var BAG_NUM int    //背包数量

var Capacity_Limit []int

var Bulletin ArtificialFish

var Dim [][]int

var Weight []int
