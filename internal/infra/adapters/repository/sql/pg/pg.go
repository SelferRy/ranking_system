package pg

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/repository"
	"github.com/SelferRy/ranking_system/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type Repo struct {
	*BannerRepo
}

//func (s *Repo) Connect(ctx context.Context) error {
//	return nil
//}
//
//func (s *Repo) Close(ctx context.Context) error {
//	return nil
//}

func New(ctx context.Context, logger logger.Logger, conf repository.Config) (*Repo, error) {
	if conf.Driver != "postgres" {
		log.Fatal("Unimplemented db driver.")
	}
	dns := func() string {
		//var b strings.Builder
		const dbProtocol = "postgres://"
		const stdDbSetup = "?sslmode=disable"
		arr := []string{
			dbProtocol,
			conf.User, ":", conf.Password,
			"@", conf.Host, ":", conf.Port,
			"/", conf.Database, stdDbSetup}
		return strings.Join(arr, "")
		//for s := range arr {
		//
		//}
		//b.WriteString("postgres://")
		//b.WriteString(conf.User)
		//b.WriteString(":")
		//b.WriteString(conf.Password)
		//b.WriteString("@")
		//b.WriteString(conf.Host)
		//b.WriteString(":")
		//b.
		//for i := 3; i >= 1; i-- {
		//	fmt.Fprintf(&b, "%d...", i)
		//}
		//
		//fmt.Println(b.String())
	}()
	// make dns. Extract into func
	//dns := "postgres://" + conf.User + ":" + conf.Password + "@" + conf.Host + ":" + conf.Port + "/" + conf.Database
	//dns = dns + "?sslmode=disable"
	//dns := 'postgres://ad-rotator:ad-rotator@localhost:5432/ad-rotator?sslmode=disable'
	driverName := "pgx"
	db, err := sqlx.ConnectContext(ctx, driverName, dns)
	if err != nil {
		return nil, fmt.Errorf("sqlx.ConnectContext(...) error: %w", err)
	}
	db.SetMaxOpenConns(conf.MaxConn)
	db.SetMaxIdleConns(conf.MaxIdleConn)
	db.SetConnMaxLifetime(conf.MaxLifetime)
	return &Repo{
		BannerRepo: NewBannerRepo(db),
	}, nil
}
