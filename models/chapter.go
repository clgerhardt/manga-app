package models

import (
	"fmt"
	"time"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type Chapter struct {
	ID                   uuid.UUID `json:"id"`
	CreatedAt            time.Time `json:"_"`
	UpdatedAt            time.Time `json:"_"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	NumberOfPages        int64     `json:"number_of_pages"`
	ChapterNumber        int64     `json:"chapter_number"`
	CollectionID         uuid.UUID `json:"collection"`
}

// func (i *Chapter) GetAllChapters(conn *pgx.Conn, collectionID uuid.UUID) ([]Chapter, error) {
// 	rows, err := conn.Query(context.Background(), "SELECT id, title, description, number_of_pages, collection_id, created_at, updated_at FROM chapter WHERE collection_id=$1", collectionID)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, fmt.Errorf("Error getting this collection's chapters")
// 	}

// 	var chapters []Chapter
// 	for rows.Next() {
// 		chapter := Chapter{}
// 		err = rows.Scan(&chapter.ID, &chapter.Title, &chapter.Description, &chapter.NumberOfPages, &chapter.ChapterNumber, &chapter.CollectionID, &chapter.CreatedAt, &chapter.UpdatedAt)
		
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
// 		chapters = append(chapters, chapter)
// 	}

// 	return chapters, nil
// }

func (i *Chapter) Create(conn *pgx.Conn) error {
	i.Title = strings.Trim(i.Title, " ")
	if len(i.Title) < 1 {
		return fmt.Errorf("Title must not be empty.")
	}

	now := time.Now()
	row := conn.QueryRow(context.Background(), "INSERT INTO chapter (title, description, number_of_pages, chapter_number, collection_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", i.Title, i.Description, i.NumberOfPages, i.ChapterNumber, i.CollectionID, now, now)

	err := row.Scan(&i.ID)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("There was an error creating the chapter")
	}

	return nil
}