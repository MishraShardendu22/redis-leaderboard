package view

import "fmt"

func podiumClass(rank int) string {
	switch rank {
	case 1:
		return "podium-gold"
	case 2:
		return "podium-silver"
	case 3:
		return "podium-bronze"
	default:
		return "podium-neutral"
	}
}

func podiumLabel(rank int) string {
	switch rank {
	case 1:
		return "Gold"
	case 2:
		return "Silver"
	case 3:
		return "Bronze"
	default:
		return "Rank"
	}
}

func movementClass(delta int64) string {
	switch {
	case delta > 0:
		return "movement-up"
	case delta < 0:
		return "movement-down"
	default:
		return "movement-flat"
	}
}

func movementLabel(delta int64) string {
	switch {
	case delta > 0:
		return fmt.Sprintf("▲ %d", delta)
	case delta < 0:
		return fmt.Sprintf("▼ %d", -delta)
	default:
		return "• 0"
	}
}
