package store

// ChartStore -
type ChartStore struct {
	db *DB
}

// NewChartStore -
func NewChartStore(db *DB) *ChartStore {
	return &ChartStore{
		db: db,
	}
}

// UserTotalChart -
func (s *ChartStore) UserTotalChart() error {
	return nil
}
