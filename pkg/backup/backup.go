package backup

type Target interface {
	Backup() ([]string, error)
}
