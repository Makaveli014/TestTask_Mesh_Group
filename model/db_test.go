package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	// Use environment variables for test DB connection or defaults
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "mesh_group_test")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Skipf("Skipping integration test: could not connect to test database: %v", err)
	}

	// Create table for testing
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS segmentation (
			id SERIAL PRIMARY KEY,
			address_sap_id VARCHAR(255) NOT NULL UNIQUE,
			adr_segment VARCHAR(16),
			segment_id BIGINT
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Clear table before test
	_, err = db.Exec("DELETE FROM segmentation")
	if err != nil {
		t.Fatalf("Failed to clear test table: %v", err)
	}

	return db
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func TestUpsertSegmentations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// 1. Initial insert
	data := []Segmentation{
		{AddressSapID: "sap1", AdrSegment: "seg1", SegmentID: 100},
		{AddressSapID: "sap2", AdrSegment: "seg2", SegmentID: 200},
	}

	err := UpsertSegmentations(db, data)
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	// 2. Upsert with one existing and one new
	updateData := []Segmentation{
		{AddressSapID: "sap1", AdrSegment: "seg1-updated", SegmentID: 101}, // Update
		{AddressSapID: "sap3", AdrSegment: "seg3", SegmentID: 300},       // New
	}

	err = UpsertSegmentations(db, updateData)
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	// Verify results
	var count int
	err = db.Get(&count, "SELECT count(*) FROM segmentation")
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("Expected 3 records in DB, got %d", count)
	}

	var seg1 Segmentation
	err = db.Get(&seg1, "SELECT * FROM segmentation WHERE address_sap_id = 'sap1'")
	if err != nil {
		t.Fatal(err)
	}
	if seg1.AdrSegment != "seg1-updated" || seg1.SegmentID != 101 {
		t.Errorf("Expected updated data for sap1, got AdrSegment=%s, SegmentID=%d", seg1.AdrSegment, seg1.SegmentID)
	}
}
