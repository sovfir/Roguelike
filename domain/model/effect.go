package model

type Effect struct {
	Item     *Item
	TimeLeft int
}

func NewEffect(item *Item, time int) Effect {
	return Effect{Item: item, TimeLeft: time}
}
