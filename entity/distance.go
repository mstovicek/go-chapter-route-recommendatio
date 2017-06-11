package entity

type Distance struct {
	FromPlaceID     string
	ToPlaceID       string
	DistanceMetres  int
	DurationSeconds float64
}

type DistanceMatrix map[string]map[string]Distance

func NewDistanceMatrix() DistanceMatrix {
	return make(DistanceMatrix)
}

func (distanceMatrix DistanceMatrix) Add(from string, to string, distance Distance) {
	_, ok := distanceMatrix[from]
	if !ok {
		distanceMatrix[from] = make(map[string]Distance)
	}
	distanceMatrix[from][to] = distance
}

func (distanceMatrix DistanceMatrix) Get(from string, to string) (Distance, bool) {
	distance, ok := distanceMatrix[from][to]
	return distance, ok
}
