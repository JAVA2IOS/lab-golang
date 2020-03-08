package controller

import (
	"net/http"
	"lab-golang/upload/controller"
)

func QqlCreateShippingOrder(w http.ResponseWriter, r *http.Request) {
	_ = uploadController.SaveFile(w, r, uploadController.LabFileTypeCSV)

}