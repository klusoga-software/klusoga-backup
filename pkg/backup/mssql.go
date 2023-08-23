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
		p := path.Join(m.parameters.Path, d+"-"+time.Now().Format(time.DateOnly)+".bak")
		sqlbackup := fmt.Sprintf("BACKUP DATABASE [%s] TO DISK = '%s';", d, p)
		_, err = db.Exec(sqlbackup)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		paths = append(paths, p)
	}

	return paths, nil
}
