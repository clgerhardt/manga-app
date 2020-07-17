package models

import (
	"fmt"
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
