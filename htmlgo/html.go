package html

import (
	"fmt"
	"io" // moderno: usar io.ReadAll
	"net/http"
	"regexp"
)

func Titulo(urls ...string) <-chan string {
	c := make(chan string)

	for _, url := range urls {
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				c <- fmt.Sprintf("Erro ao acessar %s: %v", url, err)
				return
			}
			defer resp.Body.Close() // sempre fechar!

			html, err := io.ReadAll(resp.Body)
			if err != nil {
				c <- fmt.Sprintf("Erro ao ler %s: %v", url, err)
				return
			}

			// usando raw string para regex:
			r, _ := regexp.Compile(`<title>(.*?)</title>`)
			submatch := r.FindStringSubmatch(string(html))
			if len(submatch) >= 2 {
				c <- submatch[1]
			} else {
				c <- fmt.Sprintf("Título não encontrado em %s", url)
			}
		}(url)
	}

	return c
}
