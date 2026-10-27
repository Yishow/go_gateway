package dbtarget

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	datalinkbase "go-gateway/internal/datalink"

	"github.com/microsoft/go-mssqldb/msdsn"
)

type sqlServerConnectionSettings struct {
	DSN      string
	Database string
	Schema   string
}

const sqlServerEncryptionRequiredOption = "true"

// sqlServerConnectionSettingsFor builds only structured SQL authentication.
// Raw driver strings and certificate-verification bypasses cannot override the
// saved database/schema scope or the encrypted, verified transport contract.
func sqlServerConnectionSettingsFor(config ConnectionConfig) (sqlServerConnectionSettings, error) {
	invalid := func() (sqlServerConnectionSettings, error) {
		return sqlServerConnectionSettings{}, errors.New("sqlserver connection requires explicit host/database/schema and SQL authentication with verified TLS")
	}
	if stringConfigValue(config, "dsn") != "" || stringConfigValue(config, "servercertificate") != "" {
		return invalid()
	}
	host := stringConfigValue(config, "host")
	database := stringConfigValue(config, "database")
	schemaName := stringConfigValue(config, "schema")
	user := stringConfigValue(config, "user")
	password, ok := config["password"].(string)
	if host == "" || strings.ContainsAny(host, "/\\?#@\x00") || database == "" || schemaName == "" ||
		user == "" || strings.ContainsRune(user, '\\') || !ok || password == "" {
		return invalid()
	}
	portText := stringConfigValue(config, "port")
	if portText == "" {
		portText = "1433"
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port < 1 || port > 65535 {
		return invalid()
	}
	encrypt := sqlServerEncryptionRequiredOption
	if value, present := config["encrypt"]; present {
		switch strings.ToLower(strings.TrimSpace(sqlServerOptionText(value))) {
		case sqlServerEncryptionRequiredOption, "mandatory", "yes", "1", "t":
		case "strict":
			encrypt = "strict"
		default:
			return invalid()
		}
	}
	for _, key := range []string{"trustservercertificate", "trust_server_certificate"} {
		if value, present := config[key]; present {
			value, err := strconv.ParseBool(sqlServerOptionText(value))
			if err != nil || value {
				return invalid()
			}
		}
	}
	for _, key := range []string{"use_tls", "tls"} {
		if value, present := config[key]; present {
			enabled, err := strconv.ParseBool(sqlServerOptionText(value))
			if err != nil || !enabled {
				return invalid()
			}
		}
	}
	tlsMinimum := stringConfigValue(config, "tlsmin")
	if tlsMinimum == "" {
		tlsMinimum = "1.2"
	}
	if tlsMinimum != "1.2" && tlsMinimum != "1.3" {
		return invalid()
	}
	timeout := 5 * time.Second
	if text := stringConfigValue(config, "timeout"); text != "" {
		timeout, err = time.ParseDuration(text)
		if err != nil || timeout <= 0 {
			return invalid()
		}
	}
	seconds := strconv.FormatInt(int64((timeout-1)/time.Second)+1, 10)
	query := url.Values{
		"database": {database}, "encrypt": {encrypt}, "TrustServerCertificate": {"false"},
		"tlsmin": {tlsMinimum}, "connection timeout": {seconds}, "dial timeout": {seconds},
	}
	for _, key := range []string{"certificate", "hostnameincertificate"} {
		if value, present := config[key]; present {
			text, ok := value.(string)
			if !ok || strings.TrimSpace(text) == "" {
				return invalid()
			}
			query.Set(key, text)
		}
	}
	dsn := (&url.URL{
		Scheme: "sqlserver", User: url.UserPassword(user, password),
		Host: net.JoinHostPort(strings.Trim(host, "[]"), strconv.Itoa(port)), RawQuery: query.Encode(),
	}).String()
	parsed, err := msdsn.Parse(dsn)
	if err != nil || parsed.TLSConfig == nil || parsed.TLSConfig.InsecureSkipVerify || parsed.TrustServerCertificate ||
		(parsed.Encryption != msdsn.EncryptionRequired && parsed.Encryption != msdsn.EncryptionStrict) ||
		parsed.Database != database || parsed.User != user || parsed.Password != password {
		return invalid()
	}
	return sqlServerConnectionSettings{DSN: dsn, Database: database, Schema: schemaName}, nil
}

func sqlServerOptionText(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case bool:
		return strconv.FormatBool(typed)
	default:
		return ""
	}
}

func buildSQLServerDBConfig(config ConnectionConfig) (datalinkbase.DBConfig, error) {
	settings, err := sqlServerConnectionSettingsFor(config)
	if err != nil {
		return datalinkbase.DBConfig{}, err
	}
	return datalinkbase.DBConfig{
		Type: datalinkbase.DBTypeSQLServer, DSN: settings.DSN, MaxOpenConns: 3, MaxIdleConns: 1,
	}, nil
}
