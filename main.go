package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := "hello world"
		fmt.Println(msg)
		w.Write([]byte(`
		<html>
	    <head>
			  <title>チャット</title>	
			</head>
			<body>
			  チャットしましょう！
			</body>
		</html>
		`))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
