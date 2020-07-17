package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type Collection struct {
	ID                   uuid.UUID `json:"id"`
	CreatedAt            time.Time `json:"_"`
	UpdatedAt            time.Time `json:"_"`
	Title                string    `json:"title"`
	numberOfChapters     int64     `json:"title"`
	mostRecentUploadDate time.Time `json:"_"`
}

func GetAllCollections(conn *pgx.Conn) ([]Collection, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, title, notes, number_of_chapters, most_recent_upload_date FROM collection")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Error getting items")
	}

	var collections []Collection
	for rows.Next() {
		collection := Collection{}
		err = rows.Scan(&collection.ID, &collection.Title, &collection.numberOfChapters, &collection.mostRecentUploadDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		collections = append(collections, collection)
	}

	return collections, nil
}

func (i *Collection) Create(conn *pgx.Conn, userID string) error {
	i.Title = strings.Trim(i.Title, " ")
	if len(i.Title) < 1 {
		return fmt.Errorf("Title must not be empty.")
	}
	if i.PriceInCents < 0 {
		i.PriceInCents = 0
	}
	now := time.Now()
	fmt.Println("USER ID ", userID)
	row := conn.QueryRow(context.Background(), "INSERT INTO item (title, notes, seller_id, price_in_cents, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, seller_id", i.Title, i.Notes, userID, i.PriceInCents, now, now)

	err := row.Scan(&i.ID, &i.SellerID)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("There was an error creating the item")
	}

	return nil
}
