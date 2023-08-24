package backup

import (
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log/slog"
	"net/url"
	"path"
	"time"
)

type mssqlTarget struct {
	parameters MssqlParameters
}

type MssqlParameters struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Path      string
	Databases []string
}

func NewMssqlTarget(parameters MssqlParameters) Target {
	return &mssqlTarget{
		parameters: parameters,
	}
}

func (m *mssqlTarget) Backup() ([]string, error) {
	slog.Info("Run Mssql Backup")
	var paths []string

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(m.parameters.Username, m.parameters.Password),
		Host:   fmt.Sprintf("%s:%d", m.parameters.Host, m.parameters.Port),
	}

	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	for _, d := range m.parameters.Databases {
		p := path.Join(m.parameters.Path, d+"-"+time.Now().UTC().Format("2006-01-02-15-04-05")+".bak")
		sqlbackup := fmt.Sprintf("BACKUP DATABASE [%s] TO DISK = '%s';", d, p)
		_, err = db.Exec(sqlbackup)
		if err != nil {
			slog.Error("Error while backup database", "error", err.Error(), "database", d)
			continue
		}

		paths = append(paths, p)
	}

	slog.Info("Backups finished")

	return paths, nil
}
