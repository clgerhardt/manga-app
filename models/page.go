package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

// Page struct defines the properties of a page object
type Page struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"_"`
	UpdatedAt   time.Time `json:"_"`
	Title       string    `json:"title"`
	Page        []byte    `json:"page"`
	OrderNumber int64     `json:"order_number"`
	ChapterID   uuid.UUID `json:"chapter"`
}

// Create page method
func (i *Page) Create(conn *pgx.Conn) error {
	i.Title = strings.Trim(i.Title, " ")
	if len(i.Title) < 1 {
		return fmt.Errorf("Title must not be empty")
	}

	if i.OrderNumber == -1 {
		return fmt.Errorf("Order must not be empty")
	}

	now := time.Now()
	row := conn.QueryRow(context.Background(), "INSERT INTO page (title, page, chapter_id, created_at, updated_at, order_number) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, order_number", i.Title, i.Page, i.ChapterID, now, now, i.OrderNumber)

	err := row.Scan(&i.ID, &i.OrderNumber)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("There was an error creating the page")
	}

	return nil
}
