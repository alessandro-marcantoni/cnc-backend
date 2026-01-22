package persistence

import (
	"context"
	"database/sql"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain/club"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

type SQLSeasonRepository struct {
	db *sql.DB
}

func NewSQLSeasonRepository(db *sql.DB) *SQLSeasonRepository {
	return &SQLSeasonRepository{db: db}
}

func (r *SQLSeasonRepository) GetSeasonById(seasonId int64) result.Result[club.Season] {
	query := `
		SELECT id, code, name, starts_at, ends_at
		FROM seasons
		WHERE id = $1
	`

	var season club.Season
	err := r.db.QueryRowContext(context.Background(), query, seasonId).Scan(
		&season.ID,
		&season.Code,
		&season.Name,
		&season.StartsAt,
		&season.EndsAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return result.Err[club.Season](errors.RepositoryError{Description: "season not found"})
		}
		return result.Err[club.Season](errors.RepositoryError{Description: err.Error()})
	}

	return result.Ok(season)
}
