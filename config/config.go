package config

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Data holds all API settings
type Data struct {
	API struct {
		BaseURL         string
		PathPrefix      string
		Bind            string
		SwaggerAPIPath  string
		SwaggerPath     string
		SwaggerFilePath string
		ImageFilePath   string
	}

	Connections struct {
		Logger     LoggerConnection
		PostgreSQL PostgreSQLConnection
		AMQP       AMQPConnection
		Email      EmailConfig
		PayPal     string
	}

	PaymentProviders struct {
		PayPal struct {
			ClientID string
			Secret   string
		}
	}

	EmailTemplates Templates

	Web struct {
		BaseURL string
	}
}

// EmailConfig contains all email settings
type EmailConfig struct {
	AdminEmail string
	ReplyTo    string
	SMTP       struct {
		User     string
		Password string
		Server   string
		Port     int
	}
}

// EmailTemplate holds all values of one email template
type EmailTemplate struct {
	Subject string
	Text    string
	HTML    string
}

// Templates holds all email templates
type Templates struct {
	PaymentConfirmation EmailTemplate
}

// LoggerConnection contains all of the logger settings
type LoggerConnection struct {
	Protocol string
	Address  string
}

// PostgreSQLConnection contains all of the db configuration values
type PostgreSQLConnection struct {
	User     string
	Password string
	Host     string
	Port     int
	DbName   string
	SslMode  string
}

// AMQPConnection contains all of the db configuration values
type AMQPConnection struct {
	User     string
	Password string
	Host     string
	Port     int
	Broker   string
	Queue    string
}

var (
	// Settings contains the parsed configuration values
	Settings *Data
)

// ParseSettings parses the config file
func ParseSettings() {
	logLevelStr := flag.String("loglevel", "info", "Log level")
	configFile := flag.String("configfile", "config.json", "config file in the JSON format")
	flag.Parse()

	logLevel, err := log.ParseLevel(*logLevelStr)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)

	if configFile == nil || len(*configFile) == 0 {
		log.Panic(errors.New("Did not get a config file passed in"))
	}
	log.WithField("File", *configFile).Info("Using config file")

	// Parse config file
	configData := Data{}
	handler := NewHandler(*configFile, &configData, nil)
	if handler == nil {
		log.Fatal(errors.New("Config handler is nil, cannot continue"))
	}
	if !handler.LastReadValid() {
		log.WithField(
			"File",
			*configFile,
		).Fatal(errors.New("Did not get a valid config file"))
	}

	Settings = handler.CurrentData().(*Data)
	//FIXME catch empty conf fields
}

// Marshal returns a "Connection String" with escaped values of all non-empty
// fields as described at http://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING
func (c *PostgreSQLConnection) Marshal() string {
	val := reflect.ValueOf(c).Elem()
	var out string
	l := val.NumField()

	r := strings.NewReplacer(`'`, `\'`, `\`, `\\`)

	for i := 0; i < l; i++ {
		var fieldValue string

		switch f := val.Field(i).Interface().(type) {
		case string:
			fieldValue = f
		case int:
			if f == 0 {
				continue
			}
			fieldValue = strconv.Itoa(f)
		}
		fieldType := val.Type().Field(i).Name

		if len(fieldValue) > 0 {
			out += strings.ToLower(fieldType) + "='" + r.Replace(fieldValue) + "'"
			if i < l {
				out += " "
			}
		}
	}

	return out
}

// Marshal returns an AMQP "Connection String"
func (c *AMQPConnection) Marshal() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, c.Broker)
}
