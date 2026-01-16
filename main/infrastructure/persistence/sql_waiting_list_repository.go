package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	facilityrental "github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

//go:embed queries/get_waiting_list.sql
var getWaitingListQuery string

//go:embed queries/get_next_waiting_entry.sql
var getNextWaitingEntryQuery string

//go:embed queries/add_waiting_entry.sql
var addWaitingEntryQuery string

//go:embed queries/remove_waiting_entry.sql
var removeWaitingEntryQuery string

//go:embed queries/remove_waiting_entry_by_member_and_type.sql
var removeWaitingEntryByMemberAndTypeQuery string

//go:embed queries/get_member_waiting_entry.sql
var getMemberWaitingEntryQuery string

type SQLWaitingListRepository struct {
	db *sql.DB
}

func NewSQLWaitingListRepository(db *sql.DB) *SQLWaitingListRepository {
	return &SQLWaitingListRepository{db: db}
}

func (r *SQLWaitingListRepository) GetWaitingList(facilityType domain.Id[facilityrental.FacilityType]) result.Result[facilityrental.WaitingList] {
	rows, err := r.db.Query(getWaitingListQuery, facilityType.Value)
	if err != nil {
		return result.Err[facilityrental.WaitingList](errors.RepositoryError{Description: "failed to get waiting list: " + err.Error()})
	}
	defer rows.Close()

	var entries []facilityrental.WaitingListEntry

	for rows.Next() {
		var id int64
		var memberId int64
		var facilityTypeId int64
		var queuedAt time.Time
		var notes sql.NullString

		if err := rows.Scan(&id, &memberId, &facilityTypeId, &queuedAt, &notes); err != nil {
			return result.Err[facilityrental.WaitingList](errors.RepositoryError{Description: "failed to scan waiting list entry: " + err.Error()})
		}

		entry := facilityrental.ReconstructWaitingListEntry(
			domain.Id[facilityrental.WaitingListEntry]{Value: id},
			domain.Id[membership.Member]{Value: memberId},
			domain.Id[facilityrental.FacilityType]{Value: facilityTypeId},
			queuedAt,
			notes.String,
		)

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return result.Err[facilityrental.WaitingList](errors.RepositoryError{Description: "error iterating waiting list: " + err.Error()})
	}

	// Note: We don't load the full FacilityType here, just create a placeholder
	// The caller can enrich this if needed
	waitingList := facilityrental.WaitingList{
		FacilityType: facilityrental.FacilityType{},
		Entries:      entries,
	}

	return result.Ok(waitingList)
}

func (r *SQLWaitingListRepository) GetNextEntry(facilityType domain.Id[facilityrental.FacilityType]) result.Result[facilityrental.WaitingListEntry] {
	var id int64
	var memberId int64
	var facilityTypeId int64
	var queuedAt time.Time
	var notes sql.NullString

	err := r.db.QueryRow(getNextWaitingEntryQuery, facilityType.Value).Scan(&id, &memberId, &facilityTypeId, &queuedAt, &notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return result.Err[facilityrental.WaitingListEntry](errors.NotFoundError{Description: "no entries in waiting list"})
		}
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to get next waiting entry: " + err.Error()})
	}

	entry := facilityrental.ReconstructWaitingListEntry(
		domain.Id[facilityrental.WaitingListEntry]{Value: id},
		domain.Id[membership.Member]{Value: memberId},
		domain.Id[facilityrental.FacilityType]{Value: facilityTypeId},
		queuedAt,
		notes.String,
	)

	return result.Ok(entry)
}

func (r *SQLWaitingListRepository) AddEntry(entry facilityrental.WaitingListEntry) result.Result[facilityrental.WaitingListEntry] {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to begin transaction: " + err.Error()})
	}
	defer tx.Rollback()

	var id int64
	var queuedAt time.Time
	notesValue := sql.NullString{String: entry.Notes, Valid: entry.Notes != ""}

	err = tx.QueryRowContext(
		ctx,
		addWaitingEntryQuery,
		entry.MemberId.Value,
		entry.FacilityType.Value,
		notesValue,
	).Scan(&id, &queuedAt)

	if err != nil {
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to add waiting entry: " + err.Error()})
	}

	if err := tx.Commit(); err != nil {
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to commit transaction: " + err.Error()})
	}

	savedEntry := facilityrental.ReconstructWaitingListEntry(
		domain.Id[facilityrental.WaitingListEntry]{Value: id},
		entry.MemberId,
		entry.FacilityType,
		queuedAt,
		entry.Notes,
	)

	return result.Ok(savedEntry)
}

func (r *SQLWaitingListRepository) RemoveEntry(entryId domain.Id[facilityrental.WaitingListEntry]) result.Result[facilityrental.WaitingListEntry] {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to begin transaction: " + err.Error()})
	}
	defer tx.Rollback()

	var id int64
	var memberId int64
	var facilityTypeId int64
	var queuedAt time.Time
	var notes sql.NullString

	err = tx.QueryRowContext(ctx, removeWaitingEntryQuery, entryId.Value).Scan(&id, &memberId, &facilityTypeId, &queuedAt, &notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return result.Err[facilityrental.WaitingListEntry](errors.NotFoundError{Description: "waiting list entry not found"})
		}
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to remove waiting entry: " + err.Error()})
	}

	if err := tx.Commit(); err != nil {
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to commit transaction: " + err.Error()})
	}

	entry := facilityrental.ReconstructWaitingListEntry(
		domain.Id[facilityrental.WaitingListEntry]{Value: id},
		domain.Id[membership.Member]{Value: memberId},
		domain.Id[facilityrental.FacilityType]{Value: facilityTypeId},
		queuedAt,
		notes.String,
	)

	return result.Ok(entry)
}

func (r *SQLWaitingListRepository) RemoveEntryByMemberAndType(facilityType domain.Id[facilityrental.FacilityType], memberId domain.Id[membership.Member]) result.Result[facilityrental.WaitingListEntry] {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to begin transaction: " + err.Error()})
	}
	defer tx.Rollback()

	var id int64
	var queuedAt time.Time
	var notes sql.NullString

	err = tx.QueryRowContext(ctx, removeWaitingEntryByMemberAndTypeQuery, memberId.Value, facilityType.Value).Scan(&id, &queuedAt, &notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return result.Err[facilityrental.WaitingListEntry](errors.NotFoundError{Description: "waiting list entry not found"})
		}
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to remove waiting entry: " + err.Error()})
	}

	if err := tx.Commit(); err != nil {
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to commit transaction: " + err.Error()})
	}

	entry := facilityrental.ReconstructWaitingListEntry(
		domain.Id[facilityrental.WaitingListEntry]{Value: id},
		memberId,
		facilityType,
		queuedAt,
		notes.String,
	)

	return result.Ok(entry)
}

func (r *SQLWaitingListRepository) GetMemberEntry(facilityType domain.Id[facilityrental.FacilityType], memberId domain.Id[membership.Member]) result.Result[facilityrental.WaitingListEntry] {
	var id int64
	var facilityTypeId int64
	var queuedAt time.Time
	var notes sql.NullString

	err := r.db.QueryRow(getMemberWaitingEntryQuery, memberId.Value, facilityType.Value).Scan(&id, &facilityTypeId, &queuedAt, &notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return result.Err[facilityrental.WaitingListEntry](errors.NotFoundError{Description: "member not in waiting list"})
		}
		return result.Err[facilityrental.WaitingListEntry](errors.RepositoryError{Description: "failed to get member waiting entry: " + err.Error()})
	}

	entry := facilityrental.ReconstructWaitingListEntry(
		domain.Id[facilityrental.WaitingListEntry]{Value: id},
		memberId,
		domain.Id[facilityrental.FacilityType]{Value: facilityTypeId},
		queuedAt,
		notes.String,
	)

	return result.Ok(entry)
}
