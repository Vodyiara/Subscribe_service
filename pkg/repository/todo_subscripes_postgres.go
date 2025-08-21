package repository

import (
	"Test_project_Effective_Mobile"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
	"time"
)

const layout = "2006-01-02"

type TodoSubscription struct {
	db *sqlx.DB
}

//TODO: Узнать в каком формате приходят даты

func NewTodoSubscription(db *sqlx.DB) *TodoSubscription {
	return &TodoSubscription{db: db}
}

func (r *TodoSubscription) Create(item Test_project_Effective_Mobile.Subscripe) (int, error) {
	var currentDB string
	if err := r.db.Get(&currentDB, "SELECT current_database()"); err != nil {
		log.Fatal(err)
	}

	startDate, err := ParseDate(item.Start_date)
	if err != nil {
		return 0, fmt.Errorf("invalid start_date: %w", err)
	}

	var endDate *time.Time
	if item.End_date != nil {
		parsedEnd, err := ParseDate(*item.End_date)
		if err != nil {
			return 0, fmt.Errorf("invalid end_date: %w", err)
		}
		endDate = &parsedEnd
	}

	var id int
	query := fmt.Sprintf(
		"INSERT INTO %s (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		subscriptionsTable,
	)
	row := r.db.QueryRow(query, item.Service_name, item.Price, item.User_id, startDate, endDate)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert subscription: %w", err)
	}

	// логирование удалено из репозитория
	return id, nil
}

func (r *TodoSubscription) GetAllSubs() ([]Test_project_Effective_Mobile.Subscripe, error) {
	var lists []Test_project_Effective_Mobile.Subscripe
	query := fmt.Sprintf(`
    SELECT tl.id, tl.service_name, tl.price, tl.user_id, tl.start_date, tl.end_date
    FROM %s tl
    ORDER BY tl.id
`, subscriptionsTable)

	type dbRow struct {
		Id           int        `db:"id"`
		Service_name string     `db:"service_name"`
		Price        int        `db:"price"`
		User_id      uuid.UUID  `db:"user_id"`
		Start_date   time.Time  `db:"start_date"`
		End_date     *time.Time `db:"end_date"`
	}

	var rows []dbRow
	if err := r.db.Select(&rows, query); err != nil {
		return nil, err
	}

	for _, row := range rows {
		s := Test_project_Effective_Mobile.Subscripe{
			Id:           row.Id,
			Service_name: row.Service_name,
			Price:        row.Price,
			User_id:      row.User_id,
			Start_date:   row.Start_date.Format("01-2006"),
		}
		if row.End_date != nil {
			formatted := row.End_date.Format("01-2006")
			s.End_date = &formatted
		}
		lists = append(lists, s)
	}

	return lists, nil
}

func (r *TodoSubscription) Update(subId int, item Test_project_Effective_Mobile.UpdateSubInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if item.Service_name != nil {
		setValues = append(setValues, fmt.Sprintf("service_name=$%d", argId))
		args = append(args, *item.Service_name)
		argId++
	}

	if item.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *item.Price)
		argId++
	}

	if item.Start_date != nil {
		parsedStart, err := ParseDate(*item.Start_date)
		if err != nil {
			return fmt.Errorf("invalid start_date: %w", err)
		}
		setValues = append(setValues, fmt.Sprintf("start_date=$%d", argId))
		args = append(args, parsedStart)
		argId++
	}

	if item.End_date != nil {
		parsedEnd, err := ParseDate(*item.End_date)
		if err != nil {
			return fmt.Errorf("invalid end_date: %w", err)
		}
		setValues = append(setValues, fmt.Sprintf("end_date=$%d", argId))
		args = append(args, parsedEnd)
		argId++
	}

	if item.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *item.Done)
		argId++
	}

	if len(setValues) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// последний аргумент - id
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", subscriptionsTable, strings.Join(setValues, ", "), argId)
	args = append(args, subId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %v", args)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoSubscription) Delete(subId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", subscriptionsTable)
	_, err := r.db.Exec(query, subId)
	return err
}

func (r *TodoSubscription) GetByIdSub(subId int) (Test_project_Effective_Mobile.Subscripe, error) {
	var subs Test_project_Effective_Mobile.Subscripe

	type dbRow struct {
		Id           int        `db:"id"`
		Service_name string     `db:"service_name"`
		Price        int        `db:"price"`
		User_id      uuid.UUID  `db:"user_id"`
		Start_date   time.Time  `db:"start_date"`
		End_date     *time.Time `db:"end_date"`
	}

	var row dbRow
	query := fmt.Sprintf(`
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM %s 
		WHERE id = $1
	`, subscriptionsTable)

	if err := r.db.Get(&row, query, subId); err != nil {
		return subs, fmt.Errorf("failed to get subscription by id %d: %w", subId, err)
	}

	subs = Test_project_Effective_Mobile.Subscripe{
		Id:           row.Id,
		Service_name: row.Service_name,
		Price:        row.Price,
		User_id:      row.User_id,
		Start_date:   row.Start_date.Format("01-2006"),
	}
	if row.End_date != nil {
		formatted := row.End_date.Format("01-2006")
		subs.End_date = &formatted
	}

	return subs, nil
}

func (r *TodoSubscription) GetTotalPrice(userId uuid.UUID, serviceName, startDate, endDate string) (int, error) {
	var total sql.NullInt64

	start, err := time.Parse("01-2006", startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date format: %w", err)
	}

	end, err := time.Parse("01-2006", endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid end date format: %w", err)
	}

	end = time.Date(end.Year(), end.Month()+1, 0, 23, 59, 59, 0, time.UTC)

	query := fmt.Sprintf(`
		SELECT COALESCE(SUM(price), 0)
		FROM %s
		WHERE start_date >= $1 AND (end_date IS NULL OR end_date <= $2)
	`, subscriptionsTable)

	args := []interface{}{start, end}
	argIdx := 3

	if userId != uuid.Nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, userId)
		argIdx++
	}

	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argIdx)
		args = append(args, serviceName)
	}

	err = r.db.Get(&total, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get total price: %w", err)
	}

	if total.Valid {
		return int(total.Int64), nil
	}
	return 0, nil
}

func ParseDate(input string) (time.Time, error) {
	t, err := time.Parse("01-2006", input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected MM-YYYY: %w", err)
	}
	return t, nil
}
