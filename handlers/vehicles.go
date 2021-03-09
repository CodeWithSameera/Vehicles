package handlers

import (
	"encoding/json"
	"github.com/CodeWithSameera/Vehicles/helpers"
	"github.com/CodeWithSameera/Vehicles/model"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

var db  *gorm.DB

var GetAllVehiclesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	// Here we are converting the slice of products to JSON
	vehicles := []model.Vehicle{}
	db = helpers.ConnectDB()
	db.Find(&vehicles)
	respondJSON(w, http.StatusOK, vehicles)
})
var SearchVehiclesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)

	vno := vars["vno"]
	db = helpers.ConnectDB()
	vehicle := getVehicleOr404(db, vno, w, r)
	if vehicle == nil {
		return
	}
	respondJSON(w, http.StatusOK, vehicle)
})

var CreateVehiclesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	vehicle := model.Vehicle{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&vehicle); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	db = helpers.ConnectDB()
	if err := db.Save(&vehicle).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	SaveEntity(vehicle.VNo,vehicle.Modal)

	respondJSON(w, http.StatusCreated, vehicle)
})

var UpdateVehiclesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)

	vno := vars["vno"]
	db = helpers.ConnectDB()
	vehicle := getVehicleOr404(db, vno, w, r)
	if vehicle == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&vehicle); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&vehicle).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, vehicle)
})

var DeleteVehiclesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	db = helpers.ConnectDB()
	vno := vars["vno"]
	vehicle := getVehicleOr404(db, vno, w, r)
	if vehicle == nil {
		return
	}
	if err := db.Delete(&vehicle).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
})

func getVehicleOr404(db *gorm.DB, vno string, w http.ResponseWriter, r *http.Request) *model.Vehicle {
	vehicle := model.Vehicle{}
	if err := db.First(&vehicle, model.Vehicle{VNo: vno}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &vehicle
}