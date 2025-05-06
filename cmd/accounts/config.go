package accounts

type Config struct {
	IsDev bool

	CockroachDB struct {
		DSN string
	}
}
