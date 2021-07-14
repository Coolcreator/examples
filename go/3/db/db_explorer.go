package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Column struct {
	Name    string
	Type    string
	Primary bool
	IsNull  bool
}

type Table struct {
	Name       string
	PrimaryKey string
	Columns    []Column
}

func (t *Table) NewRow() []interface{} {
	return make([]interface{}, len(t.Columns))
}

func (t *Table) SetPointers(row []interface{}) []interface{} {
	pointers := make([]interface{}, len(t.Columns))
	for i := range row {
		pointers[i] = &row[i]
	}
	return pointers
}

type Handler struct {
	DB     *sql.DB
	Tables map[string]*Table
}

type Query struct {
	Method string
	Body   io.ReadCloser
	Table  string
	Id     int
	Offset int
	Limit  int
}

type Response struct {
	Data  interface{} `json:"response,omitempty"`
	Error string      `json:"error,omitempty"`
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB:     db,
		Tables: map[string]*Table{},
	}
}

func NewDbExplorer(db *sql.DB) (*Handler, error) {
	h, err := NewHandler(db).Init()
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *Handler) Init() (*Handler, error) {
	names, err := h.GetTableNames()
	if err != nil {
		return nil, err
	}

	for _, name := range names {
		columns, err := h.GetTableColumns(name)
		if err != nil {
			return nil, err
		}

		h.Tables[name] = &Table{Name: name, Columns: columns}
		for _, col := range h.Tables[name].Columns {
			if col.Primary == true {
				h.Tables[name].PrimaryKey = col.Name
			}
		}
	}

	return h, nil
}

func (h *Handler) GetTableNames() (tables []string, err error) {
	rows, err := h.DB.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var table string
	for rows.Next() {
		rows.Scan(&table)
		tables = append(tables, table)
	}

	return
}

func (h *Handler) GetTableColumns(table string) (columns []Column, err error) {
	rows, err := h.DB.Query("SHOW FULL COLUMNS FROM " + table)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var (
		Type    string
		Primary string
		Null    string
		Excess  interface{}
		column  Column = Column{}
	)

	for rows.Next() {

		err := rows.Scan(
			&column.Name,
			&Type,
			&Excess,
			&Null,
			&Primary,
			&Excess,
			&Excess,
			&Excess,
			&Excess,
		)

		if err != nil {
			return nil, err
		}

		if strings.Contains(Type, "int") {
			column.Type = "int"
		} else if strings.Contains(Type, "varchar") || strings.Contains(Type, "text") {
			column.Type = "string"
		}

		column.IsNull = Null == "YES"
		column.Primary = Primary == "PRI"
		columns = append(columns, column)
	}

	return
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query, err := h.ParseQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := h.Tables[query.Table]; !ok && query.Table != "All" {
		resp, err := json.Marshal(Response{nil, "unknown table"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(resp)
		return
	}

	table := h.Tables[query.Table]

	switch query.Method {
	case http.MethodGet:
		h.Get(w, table, query)
	case http.MethodPut:
		h.Put(w, table, query)
	case http.MethodPost:
		h.Post(w, table, query)
	case http.MethodDelete:
		h.Delete(w, table, query)
	}
}

func (h *Handler) ParseQuery(r *http.Request) (query *Query, err error) {
	if r.URL.Path == "/" && r.Method == http.MethodGet {
		return &Query{Method: http.MethodGet, Table: "All"}, nil
	}

	query = &Query{
		Method: r.Method,
		Body:   r.Body,
		Offset: 0,
		Limit:  5,
	}

	urlParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(urlParts) == 1 {
		query.Table = urlParts[0]
	} else if len(urlParts) == 2 {
		query.Table = urlParts[0]
		query.Id, err = strconv.Atoi(urlParts[1])
		if err != nil {
			return nil, err
		}
	}

	queryParts := r.URL.Query()

	if limit, ok := queryParts["limit"]; ok {
		query.Limit, err = strconv.Atoi(limit[0])
		if err != nil {
			query.Limit = 5
		}
	}

	if offset, ok := queryParts["offset"]; ok {
		query.Offset, err = strconv.Atoi(offset[0])
		if err != nil {
			query.Offset = 0
		}
	}

	return query, nil

}

func (h *Handler) Get(w http.ResponseWriter, table *Table, q *Query) {

	if q.Table == "All" {
		tables, err := h.GetTableNames()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(Response{map[string]interface{}{"tables": tables}, ""})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return
	}

	if q.Id == 0 {
		rows, err := h.DB.Query("SELECT * FROM "+q.Table+" LIMIT ? OFFSET ?", q.Limit, q.Offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var records []map[string]interface{}
		for rows.Next() {
			row := table.NewRow()
			pointers := table.SetPointers(row)
			if err := rows.Scan(pointers...); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var record map[string]interface{} = map[string]interface{}{}
			for idx, col := range table.Columns {
				value := row[idx]
				if bytes, ok := row[idx].([]byte); ok {
					value = string(bytes)
				}
				record[col.Name] = value
			}
			records = append(records, record)
		}

		resp := Response{map[string]interface{}{"records": records}, ""}
		j, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(j)
		return
	} else {
		rows := h.DB.QueryRow("SELECT * FROM "+table.Name+" WHERE "+table.PrimaryKey+" = ?", q.Id)
		row := table.NewRow()
		pointers := table.SetPointers(row)
		if err := rows.Scan(pointers...); err != nil {
			if err == sql.ErrNoRows {
				resp, err := json.Marshal(Response{Error: "record not found"})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNotFound)
				w.Write(resp)
				return
			}
			return
		}

		record := make(map[string]interface{})
		for idx, col := range table.Columns {
			value := row[idx]
			if bytes, ok := row[idx].([]byte); ok {
				value = string(bytes)
			}
			record[col.Name] = value
		}

		resp, err := json.Marshal(Response{map[string]interface{}{"record": record}, ""})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return
	}

}

func (h *Handler) Put(w http.ResponseWriter, table *Table, q *Query) {

	if q.Id != 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(q.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	body := make(map[string]interface{})
	json.Unmarshal(b, &body)

	names := []string{}
	values := []interface{}{}
	placeholders := []string{}

	for _, col := range table.Columns {
		if value, matched := body[col.Name]; matched {

			if col.Primary {
				continue
			}

			switch value.(type) {
			case float64:
				if col.Type == "string" {
					FierldError(col.Name, w)
					return
				}
			case string:
				if col.Type == "int" {
					FierldError(col.Name, w)
					return
				}
			}

			values = append(values, value)
		} else {
			if col.IsNull {
				value = nil
			} else {
				if col.Type == "string" {
					value = ""
				} else if col.Type == "int" {
					value = 0
				}
			}
			values = append(values, value)
		}
		names = append(names, col.Name)
		placeholders = append(placeholders, "?")
	}

	result, err := h.DB.Exec("INSERT INTO "+table.Name+" ("+strings.Join(names, ", ")+") VALUES ("+strings.Join(placeholders, ", ")+")", values...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(Response{map[string]interface{}{table.PrimaryKey: id}, ""})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return

}

func (h *Handler) Post(w http.ResponseWriter, table *Table, q *Query) {

	if q.Id == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(q.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := make(map[string]interface{})
	json.Unmarshal(b, &body)

	names := []string{}
	values := []interface{}{}

	for _, col := range table.Columns {
		if value, matched := body[col.Name]; matched {

			if col.Primary {
				FierldError(col.Name, w)
				return
			}

			switch value.(type) {
			case float64:
				if col.Type == "string" {
					FierldError(col.Name, w)
					return
				}
			case string:
				if col.Type == "int" {
					FierldError(col.Name, w)
					return
				}
			case nil:
				if !col.IsNull {
					FierldError(col.Name, w)
					return
				}
			}

			names = append(names, col.Name+" = ?")
			values = append(values, value)
		}
	}

	values = append(values, q.Id)
	result, err := h.DB.Exec("UPDATE "+table.Name+" SET "+strings.Join(names, ",")+" WHERE "+table.PrimaryKey+" = ?", values...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(Response{map[string]interface{}{"updated": affected}, ""})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func FierldError(name string, w http.ResponseWriter) {
	j, err := json.Marshal(Response{nil, "field " + name + " have invalid type"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write(j)
	return
}

func (h *Handler) Delete(w http.ResponseWriter, table *Table, q *Query) {
	if q.Id == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec("DELETE FROM "+table.Name+" WHERE "+table.PrimaryKey+" = ?", q.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(Response{map[string]interface{}{"deleted": affected}, ""})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return
}
