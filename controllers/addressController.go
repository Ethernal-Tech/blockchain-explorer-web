package controllers

import (
	"webbc/services"

	"github.com/gin-gonic/gin"
)

type AddressController struct {
	AddressService services.IAddressService
}

func NewAddressController(addressService services.IAddressService) AddressController {
	return AddressController{AddressService: addressService}
}

func (ac *AddressController) GetAddress(context *gin.Context) {
	address, _ := ac.AddressService.GetAddress(context.Param("addresshash"))
	data := gin.H{
		"address": address,
	}
	context.HTML(200, "address.html", data)
}

func (ac *AddressController) UploadABI(context *gin.Context) {
	ac.AddressService.UploadABI(context)
}
