package handlers

import (
	"adserver/db"
	model "adserver/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// AdHandler handles the ad request
func AdHandler(w http.ResponseWriter, r *http.Request, db *db.DB) {
	// Parse the JSON body
	var request model.AdRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to parse JSON: %s", err.Error())
		return
	}

	// Generate a user ID if not provided in the request
	if request.UserID == "" {
		request.UserID = generateUserID()
	}

	// Retrieve the AdUnit from the cache
	adUnit, err := db.GetAdUnitByID(request.AdUnitID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to retrieve AdUnit: %s", err.Error())
		return
	}

	// Retrieve the Creatives from the cache
	creatives, err := db.GetCreatives()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to retrieve Creatives: %s", err.Error())
		return
	}

	// Find the applicable creative with the highest price
	selectedCreative := findApplicableCreative(adUnit, creatives)
	if selectedCreative == nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "No applicable creative found")
		return
	}

	// Prepare the response
	response := prepareResponse(selectedCreative, request.UserID)

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to encode JSON response: %s", err.Error())
		return
	}
}

// RefreshHandler handles the refresh request to update the cache values
func RefreshHandler(w http.ResponseWriter, r *http.Request, db *db.DB) {
	err := db.RefreshCache()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to refresh cache: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Cache refreshed successfully")
}

// findApplicableCreative finds the applicable creative with the highest price
func findApplicableCreative(adUnit model.AdUnit, creatives []model.Creative) *model.Creative {
	var selectedCreative *model.Creative
	highestPrice := 0.0

	for _, creative := range creatives {
		if (creative.Format == adUnit.Format) &&
			(creative.Width == adUnit.Width) &&
			(creative.Height == adUnit.Height) {
			if creative.Price > highestPrice {
				tempCreative := creative // Create a temporary copy
				selectedCreative = &tempCreative
				highestPrice = creative.Price
			}
		}
	}

	if selectedCreative == nil {
		return nil
	}

	return selectedCreative
}

// prepareResponse prepares the response in the required format
func prepareResponse(selectedCreative *model.Creative, userID string) model.AdResponse {
	response := model.AdResponse{
		CreativeID: selectedCreative.ID,
		Content:    selectedCreative.Content,
		Price:      selectedCreative.Price,
		UserID:     userID,
	}

	return response
}

// generateUserID generates a new user ID using UUID
func generateUserID() string {
	userID := uuid.New().String()
	return userID
}
