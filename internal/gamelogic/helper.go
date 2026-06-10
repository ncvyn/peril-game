package gamelogic

func getUnitPlurality(units map[int]Unit) string {
	if len(units) == 1 {
		return "unit"
	}
	return "units"
}

func getArticle(rank string) string {
	switch rank {
	case string(RankCavalry):
		return "a"
	default:
		return "an"
	}
}
