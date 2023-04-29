package mysql

import (
	"errors"

	"github.com/IceFireDB/IceFireDB-SQLProxy/pkg/mysql/client"
	"github.com/IceFireDB/IceFireDB-SQLProxy/pkg/mysql/mysql"
	"github.com/IceFireDB/IceFireDB-SQLProxy/pkg/mysql/server"
)

type Handle struct {
	conn *client.Conn
}

func (h *Handle) CloseConn(c *server.Conn) error {
	if c.IsInTransaction() || c.IsAutoCommit() {
		return h.conn.Commit()
	}
	return nil
}

func (h *Handle) UseDB(c *server.Conn, dbName string) error {
	return h.conn.UseDB(dbName)
}

func (h *Handle) HandleQuery(c *server.Conn, query string) (res *mysql.Result, err error) {
	res, err = h.conn.Execute(query)
	if err == nil {
		broadcast(query)
	}
	return
}

func (h *Handle) HandleFieldList(c *server.Conn, table string, fieldWildcard string) ([]*mysql.Field, error) {
	return h.conn.FieldList(table, fieldWildcard)
}

func (h *Handle) HandleStmtPrepare(c *server.Conn, query string) (int, int, interface{}, error) {
	stmt, err := h.conn.Prepare(query)
	if err != nil {
		return 0, 0, nil, err
	}
	return stmt.ParamNum(), stmt.ColumnNum(), stmt, nil
}

func (h *Handle) HandleStmtExecute(c *server.Conn, context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	stmt, ok := context.(*client.Stmt)
	if !ok {
		return nil, errors.New("other error")
	}
	res, err := stmt.Execute(args...)
	if err == nil {
		broadcast(query)
	}
	return res, err
}

func (h *Handle) HandleStmtClose(c *server.Conn, context interface{}) error {
	stmt, ok := context.(*client.Stmt)
	if !ok {
		return errors.New("other error")
	}
	return stmt.Close()
}

// sq
func (h *Handle) HandleOtherCommand(c *server.Conn, cmd byte, data []byte) error {
	return errors.New("command %d is not supported now" + string(cmd))
}
