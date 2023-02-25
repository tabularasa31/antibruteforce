package repo

import (
	"context"
	"fmt"
	"net"

	"github.com/seancfoley/ipaddress-go/ipaddr"
	"github.com/tabularasa31/antibruteforce/pkg/postgres"
)

// ListRepo -.
type ListRepo struct {
	*postgres.Postgres
}

func NewListRepo(pg *postgres.Postgres) *ListRepo {
	return &ListRepo{pg}
}

func (lr *ListRepo) SaveToList(ctx context.Context, subnet, color string) (string, error) {
	if message, err := lr.iterateSubnets(ctx, subnet, color); err != nil {
		return "", fmt.Errorf("repo - SaveToList - lr.iterateSubnets: %w", err)
	} else if message != "" {
		return message, nil
	}

	sql, args, err := lr.Postgres.Builder.
		Insert("lists").
		Columns("subnet, list_type").
		Values(subnet, color).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("repo - SaveToList - lr.Builder: %w", err)
	}

	_, err = lr.Postgres.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return "", fmt.Errorf("repo - SaveToList - lr.Pool.Exec: %w", err)
	}
	return "", nil
}

// DeleteFromList subnet from lists -.
func (lr *ListRepo) DeleteFromList(ctx context.Context, subnet, color string) error {
	row := lr.Postgres.Pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM lists WHERE subnet = $1 AND list_type = $2)`, subnet, color)
	var found string
	if err := row.Scan(&found); err != nil {
		return err
	}
	if found == "" {
		return fmt.Errorf("there is no subnet %s in %slist", subnet, color)
	}

	_, err := lr.Postgres.Pool.Exec(ctx, `delete from lists where subnet = $1`, subnet)
	if err != nil {
		return err
	}
	return nil
}

// CheckColor checks if a subnet exists in the specified list.
func (lr *ListRepo) CheckColor(ctx context.Context, subnet string) (string, error) {
	sql, args, err := lr.Postgres.Builder.Select("list_type").
		From("lists").
		Where("subnet=?", subnet).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("repo - CheckColor - lr.Builder: %w", err)
	}

	row := lr.Postgres.Pool.QueryRow(ctx, sql, args...)
	if err != nil {
		return "", fmt.Errorf("repo - CheckColor - lr.Pool.QueryRow: %w", err)
	}

	color := ""
	if er := row.Scan(&color); er != nil {
		return "", fmt.Errorf("repo - CheckColor - row.Scan: %w", er)
	}

	return color, nil
}

// SearchIPInList check if it is a given IP in the lists, returns color.
func (lr *ListRepo) SearchIPInList(ctx context.Context, ip net.IP) string {
	row := lr.Postgres.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT list_type FROM lists WHERE subnet >>= $1)`, ip)
	var found string
	if err := row.Scan(&found); err != nil {
		return ""
	}
	return found
}

// iterateSubnets checks if IP address ranges overlap.
// Returns message if there is an overlap conflict, empty string - if not.
func (lr *ListRepo) iterateSubnets(ctx context.Context, subnetA, color string) (string, error) {
	rows, err := lr.Postgres.Pool.Query(ctx, "SELECT subnet, list_type FROM lists")
	if err != nil {
		return "", fmt.Errorf("lr.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row struct {
			subnetB string
			color   string
		}
		var temp interface{}
		if er := rows.Scan(&temp, &row.color); er != nil {
			return "", fmt.Errorf("rows.Scan: %w", er)
		}

		row.subnetB = fmt.Sprintf("%v", temp)

		if ipaddr.NewIPAddressString(row.subnetB).GetAddress().
			Contains(ipaddr.NewIPAddressString(subnetA).GetAddress()) {
			if row.color != color {
				return "",
					fmt.Errorf("lists conflict: subnet %v in %slist already include given subnet %v",
						row.subnetB, row.color, subnetA)
			}
			return fmt.Sprintf("subnet %v already in %slist because it is included in subnet %v",
				subnetA, color, row.subnetB), nil
		} else if ipaddr.NewIPAddressString(subnetA).GetAddress().
			Contains(ipaddr.NewIPAddressString(row.subnetB).GetAddress()) {
			if row.color != color {
				return "",
					fmt.Errorf("lists conflict: given subnet %v already include subnet %v in different %slist",
						subnetA, row.subnetB, row.color)
			}
			_, er := lr.Postgres.Pool.Exec(ctx, `delete from lists where subnet = $1`, row.subnetB)
			if er != nil {
				return "", fmt.Errorf("AcontainingB - lr.Pool.Exec - delete: %w", er)
			}
			return "", nil
		}
	}
	return "", nil
}
