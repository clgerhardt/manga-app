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
	UserID               uuid.UUID `json:"user"`
}

func GetAllCollections(conn *pgx.Conn) ([]Collection, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, title, number_of_chapters, most_recent_upload_date, user_id, created_at, updated_at FROM collection")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Error getting collections")
	}

	var collections []Collection
	for rows.Next() {
		collection := Collection{}
		err = rows.Scan(&collection.ID, &collection.Title, &collection.numberOfChapters, &collection.mostRecentUploadDate, &collection.UserID, &collection.CreatedAt, &collection.UpdatedAt)
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

	now := time.Now()
	row := conn.QueryRow(context.Background(), "INSERT INTO collection (title, number_of_chapters, most_recent_upload_date, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, user_id", i.Title, i.numberOfChapters, now, userID, now, now)

	err := row.Scan(&i.ID, &i.UserID)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("There was an error creating the collection")
	}

	return nil
}

func (i *Collection) Update(conn *pgx.Conn) error {
	i.Title = strings.Trim(i.Title, " ")
	if len(i.Title) < 1 {
		return fmt.Errorf("Title must not be empty")
	}

	now := time.Now()
	_, err := conn.Exec(context.Background(), "UPDATE collection SET title=$1, number_of_chapters=$2, most_recent_upload_date=$3, updated_at=$4 WHERE id=$5", i.Title, i.numberOfChapters, now, now, i.ID)

	if err != nil {
		fmt.Printf("Error updating collection: (%v)", err)
		return fmt.Errorf("Error updating collection")
	}

	return nil
}

func FindCollectionById(id uuid.UUID, conn *pgx.Conn) (Collection, error) {
	row := conn.QueryRow(context.Background(), "SELECT title, number_of_chapters, user_id, most_recent_upload_date FROM collection WHERE id=$1", id)
	collection := Collection{
		ID: id,
	}
	err := row.Scan(&collection.Title, &collection.numberOfChapters, &collection.UserID, &collection.mostRecentUploadDate)
	if err != nil {
		return collection, fmt.Errorf("The collection doesn't exist")
	}

	return collection, nil
}

func (i *Collection) Delete(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM collection WHERE id=$5", i.ID)

	if err != nil {
		fmt.Printf("Error deleting collection: (%v)", err)
		return fmt.Errorf("Error deleting collection")
	}

	return nil
}
