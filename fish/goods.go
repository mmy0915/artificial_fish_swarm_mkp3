package fish

import (
	"bufio"
	"fmt"
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
	unitValue := float64(value)
	if weight != 0 {
		unitValue = float64(value) / float64(weight)
	}
	goods := &Goods{
		Weight:    weight,
		Value:     value,
		UnitValue: unitValue,
	}
	return goods
}

func InitGoods(path string) error {

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("failed to open file " + path)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	fmt.Printf("processing file %s\n", scanner.Text())

	scanner.Scan()
	BAG_NUM, _ = strconv.Atoi(strings.TrimSpace(scanner.Text()))
	scanner.Scan()
	OBJECT_NUM, _ = strconv.Atoi(strings.TrimSpace(scanner.Text()))

	fmt.Printf("there are %d bags, %d objects\n", BAG_NUM, OBJECT_NUM)

	Weight = make([]int, OBJECT_NUM)

	for i := 0; i < OBJECT_NUM; i++ {
		scanner.Scan()
		Weight[i], _ = strconv.Atoi(scanner.Text())
	}

	fmt.Printf("objects Weight is %v\n", Weight)

	Capacity_Limit = make([]int, BAG_NUM)
	for i := 0; i < BAG_NUM; i++ {
		scanner.Scan()
		Capacity_Limit[i], _ = strconv.Atoi(scanner.Text())
	}

	fmt.Printf("bag capacity is %v\n", Capacity_Limit)

	Dim = make([][]int, BAG_NUM)
	for i := 0; i < BAG_NUM; i++ {
		Dim[i] = make([]int, OBJECT_NUM)
	}
	for i := 0; i < BAG_NUM; i++ {
		for j := 0; j < OBJECT_NUM; j++ {
			scanner.Scan()
			Dim[i][j], _ = strconv.Atoi(scanner.Text())
		}
	}

/*	for i := 0; i < OBJECT_NUM; i++ {
		scanner.Scan()
		Value[i], _ = strconv.Atoi(scanner.Text())
		for j := 0; j < BAG_NUM; j++ {
			scanner.Scan()
			Dim[j][i], _ = strconv.Atoi(scanner.Text())
		}
	}

	Capacity_Limit = make([]int, BAG_NUM)
	for i := 0; i < BAG_NUM; i++ {
		scanner.Scan()
		Capacity_Limit[i], _ = strconv.Atoi(scanner.Text())
	}*/

	//sortGoods()

	return nil
}

/*func sortGoods() {
	for i := 0; i < OBJECT_NUM; i++ {
		for j := 0; j < OBJECT_NUM-i-1; j++ {
			if AllGoods[j].UnitValue < AllGoods[j+1].UnitValue {
				AllGoods[j], AllGoods[j+1] = AllGoods[j+1], AllGoods[j]
			}
		}
	}
}*/
