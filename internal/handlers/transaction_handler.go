package handlers

import (
	"encoding/json"
	"kasir-api/internal/models"
	"kasir-api/internal/service"
	"net/http"
	"time"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Checkout godoc
// @Summary Create a new transaction (Checkout)
// @Description Create a new transaction with multiple items
// @Tags transaction
// @Accept json
// @Produce json
// @Param checkout body models.CheckoutRequest true "Checkout Data"
// @Success 201 {object} models.Transaction
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /checkout [post]
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// GetDailyReport godoc
// @Summary Get daily sales report
// @Description Get total revenue, total transactions, and best selling product for today
// @Tags report
// @Accept json
// @Produce json
// @Success 200 {object} models.SalesReport
// @Failure 500 {string} string "Internal Server Error"
// @Router /report/hari-ini [get]
func (h *TransactionHandler) HandleDailyReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	report, err := h.service.GetDailyReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetReport godoc
// @Summary Get sales report by date range
// @Description Get total revenue, total transactions, and best selling product for a specific date range
// @Tags report
// @Accept json
// @Produce json
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {object} models.SalesReport
// @Failure 400 {string} string "Invalid date format"
// @Failure 500 {string} string "Internal Server Error"
// @Router /report [get]
func (h *TransactionHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid start_date format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid end_date format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Set end date time to 23:59:59
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	report, err := h.service.GetReport(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
