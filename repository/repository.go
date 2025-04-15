package repository

type Repository interface {
	Close() error
  Ping() error
}
