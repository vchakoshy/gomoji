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

func StringToLines(s string) []string {
	var lines []string
	var allRow []string

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		replacer := strings.NewReplacer("+", "000", "&", "and")
		result := replacer.Replace(scanner.Text())

		if len(result) != 0 {
			line := strings.TrimSpace(result)
			lines = append(lines, line)
		}
	}
	allRow = append(allRow, strings.Join(lines, ", "))

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return allRow
}

func processTr(tr *goquery.Selection, fRstOutput *os.File) {
	a := []string{}
	header := []string{"Count", "Code", "Browser", "B&W*", "Apple", "Andr",
		"One", "Twit", "Wind", "GMail", "DCM", "KDDI", "SB",
		"Name", "Version", "Default", "Annotations"}
	dict := map[string]string{}

	tr.Find("td").Each(func(indexOfTd int, td *goquery.Selection) {
		lines := StringToLines(td.Text())
		for _, line := range lines {
			if line != " " {
				a = append(a, line)
			}
		}
	})
	for k, v := range header {
		if len(a) != 0 {
			dict[v] = a[k]
		}
	}
	replacer := strings.NewReplacer(" ", "_", ":", "")
	name := replacer.Replace(dict["Annotations"])
	code := strings.Replace(dict["Code"], "U", "\\U", -1)
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
