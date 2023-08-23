package backup

type Target interface {
	Backup() error
}
