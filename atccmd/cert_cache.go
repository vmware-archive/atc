package atccmd

import (
	"context"
	"database/sql"

	"github.com/concourse/atc/db"
	"golang.org/x/crypto/acme/autocert"
)

type dbCache struct {
	get, put, delete *sql.Stmt
	es               db.EncryptionStrategy
}

func newDbCache(conn db.Conn) (autocert.Cache, error) {
	c := new(dbCache)
	c.es = conn.EncryptionStrategy()
	var err error
	c.get, err = conn.Prepare("SELECT data, nonce FROM cert_cache WHERE key = $1")
	if err != nil {
		return nil, err
	}
	c.put, err = conn.Prepare("INSERT INTO cert_cache (key, data, nonce) VALUES ($1, $2, $3)")
	if err != nil {
		return nil, err
	}
	c.delete, err = conn.Prepare("DELETE FROM cert_cache WHERE key = $1")
	return c, err
}

func (c *dbCache) Get(ctx context.Context, key string) ([]byte, error) {
	var ciphertext string
	var nonce sql.NullString
	err := c.get.QueryRowContext(ctx, key).Scan(&ciphertext, &nonce)
	if err == sql.ErrNoRows {
		err = autocert.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	var noncense *string
	if nonce.Valid {
		noncense = &nonce.String
	}
	return c.es.Decrypt(ciphertext, noncense)
}

func (c *dbCache) Put(ctx context.Context, key string, data []byte) error {
	ciphertext, nonce, err := c.es.Encrypt(data)
	if err != nil {
		return err
	}
	_, err = c.put.ExecContext(ctx, key, ciphertext, nonce)
	return err
}

func (c *dbCache) Delete(ctx context.Context, key string) error {
	_, err := c.delete.ExecContext(ctx, key)
	return err
}
