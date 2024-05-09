package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
)

const HttpServePort = 3000
const DbReportTableName = "reports"

// CSPReport represents the structure of a CSP report
type CSPReport struct {
	DocumentURI        string `json:"document-uri"`
	Referrer           string `json:"referrer,omitempty"`
	BlockedURI         string `json:"blocked-uri"`
	ViolatedDirective  string `json:"violated-directive,omitempty"`
	EffectiveDirective string `json:"effective-directive"`
	OriginalPolicy     string `json:"original-policy,omitempty"`
	Disposition        string `json:"disposition,omitempty"`
	StatusCode         int    `json:"status-code,omitempty"`
	LineNumber         int    `json:"line-number,omitempty"`
	ColumnNumber       int    `json:"column-number,omitempty"`
	SourceFile         string `json:"source-file,omitempty"`
	ScriptSample       string `json:"script-sample,omitempty"`
}

type Config struct {
	DbDsn string `json:"datasource"`
}

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func getDataSource() (*Config, error) {
	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		return nil, errors.New("DB_DSN environment variable not set")
	}
	config := Config{
		DbDsn: dbDsn,
	}
	return &config, nil
}

func writeResponse(w http.ResponseWriter, msg string, success bool, statusCode int) error {
	response := &Response{
		Message: msg,
		Success: success,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(response)
}

func validateRequest(r *http.Request) (err error) {
	// restrict to only POST method
	if r.Method != "POST" {
		err = errors.New("Invalid Method")
		return
	}

	// check the content type of the request
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/csp-report" {
		err = errors.New("Invalid Content-Type")
	}

	return
}

func main() {

	config, configErr := getDataSource()
	if configErr != nil {
		fmt.Println("Database DSN Not Supplied")
		return
	}
	db, err := sql.Open("mysql", config.DbDsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Welcome To CSP-Notifier")
		if err != nil {
			return
		}
	})

	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// validate the request
		if err := validateRequest(r); err != nil {
			_ = writeResponse(w, err.Error(), false, http.StatusMethodNotAllowed)
			return
		}
		var jsonPayload struct {
			Report CSPReport `json:"csp-report"`
		}

		if err := json.NewDecoder(r.Body).Decode(&jsonPayload); err != nil {
			_ = writeResponse(w, err.Error(), false, http.StatusBadRequest)
			return
		}

		violationReport := jsonPayload.Report

		if _, err := db.Exec("INSERT INTO "+DbReportTableName+"( "+
			"document_uri, referrer, blocked_uri, violated_directive, "+
			"effective_directive, original_policy, disposition, "+
			"line_number, column_number, source_file, status_code, script_sample) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			violationReport.DocumentURI, violationReport.Referrer, violationReport.BlockedURI,
			violationReport.ViolatedDirective, violationReport.EffectiveDirective, violationReport.OriginalPolicy,
			violationReport.Disposition, violationReport.LineNumber, violationReport.ColumnNumber,
			violationReport.SourceFile, violationReport.StatusCode, violationReport.ScriptSample); err != nil {

			_ = writeResponse(w, err.Error(), false, http.StatusInternalServerError)

			return
		}

		_ = writeResponse(w, "Violation Reported Successfully", true, http.StatusCreated)
	})

	httpAddr := fmt.Sprintf(":%d", HttpServePort)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		panic(err)
	}
}
