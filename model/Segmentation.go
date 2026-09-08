package model

import (
	"github.com/jmoiron/sqlx"
)

// Segmentation represents the data structure for the segmentation table
type Segmentation struct {
	ID           int    `db:"id"`
	AddressSapID string `db:"address_sap_id"`
	AdrSegment   string `db:"adr_segment"`
	SegmentID    int64  `db:"segment_id"`
}

// UpsertSegmentations inserts a batch of segmentations into the database,
// updating existing records if address_sap_id already exists.
func UpsertSegmentations(db *sqlx.DB, data []Segmentation) error {
	if len(data) == 0 {
		return nil
	}

	// Using a transaction for batch insertion
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO segmentation (address_sap_id, adr_segment, segment_id)
		VALUES (:address_sap_id, :adr_segment, :segment_id)
		ON CONFLICT (address_sap_id)
		DO UPDATE SET
			adr_segment = EXCLUDED.adr_segment,
			segment_id = EXCLUDED.segment_id
	`

	_, err = tx.NamedExec(query, data)
	if err != nil {
		return err
	}

	return tx.Commit()
}
