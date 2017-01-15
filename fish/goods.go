package fish

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Goods struct {
	Weight    int
	Value     int
	UnitValue float64
}

func NewGoods(weight, value int) *Goods {
	goods := &Goods{
		Weight:    weight,
		Value:     value,
		UnitValue: float64(value) / float64(weight),
	}
	return goods
}

func InitGoods(path string) error {

	OBJECT_NUM = 0
	BAG_NUM = 0

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("failed to open file " + path)
		return err
	}
	defer file.Close()

	bfRd := bufio.NewReader(file)
	for {
		line, err := bfRd.ReadString('\n')

		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				break
			}
			return err
		}
		items := strings.Split(line, ",")

		if len(items) == 2 {
			weight, _ := strconv.Atoi(strings.TrimSpace(items[0]))
			value, _ := strconv.Atoi(strings.TrimSpace(items[1]))
			AllGoods = append(AllGoods, NewGoods(weight, value))
			OBJECT_NUM++
		} else if len(items) == 1 && strings.TrimSpace(items[0]) != "" {
			cap, _ := strconv.Atoi(strings.TrimSpace(items[0]))
			Capacity_Limit = append(Capacity_Limit, cap)
			BAG_NUM++
		}
	}

	sortGoods()

	return nil
}

func sortGoods() {
	for i := 0; i < OBJECT_NUM; i++ {
		for j := 0; j < OBJECT_NUM-i-1; j++ {
			if AllGoods[j].UnitValue < AllGoods[j+1].UnitValue {
				AllGoods[j], AllGoods[j+1] = AllGoods[j+1], AllGoods[j]
			}
		}
	}
}
