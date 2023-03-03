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

// SaveToList return (true, "", nil) if it is successfully added,
// return false if something went wrong,
// such as 1. internal error - return (false, empty string, error)
// 2. overlap conflict - return (false, message about conflict, nil)
func (lr *ListRepo) SaveToList(ctx context.Context, subnet, color string) (bool, string, error) {
	if ok, message, err := lr.IterateSubnets(ctx, subnet, color); err != nil {
		return false, "", fmt.Errorf("repo - SaveToList - lr.IterateSubnets: %w", err)
	} else if !ok {
		return ok, message, nil
	}

	sql, args, err := lr.Postgres.Builder.
		Insert("lists").
		Columns("subnet, list_type").
		Values(subnet, color).
		ToSql()
	if err != nil {
		return false, "", fmt.Errorf("repo - SaveToList - lr.Builder: %w", err)
	}

	_, err = lr.Postgres.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return false, "", fmt.Errorf("repo - SaveToList - lr.Pool.Exec: %w", err)
	}
	return true, "", nil
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

// IterateSubnets checks if IP address ranges overlap.
// Returns false and message if there is an overlap conflict, true and empty string - if not.
func (lr *ListRepo) IterateSubnets(ctx context.Context, subnetA, color string) (bool, string, error) {
	rows, err := lr.Postgres.Pool.Query(ctx, "SELECT subnet, list_type FROM lists")
	if err != nil {
		return false, "", fmt.Errorf("lr.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row struct {
			subnetB string
			color   string
		}
		var temp interface{}
		if er := rows.Scan(&temp, &row.color); er != nil {
			return false, "", fmt.Errorf("rows.Scan: %w", er)
		}

		row.subnetB = fmt.Sprintf("%v", temp)

		// if given subnet A already in list
		if ipaddr.NewIPAddressString(row.subnetB).GetAddress().
			Equal(ipaddr.NewIPAddressString(subnetA).GetAddress()) {
			// in case we try to add given subnet A in the different color list
			if row.color != color {
				return false,
					fmt.Sprintf("list conflict: can't add given subnet %v in %vlist because it is already in %vlist",
						subnetA, color, row.color), nil
			}
			// in case we try to add given subnet A in the same color list
			return false, fmt.Sprintf("given subnet %v already in %slist",
				subnetA, color), nil
		}

		// if given subnet A already in list as a part of bigger subnet B
		if ipaddr.NewIPAddressString(row.subnetB).GetAddress().
			Contains(ipaddr.NewIPAddressString(subnetA).GetAddress()) {
			// in case we try to add given subnet A to the different list
			if row.color != color {
				return false,
					fmt.Sprintf("lists conflict: subnet %v in %slist already include given subnet %v",
						row.subnetB, row.color, subnetA), nil
			}
			// in case we try to add given subnet A in the same subnet B color list
			return false, fmt.Sprintf("subnet %v already in %slist because it is included in subnet %v",
				subnetA, color, row.subnetB), nil
		} else if ipaddr.NewIPAddressString(subnetA).GetAddress(). // if given subnet A already include smaller subnet B in list
										Contains(ipaddr.NewIPAddressString(row.subnetB).GetAddress()) {
			// in case we try to add given subnet A to the different list
			if row.color != color {
				return false,
					fmt.Sprintf("lists conflict: given subnet %v already include subnet %v in different %slist",
						subnetA, row.subnetB, row.color), nil
			}
			// in case we try to add given subnet A in the same subnet B color list
			_, er := lr.Postgres.Pool.Exec(ctx, `delete from lists where subnet = $1`, row.subnetB)
			if er != nil {
				return false, "", fmt.Errorf("AcontainingB - lr.Pool.Exec - delete: %w", er)
			}
			return true, "", nil
		}
	}
	return true, "", nil
}

func (lr *ListRepo) Up() error {
	ctx := context.Background()

	query :=
		"CREATE TABLE IF NOT EXISTS lists (subnet CIDR PRIMARY KEY, list_type TEXT NOT NULL)"
	_, err := lr.Postgres.Pool.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

// Drop attaches the provider and drop the table
func (lr *ListRepo) Drop() error {
	ctx := context.Background()

	query := "DROP TABLE IF EXISTS lists"
	_, err := lr.Postgres.Pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
