package constants

type Directions struct {
	DeltaY int
	DeltaX int
}

func LeftUp() Directions {
	return Directions{-1, -1}
}

func Up() Directions {
	return Directions{-1, 0}
}

func RightUp() Directions {
	return Directions{-1, 1}
}

func LeftDown() Directions {
	return Directions{1, -1}
}

func Down() Directions {
	return Directions{1, 0}
}

func RightDown() Directions {
	return Directions{1, 1}
}

func Left() Directions {
	return Directions{0, -1}
}

func Right() Directions {
	return Directions{0, 1}
}
