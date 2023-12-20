package main

import (
	"database/sql"
)

// Repository представляет собой структуру для работы с базой данных.
type Repository struct {
	db *sql.DB
}

// NewRepository создаёт новый экземпляр Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CheckAddressInDatabase проверяет наличие адреса в базе данных.
func (r *Repository) CheckAddressInDatabase(address string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM address WHERE address_text = $1"
	err := r.db.QueryRow(query, address).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// SaveSearchHistory сохраняет запрос поиска в истории поиска.
func (r *Repository) SaveSearchHistory(searchQuery string) error {
	query := "INSERT INTO search_history (search_query) VALUES ($1)"
	_, err := r.db.Exec(query, searchQuery)
	return err
}

// SaveAddress сохраняет адрес и возвращает его идентификатор.
func (r *Repository) SaveAddress(address string) (int, error) {
	var addressID int
	query := "INSERT INTO address (address_text) VALUES ($1) RETURNING id"
	err := r.db.QueryRow(query, address).Scan(&addressID)
	return addressID, err
}

// LinkAddressToSearchHistory связывает адрес с историей поиска.
func (r *Repository) LinkAddressToSearchHistory(searchHistoryID, addressID int) error {
	query := "INSERT INTO history_search_address (search_history_id, address_id) VALUES ($1, $2)"
	_, err := r.db.Exec(query, searchHistoryID, addressID)
	return err
}

// GetAddressesFromSearchHistory возвращает адреса, соответствующие запросу поиска.
func (r *Repository) GetAddressesFromSearchHistory(searchQuery string, similarityThreshold float64) ([]string, error) {
	query := `
        SELECT address.address_text
        FROM search_history
        JOIN history_search_address ON search_history.id = history_search_address.search_history_id
        JOIN address ON history_search_address.address_id = address.id
        WHERE search_history.search_query ILIKE $1
          AND similarity(search_history.search_query, $2) > $3
    `
	rows, err := r.db.Query(query, "%"+searchQuery+"%", searchQuery, similarityThreshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []string
	for rows.Next() {
		var address string
		if err := rows.Scan(&address); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}
