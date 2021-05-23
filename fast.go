package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

var dataPool = sync.Pool{
	New: func() interface{} {
		return &User{}
	},
}

func FastSearch(out io.Writer) {
	//SlowSearch(out)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	seenBrowsers := make(map[string]string)
	var isAndroid, isMSIE bool
	b := strings.Builder{}
	i := 0
	in := bufio.NewScanner(file)

	for in.Scan() {
		i++
		user := dataPool.Get().(*User)
		err := user.UnmarshalJSON(in.Bytes())
		dataPool.Put(user)
		if err != nil {
			return
		}

		isAndroid = false
		isMSIE = false

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
			} else if strings.Contains(browser, "MSIE") {
				isMSIE = true
			} else {
				continue
			}
			seenBrowsers[browser] = browser
		}
		if isAndroid && isMSIE {
			email := strings.Replace(user.Email, "@", " [at] ", 1)
			b.WriteString("[" + strconv.Itoa(i-1) + "] " + user.Name + " <" + email + ">\n")
		}
	}

	fmt.Fprint(out, "found users:\n"+b.String()+"\n")
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}
