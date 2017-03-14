package gomoji

/*
Extract the full list of emoji and names from the Unicode Consortium and
apply as much formatting as possible so the codes can be dropped into the
emoji registry file.

Written and run htmlTableToMap(url, outputFilePath) with this parameters

    url := "http://www.unicode.org/emoji/charts/full-emoji-list.html"
    outputFilePath := "path/your/local"

*/

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strings"
)

func StringToLines(text string) []string {
	var lines []string
	var allRow []string
	var replacer *strings.Replacer

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {

		if ok := strings.HasPrefix(scanner.Text(), "U+"); ok {
			unicodes := strings.Split(scanner.Text(), " ")
			for _, v := range unicodes {
				if len(v) == 6 {
					replacer = strings.NewReplacer(`U+`, `\u`)
				} else {
					replacer = strings.NewReplacer(`+`, `000`, `U`, `\U`)
				}
				result := replacer.Replace(v)
				lines = append(lines, strings.TrimSpace(result))
			}

		} else {

			replacer = strings.NewReplacer("&", "and", "+", "000", "&", "and", " ", "_", ":", "", ",", "", "“", "", "”", "")
			result := replacer.Replace(scanner.Text())
			lines = append(lines, strings.TrimSpace(result))
		}
	}
	allRow = append(allRow, strings.Join(lines, ""))

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return allRow
}

func processTr(tr *goquery.Selection, fRstOutput *os.File) {
	a := []string{}
	header := []string{"Count", "Code", "Appl", "Goog", "Twtr.", "One", "FB", "FBM", "Sams.", "Wind.", "GMail", "SB", "DCM", "KDDI", "Name", "Date", "Keywords"}
	dict := map[string]string{}

	tr.Find("td").Each(func(indexOfTd int, td *goquery.Selection) {
		lines := StringToLines(td.Text())
		for _, line := range lines {
			a = append(a, line)
		}
	})

	for k, v := range header {
		if len(a) == 17 {
			dict[v] = a[k]
		}
	}

	name := dict["Name"]
	code := dict["Code"]
	str := name + ":" + code
	if str != ":" {
		name = ":" + name + ":"
		fmt.Fprintf(fRstOutput, "%q:%q,\n", name, code)
	}

}

func htmlTableToMap(url, outputFilePath string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	fRstOutput, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(fRstOutput, "{")
	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(_ int, tr *goquery.Selection) {
			processTr(tr, fRstOutput)
		})

	})
	fmt.Fprintf(fRstOutput, "}")
}
