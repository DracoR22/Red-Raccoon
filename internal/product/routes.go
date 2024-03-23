package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DracoR22/Red-Raccoon/types"
	"github.com/DracoR22/Red-Raccoon/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

// DECLARE ROUTES
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products/{productID}", h.handleGetProduct).Methods(http.MethodGet)

	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)
}

// GET ALL PRODUCTS
func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	// Get the products
	products, err := h.store.GetProducts()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

// GET PRODUCT BY ID
func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID
	vars := mux.Vars(r)
	str, ok := vars["productID"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	// Validate product ID
	productID, err := strconv.Atoi(str)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	// Get product by ID
	product, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

// CREATE PRODUCT
func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// Get the payload
	var product types.CreateProductPayload
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate the payload
	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// Create the product
	err := h.store.CreateProduct(product)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product)
}
