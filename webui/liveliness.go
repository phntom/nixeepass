package webui

import (
	"github.com/phntom/nixeepass/orm"
	"io"
	"net/http"
)

func livelinessHandler(writer io.Writer, request *http.Request) error {
	err := orm.Alive()
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte("alive"))
	return err
}

func readinessHandler(writer io.Writer, request *http.Request) error {
	_, err := writer.Write([]byte("ready"))
	return err
}
