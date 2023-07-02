package plugin

import (
	"io"
	"net"
)

type LogstashConfig struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

type logstashWriter struct {
	conn net.Conn
}

func NewLogstash(c LogstashConfig) io.Writer {
	conn, err := net.Dial("tcp", c.Host+":"+c.Port)
	if err != nil {
		panic(err)
	}
	return &logstashWriter{conn: conn}
}

func (l *logstashWriter) Write(p []byte) (n int, err error) {
	return l.conn.Write(p)
}
