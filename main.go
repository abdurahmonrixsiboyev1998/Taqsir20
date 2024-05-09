package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

const (
	dbConnectionString = "postgres://postgres:14022014@localhost/demo?sslmode=disable"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var wg sync.WaitGroup
	num, ctx := errgroup.WithContext(ctx)

	
	wg.Add(1)
	num.Go(func() error {
		defer wg.Done()

		_, err := db.ExecContext(ctx, "INSERT INTO large_dataset(generated) VALUES ($1)", 100)
		if err != nil {
			return err
		}

		fmt.Println("Insert operatsiyasi muvaffaqiyatli amalga oshirildi.")
		return nil
	})

	wg.Add(1)
	num.Go(func() error {
		defer wg.Done()

		r, err := db.QueryContext(ctx, "SELECT id, generated FROM large_dataset")
		if err != nil {
			return err
		}
		defer r.Close()

		fmt.Println("Select operation results:")
		for r.Next() {
			var id, generated int
			if err := r.Scan(&id, &generated); err != nil {
				return err
			}
			fmt.Printf("ID: %d, Generated: %d\n", id, generated)
		}

		return r.Err()
	})

	wg.Add(1)
	num.Go(func() error {
		defer wg.Done()

		_, err := db.ExecContext(ctx, "UPDATE large_dataset SET generated = $1 WHERE id = $2", 200, 1)
		if err != nil {
			return err
		}

		fmt.Println("Update the operation was successful.")
		return nil
	})


	wg.Wait()
	if err := num.Wait(); err != nil {
		fmt.Println("Error:", err)
	}
}
