package repositories

import (
	"database/sql"
	"cashier/models"
	"fmt"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		query := "SELECT name, price, stock FROM products WHERE id = $1"
		if useLock {
			query += " FOR UPDATE"
		}

		err := tx.QueryRow(query, item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s (available: %d, requested: %d)", productName, stock, item.Quantity)
		}
  
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetReport(startDate, endDate time.Time) (*models.Report, error) {
	var report models.Report

	// 1. Total Revenue & Total Transactions
	err := repo.db.QueryRow(`
        SELECT COALESCE(SUM(total_amount), 0), COUNT(id)
        FROM transactions
        WHERE created_at BETWEEN $1 AND $2
    `, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaction)
	if err != nil {
		return nil, err
	}

	// 2. Most Sold Product
	err = repo.db.QueryRow(`
        SELECT p.name, COALESCE(SUM(td.quantity), 0) as sold_qty
        FROM transaction_details td
        JOIN products p ON td.product_id = p.id
        JOIN transactions t ON td.transaction_id = t.id
        WHERE t.created_at BETWEEN $1 AND $2
        GROUP BY p.name
        ORDER BY sold_qty DESC
        LIMIT 1
    `, startDate, endDate).Scan(&report.MostSoldProduct.Name, &report.MostSoldProduct.SoldQty)

	if err == sql.ErrNoRows {
		report.MostSoldProduct = models.MostSoldProduct{Name: "", SoldQty: 0}
	} else if err != nil {
		return nil, err
	}

	return &report, nil
}