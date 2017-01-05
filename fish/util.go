package fish

func DeepCopy(dst, src *ArtificialFish) {
	copy(dst.Capacity, src.Capacity)
	copy(dst.NeighborFitness, src.NeighborFitness)
	copy(dst.Object, src.Object)
	dst.FoodConsistence = src.FoodConsistence
	dst.NeighborNum = src.NeighborNum
//	dst.Visual = src.Visual
}
