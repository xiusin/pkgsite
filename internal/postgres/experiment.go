// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package postgres

import (
	"context"
	"database/sql"

	"golang.org/x/discovery/internal"
	"golang.org/x/discovery/internal/derrors"
	"golang.org/x/xerrors"
)

// GetActiveExperiments fetches all experiments where rollout is greater than 0.
func (db *DB) GetActiveExperiments(ctx context.Context) (_ []*internal.Experiment, err error) {
	defer derrors.Wrap(&err, "DB.GetActiveExperiments(ctx)")

	query := "SELECT name, rollout, description FROM experiments WHERE rollout > 0;"
	var experiments []*internal.Experiment
	err = db.db.RunQuery(ctx, query, func(rows *sql.Rows) error {
		var e internal.Experiment
		if err := rows.Scan(&e.Name, &e.Rollout, &e.Description); err != nil {
			return err
		}
		experiments = append(experiments, &e)
		return nil
	})
	return experiments, nil
}

// insertExperiment inserts a row into the experiments table.
func (db *DB) insertExperiment(ctx context.Context, e *internal.Experiment) (err error) {
	defer derrors.Wrap(&err, "DB.insertExperiment(ctx, %v)", e)
	if e.Name == "" || e.Description == "" {
		return xerrors.Errorf("neither name nor description can be empty: %w", derrors.InvalidArgument)
	}

	_, err = db.db.Exec(ctx,
		`INSERT INTO experiments
		(name, rollout, description) VALUES ($1, $2, $3);`,
		e.Name, e.Rollout, e.Description)
	return err
}

// updateExperiment updates the specified experiment with the provided rollout value.
func (db *DB) updateExperiment(ctx context.Context, e *internal.Experiment) (err error) {
	defer derrors.Wrap(&err, "DB.updateExperimentRollout(ctx, %v)", e)
	if e.Name == "" || e.Description == "" {
		return xerrors.Errorf("neither name nor description can be empty: %w", derrors.InvalidArgument)
	}

	query := `UPDATE experiments
		SET rollout = $2, description = $3
		WHERE name = $1;`
	_, err = db.db.Exec(ctx, query, e.Name, e.Rollout, e.Description)
	return err
}
