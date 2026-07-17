package db

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/joaovictornsv/books-cli/internal/models"
)

var bookCSVColumns = []string{
	"id", "title", "author", "category", "status",
	"priority_to_buy", "eligible_to_donate", "donated",
	"notes", "description", "added_at", "started_at", "finished_at",
}

type ImportResult struct {
	Created int
	Updated int
	Total   int
	Errors  []string
	DryRun  bool
}

type exportDocument struct {
	Books []models.Book `json:"books"`
	Total int           `json:"total"`
}

func (r *Repository) ListAll(ctx context.Context, includeArchived bool) ([]models.Book, error) {
	result, err := r.List(ctx, ListFilter{IncludeArchived: includeArchived})
	if err != nil {
		return nil, err
	}
	return result.Books, nil
}

func MarshalBooksJSON(books []models.Book) ([]byte, error) {
	if books == nil {
		books = []models.Book{}
	}
	doc := exportDocument{Books: books, Total: len(books)}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal books json: %w", err)
	}
	return append(data, '\n'), nil
}

func UnmarshalBooksJSON(data []byte) ([]models.Book, error) {
	var doc exportDocument
	if err := json.Unmarshal(data, &doc); err == nil && doc.Books != nil {
		return doc.Books, nil
	}

	var books []models.Book
	if err := json.Unmarshal(data, &books); err != nil {
		return nil, fmt.Errorf("parse books json: %w", err)
	}
	return books, nil
}

func WriteBooksCSV(w io.Writer, books []models.Book) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(bookCSVColumns); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}
	for _, book := range books {
		if err := cw.Write(bookToCSVRecord(book)); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		return fmt.Errorf("flush csv: %w", err)
	}
	return nil
}

func ReadBooksCSV(r io.Reader) ([]models.Book, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1

	records, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}
	if len(records) == 0 {
		return []models.Book{}, nil
	}

	header := records[0]
	colIndex, err := csvColumnIndex(header)
	if err != nil {
		return nil, err
	}

	books := make([]models.Book, 0, len(records)-1)
	for i, record := range records[1:] {
		book, err := csvRecordToBook(record, colIndex)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i+2, err)
		}
		books = append(books, book)
	}
	return books, nil
}

func csvColumnIndex(header []string) (map[string]int, error) {
	index := make(map[string]int, len(header))
	for i, name := range header {
		index[strings.TrimSpace(name)] = i
	}
	for _, col := range bookCSVColumns {
		if _, ok := index[col]; !ok {
			return nil, fmt.Errorf("csv missing required column %q", col)
		}
	}
	return index, nil
}

func bookToCSVRecord(book models.Book) []string {
	return []string{
		strconv.FormatInt(book.ID, 10),
		book.Title,
		ptrToString(book.Author),
		categoryToString(book.Category),
		book.Status.String(),
		strconv.Itoa(book.PriorityToBuy),
		strconv.Itoa(book.EligibleToDonate),
		strconv.Itoa(book.Donated),
		ptrToString(book.Notes),
		ptrToString(book.Description),
		book.AddedAt,
		ptrToString(book.StartedAt),
		ptrToString(book.FinishedAt),
	}
}

func csvRecordToBook(record []string, colIndex map[string]int) (models.Book, error) {
	get := func(col string) string {
		i := colIndex[col]
		if i >= len(record) {
			return ""
		}
		return record[i]
	}

	idRaw := strings.TrimSpace(get("id"))
	var id int64
	var err error
	if idRaw != "" {
		id, err = strconv.ParseInt(idRaw, 10, 64)
		if err != nil {
			return models.Book{}, fmt.Errorf("invalid id: %w", err)
		}
	}

	status, err := models.ParseStatus(get("status"))
	if err != nil {
		return models.Book{}, err
	}

	priority, err := strconv.Atoi(strings.TrimSpace(get("priority_to_buy")))
	if err != nil {
		return models.Book{}, fmt.Errorf("invalid priority_to_buy: %w", err)
	}
	eligible, err := strconv.Atoi(strings.TrimSpace(get("eligible_to_donate")))
	if err != nil {
		return models.Book{}, fmt.Errorf("invalid eligible_to_donate: %w", err)
	}
	donated, err := strconv.Atoi(strings.TrimSpace(get("donated")))
	if err != nil {
		return models.Book{}, fmt.Errorf("invalid donated: %w", err)
	}

	book := models.Book{
		ID:               id,
		Title:            get("title"),
		Status:           status,
		PriorityToBuy:    priority,
		EligibleToDonate: eligible,
		Donated:          donated,
		AddedAt:          get("added_at"),
	}
	book.Author = optionalString(get("author"))
	book.Notes = optionalString(get("notes"))
	book.Description = optionalString(get("description"))
	book.StartedAt = optionalString(get("started_at"))
	book.FinishedAt = optionalString(get("finished_at"))

	if cat := strings.TrimSpace(get("category")); cat != "" {
		category, err := models.ParseCategory(cat)
		if err != nil {
			return models.Book{}, err
		}
		book.Category = &category
	}

	return book, nil
}

func optionalString(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
}

func ptrToString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func categoryToString(v *models.Category) string {
	if v == nil {
		return ""
	}
	return v.String()
}

func (r *Repository) ImportBooks(ctx context.Context, books []models.Book, dryRun bool) (ImportResult, error) {
	result := ImportResult{Total: len(books), DryRun: dryRun}

	for i, book := range books {
		if err := book.Validate(); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", i+1, err))
		}
	}
	if len(result.Errors) > 0 {
		return result, nil
	}

	if dryRun {
		for _, book := range books {
			if book.ID > 0 {
				_, err := r.GetByID(ctx, book.ID)
				if err == nil {
					result.Updated++
					continue
				}
				if err != ErrNotFound {
					return ImportResult{}, err
				}
			}
			result.Created++
		}
		return result, nil
	}

	tx, err := r.db.sql.BeginTx(ctx, nil)
	if err != nil {
		return ImportResult{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, book := range books {
		if book.ID > 0 {
			_, err := getByIDWithQuerier(ctx, tx, book.ID)
			if err == nil {
				if err := replaceBookWithQuerier(ctx, tx, book.ID, book); err != nil {
					return ImportResult{}, err
				}
				result.Updated++
				continue
			}
			if !isNotFound(err) {
				return ImportResult{}, err
			}
			if _, err := createWithIDQuerier(ctx, tx, book); err != nil {
				return ImportResult{}, err
			}
			result.Created++
			continue
		}

		if _, err := createWithQuerier(ctx, tx, book); err != nil {
			return ImportResult{}, err
		}
		result.Created++
	}

	if err := tx.Commit(); err != nil {
		return ImportResult{}, fmt.Errorf("commit transaction: %w", err)
	}
	return result, nil
}

func isNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func replaceBookWithQuerier(ctx context.Context, q querier, id int64, book models.Book) error {
	book.ID = id
	if err := book.Validate(); err != nil {
		return err
	}
	_, err := q.ExecContext(ctx, `
		UPDATE books SET
			title = ?, author = ?, category = ?, status = ?, priority_to_buy = ?, eligible_to_donate = ?,
			donated = ?, notes = ?, description = ?, added_at = ?, started_at = ?, finished_at = ?
		WHERE id = ?`,
		book.Title,
		nullString(book.Author),
		nullCategory(book.Category),
		book.Status.String(),
		book.PriorityToBuy,
		book.EligibleToDonate,
		book.Donated,
		nullString(book.Notes),
		nullString(book.Description),
		book.AddedAt,
		nullString(book.StartedAt),
		nullString(book.FinishedAt),
		id,
	)
	if err != nil {
		return fmt.Errorf("replace book: %w", err)
	}
	return nil
}

func createWithIDQuerier(ctx context.Context, q querier, book models.Book) (models.Book, error) {
	if book.AddedAt == "" {
		book.AddedAt = models.NowTimestamp()
	}
	if err := book.Validate(); err != nil {
		return models.Book{}, err
	}
	_, err := q.ExecContext(ctx, `
		INSERT INTO books (
			id, title, author, category, status, priority_to_buy, eligible_to_donate, donated,
			notes, description, added_at, started_at, finished_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		book.ID,
		book.Title,
		nullString(book.Author),
		nullCategory(book.Category),
		book.Status.String(),
		book.PriorityToBuy,
		book.EligibleToDonate,
		book.Donated,
		nullString(book.Notes),
		nullString(book.Description),
		book.AddedAt,
		nullString(book.StartedAt),
		nullString(book.FinishedAt),
	)
	if err != nil {
		return models.Book{}, fmt.Errorf("insert book with id: %w", err)
	}
	return book, nil
}

func createWithQuerier(ctx context.Context, q querier, book models.Book) (models.Book, error) {
	if book.AddedAt == "" {
		book.AddedAt = models.NowTimestamp()
	}
	if err := book.Validate(); err != nil {
		return models.Book{}, err
	}
	res, err := q.ExecContext(ctx, `
		INSERT INTO books (
			title, author, category, status, priority_to_buy, eligible_to_donate, donated,
			notes, description, added_at, started_at, finished_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		book.Title,
		nullString(book.Author),
		nullCategory(book.Category),
		book.Status.String(),
		book.PriorityToBuy,
		book.EligibleToDonate,
		book.Donated,
		nullString(book.Notes),
		nullString(book.Description),
		book.AddedAt,
		nullString(book.StartedAt),
		nullString(book.FinishedAt),
	)
	if err != nil {
		return models.Book{}, fmt.Errorf("insert book: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Book{}, fmt.Errorf("last insert id: %w", err)
	}
	book.ID = id
	return book, nil
}
