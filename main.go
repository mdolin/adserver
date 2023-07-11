package main

import (
	"adserver/db"
	"adserver/handlers"
	"adserver/models"
	"log"
	"net/http"
)

func main() {
	// Create a new database connection
	db, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Populate the database with sample AdUnits and Creatives
	err = populateDB(db)
	if err != nil {
		log.Fatal(err)
	}

	// Register the AdHandler and RefreshHandler with the HTTP server
	http.HandleFunc("/adrequest", func(w http.ResponseWriter, r *http.Request) {
		handlers.AdHandler(w, r, db)
	})
	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		handlers.RefreshHandler(w, r, db)
	})

	// Start the HTTP server
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// populateDB inserts sample AdUnits and Creatives into the database
func populateDB(db *db.DB) error {
	adUnit1 := models.AdUnit{
		ID:     "adunit1",
		Format: models.Banner,
		Width:  300,
		Height: 250,
	}

	adUnit2 := models.AdUnit{
		ID:     "adunit2",
		Format: models.Interstitial,
		Width:  1024,
		Height: 768,
	}

	adUnit3 := models.AdUnit{
		ID:     "adunit3",
		Format: models.Video,
		Width:  1000,
		Height: 700,
	}

	creative1 := models.Creative{
		ID:      "creative1",
		Format:  models.Banner,
		Width:   300,
		Height:  250,
		Content: "Sample Banner Ad",
		Price:   1.5,
	}

	creative2 := models.Creative{
		ID:      "creative2",
		Format:  models.Interstitial,
		Width:   1024,
		Height:  768,
		Content: "Sample Interstitial Ad",
		Price:   3.0,
	}

	err := db.InsertAdUnit(adUnit1)
	if err != nil {
		return err
	}

	err = db.InsertAdUnit(adUnit2)
	if err != nil {
		return err
	}

	err = db.InsertAdUnit(adUnit3)
	if err != nil {
		return err
	}

	err = db.InsertCreative(creative1)
	if err != nil {
		return err
	}

	err = db.InsertCreative(creative2)
	if err != nil {
		return err
	}

	return nil
}
