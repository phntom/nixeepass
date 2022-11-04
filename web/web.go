package web

import (
	"net/http"
)

func RunWebUI() error {
	fs := http.FileServer(http.Dir("./ui/build"))
	http.Handle("/", fs)
	err := http.ListenAndServe(":8999", nil)
	return err
}
