package postgresrepo

import (
	"company-service/internal/company"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompanyRepo struct {
	db *pgxpool.Pool
}

func NewCompanyRepo(pool *pgxpool.Pool) *CompanyRepo {
	return &CompanyRepo{db: pool}
}

func (r *CompanyRepo) Create(ctx context.Context, c *company.Company) error {
	_, err := r.db.Exec(ctx, `INSERT INTO companies (id, name, description, amount_of_employees, registered, type)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		c.ID, c.Name, c.Description, c.AmountOfEmployees, c.Registered, c.Type)
	return err
}

func (r *CompanyRepo) Get(ctx context.Context, id uuid.UUID) (*company.Company, error) {
	row := r.db.QueryRow(ctx, `SELECT id, name, description, amount_of_employees, registered, type FROM companies WHERE id=$1`, id)
	var c company.Company
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.AmountOfEmployees, &c.Registered, &c.Type)
	return &c, err
}

func (r *CompanyRepo) Update(ctx context.Context, c *company.Company) error {
	_, err := r.db.Exec(ctx, `UPDATE companies SET name=$1, description=$2, amount_of_employees=$3, registered=$4, type=$5 WHERE id=$6`,
		c.Name, c.Description, c.AmountOfEmployees, c.Registered, c.Type, c.ID)
	return err
}

func (r *CompanyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM companies WHERE id=$1`, id)
	return err
}
