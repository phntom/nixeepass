package webui

import (
	"errors"
	"github.com/phntom/nixeepass/orm"
	"io"
	"net/http"
)

func ensureSecret(request *http.Request) error {
	if request.Header.Get("X-Secret") != cfgHttp.LivelinessReadinessSecret {
		return errors.New("invalid X-Secret header")
	}
	return nil
}

func livelinessHandler(writer io.Writer, request *http.Request) error {
	err := ensureSecret(request)
	if err != nil {
		return err
	}
	err = orm.Alive()
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte("alive"))
	return err
}

func readinessHandler(writer io.Writer, request *http.Request) error {
	err := ensureSecret(request)
	if err != nil {
		return err
	}
	err = orm.Alive()
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte("ready"))
	return err
}
