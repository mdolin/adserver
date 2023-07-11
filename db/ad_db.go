package db

import (
	model "adserver/models"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mattn/go-sqlite3"
)

// DB represents the database connection and cache
type DB struct {
	Connection *sql.DB
	adUnits    []model.AdUnit
	creatives  []model.Creative
	cacheMutex sync.RWMutex
}

// NewDB creates a new database connection and initializes the cache
func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite3", "./ad.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Create the necessary tables if they do not exist
	err = createTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Create the cache and populate it initially
	adUnits, err := getAllAdUnits(db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve AdUnits: %w", err)
	}

	creatives, err := getAllCreatives(db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Creatives: %w", err)
	}

	adDB := &DB{
		Connection: db,
		adUnits:    adUnits,
		creatives:  creatives,
	}

	// Start a goroutine to periodically refresh the cache
	go adDB.refreshCacheRoutine()

	return adDB, nil
}

// refreshCacheRoutine periodically refreshes the cache from the database
func (db *DB) refreshCacheRoutine() {
	for {
		err := db.RefreshCache()
		if err != nil {
			log.Printf("Failed to refresh cache: %v", err)
		}

		// Wait for 5 minutes before the next cache refresh
		time.Sleep(5 * time.Minute)
	}
}

// refreshCache refreshes the cache by retrieving AdUnits and Creatives from the database
func (db *DB) RefreshCache() error {
	adUnits, err := getAllAdUnits(db.Connection)
	if err != nil {
		return fmt.Errorf("failed to retrieve AdUnits: %w", err)
	}

	creatives, err := getAllCreatives(db.Connection)
	if err != nil {
		return fmt.Errorf("failed to retrieve Creatives: %w", err)
	}

	db.cacheMutex.Lock()
	db.adUnits = adUnits
	db.creatives = creatives
	db.cacheMutex.Unlock()

	return nil
}

// InsertAdUnit inserts an AdUnit into the database
func (db *DB) InsertAdUnit(adUnit model.AdUnit) error {
	_, err := db.Connection.Exec("INSERT INTO AdUnits (ID, Format, Width, Height) VALUES (?, ?, ?, ?)",
		adUnit.ID, adUnit.Format, adUnit.Width, adUnit.Height)
	if err != nil {
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.Code == sqlite3.ErrConstraint {
			return fmt.Errorf("a AdUnit with ID '%s' already exists", adUnit.ID)
		}
		return fmt.Errorf("failed to insert AdUnit into the database: %w", err)
	}

	// Update the cache
	db.cacheMutex.Lock()
	db.adUnits = append(db.adUnits, adUnit)
	db.cacheMutex.Unlock()

	return nil
}

// InsertCreative inserts a Creative into the database
func (db *DB) InsertCreative(creative model.Creative) error {
	_, err := db.Connection.Exec("INSERT INTO Creatives (ID, Format, Width, Height, Content, Price) VALUES (?, ?, ?, ?, ?, ?)",
		creative.ID, creative.Format, creative.Width, creative.Height, creative.Content, creative.Price)
	if err != nil {
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.Code == sqlite3.ErrConstraint {
			return fmt.Errorf("a Creative with ID '%s' already exists", creative.ID)
		}
		return fmt.Errorf("failed to insert Creative into the database: %w", err)
	}

	// Update the cache
	db.cacheMutex.Lock()
	db.creatives = append(db.creatives, creative)
	db.cacheMutex.Unlock()

	return nil
}

// GetAdUnitByID retrieves an AdUnit from the cache by ID
func (db *DB) GetAdUnitByID(id string) (model.AdUnit, error) {
	db.cacheMutex.RLock()
	defer db.cacheMutex.RUnlock()

	for _, adUnit := range db.adUnits {
		if adUnit.ID == id {
			return adUnit, nil
		}
	}

	return model.AdUnit{}, fmt.Errorf("AdUnit not found")
}

// GetCreativeByID retrieves a Creative from the cache by ID
func (db *DB) GetCreativeByID(id string) (model.Creative, error) {
	db.cacheMutex.RLock()
	defer db.cacheMutex.RUnlock()

	for _, creative := range db.creatives {
		if creative.ID == id {
			return creative, nil
		}
	}

	return model.Creative{}, fmt.Errorf("creative not found")
}

// GetCreatives retrieves all Creatives from the cache
func (db *DB) GetCreatives() ([]model.Creative, error) {
	db.cacheMutex.RLock()
	defer db.cacheMutex.RUnlock()

	if len(db.creatives) == 0 {
		return nil, fmt.Errorf("no Creatives found in cache")
	}

	creatives := make([]model.Creative, len(db.creatives))
	copy(creatives, db.creatives)

	return creatives, nil
}

// Close the database connection
func (db *DB) Close() error {
	err := db.Connection.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// CreateTables create the necessary tables if they do not exist
func createTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS AdUnits (
		ID TEXT PRIMARY KEY,
		Format TEXT,
		Width INT,
		Height INT
	)`)
	if err != nil {
		return fmt.Errorf("failed to create AdUnits table: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Creatives (
		ID TEXT PRIMARY KEY,
		Format TEXT,
		Width INT,
		Height INT,
		Content TEXT,
		Price REAL
	)`)
	if err != nil {
		return fmt.Errorf("failed to create Creatives table: %w", err)
	}

	return nil
}

// getAllAdUnits retrieves all AdUnits from the database
func getAllAdUnits(db *sql.DB) ([]model.AdUnit, error) {
	rows, err := db.Query("SELECT ID, Format, Width, Height FROM AdUnits")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve AdUnits: %w", err)
	}
	defer rows.Close()

	var adUnits []model.AdUnit

	for rows.Next() {
		var adUnit model.AdUnit
		err := rows.Scan(&adUnit.ID, &adUnit.Format, &adUnit.Width, &adUnit.Height)
		if err != nil {
			return nil, fmt.Errorf("failed to scan AdUnit: %w", err)
		}

		adUnits = append(adUnits, adUnit)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}

	return adUnits, nil
}

// getAllCreatives retrieves all Creatives from the database
func getAllCreatives(db *sql.DB) ([]model.Creative, error) {
	rows, err := db.Query("SELECT ID, Format, Width, Height, Content, Price FROM Creatives")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Creatives: %w", err)
	}
	defer rows.Close()

	var creatives []model.Creative

	for rows.Next() {
		var creative model.Creative
		err := rows.Scan(&creative.ID, &creative.Format, &creative.Width, &creative.Height, &creative.Content, &creative.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan Creative: %w", err)
		}

		creatives = append(creatives, creative)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}

	return creatives, nil
}
