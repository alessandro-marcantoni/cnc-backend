package club

import "github.com/alessandro-marcantoni/cnc-backend/main/shared/result"

// SeasonRepository defines the interface for accessing season data
type SeasonRepository interface {
	// GetSeasonById retrieves a season by its ID
	GetSeasonById(seasonId int64) result.Result[Season]
}
