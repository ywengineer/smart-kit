package redisx

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"

	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/gookit/goutil/errorx"
	"github.com/redis/go-redis/v9"
)

// NewRedis redisURL format redis://user:password@host:port/?db=0&node=host:port&node=host:port
//
// - when redis running in cluster mode, add query param node.
// - when redis running in simple mode, add query param db.
func NewRedis(redisURL string) redis.UniversalClient {
	u, err := url.Parse(redisURL)
	if err != nil {
		panic(errorx.With(err, "parse redis url failed"))
		return nil
	}
	query := u.Query()
	opt := &redis.UniversalOptions{
		Addrs: []string{u.Host},
		DB:    0,
	}
	// multi nodes
	if len(query["node"]) > 0 {
		opt.Addrs = append(opt.Addrs, query["node"]...)
	}
	// use auth
	if u.User != nil {
		opt.Username = u.User.Username()
		opt.Password, _ = u.User.Password()
	}
	if u.Scheme == "rediss" {
		opt.TLSConfig = &tls.Config{
			ServerName: u.Hostname(),
			MinVersion: tls.VersionTLS12,
		}
	}
	//
	opt, err = setupConnParams(u, opt)
	if err != nil {
		panic(errorx.With(err, "setup redis connection params failed"))
		return nil
	}
	//
	opt.MinIdleConns = utilk.Max(1, opt.MinIdleConns)
	//
	return redis.NewUniversalClient(opt)
}

// setupConnParams converts query parameters in u to option value in o.
func setupConnParams(u *url.URL, o *redis.UniversalOptions) (*redis.UniversalOptions, error) {
	q := queryOptions{q: u.Query()}

	o.DB = q.int("db")
	o.Protocol = q.int("protocol")
	o.ClientName = q.string("client_name")
	o.MaxRetries = q.int("max_retries")
	o.MinRetryBackoff = q.duration("min_retry_backoff")
	o.MaxRetryBackoff = q.duration("max_retry_backoff")
	o.DialTimeout = q.duration("dial_timeout")
	o.ReadTimeout = q.duration("read_timeout")
	o.WriteTimeout = q.duration("write_timeout")
	o.PoolFIFO = q.bool("pool_fifo")
	o.PoolSize = q.int("pool_size")
	o.PoolTimeout = q.duration("pool_timeout")
	o.MinIdleConns = q.int("min_idle_conns")
	o.MaxIdleConns = q.int("max_idle_conns")
	o.MaxActiveConns = q.int("max_active_conns")
	if q.has("conn_max_idle_time") {
		o.ConnMaxIdleTime = q.duration("conn_max_idle_time")
	} else {
		o.ConnMaxIdleTime = q.duration("idle_timeout")
	}
	if q.has("conn_max_lifetime") {
		o.ConnMaxLifetime = q.duration("conn_max_lifetime")
	} else {
		o.ConnMaxLifetime = q.duration("max_conn_age")
	}
	if q.err != nil {
		return nil, q.err
	}
	if o.TLSConfig != nil && q.has("skip_verify") {
		o.TLSConfig.InsecureSkipVerify = q.bool("skip_verify")
	}

	// any parameters left?
	if r := q.remaining(); len(r) > 0 {
		return nil, fmt.Errorf("redis: unexpected option: %s", strings.Join(r, ", "))
	}

	return o, nil
}
