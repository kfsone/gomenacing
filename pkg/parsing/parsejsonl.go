// Parser for JSONL file translation with parallelism.

package parsing

import (
	"bufio"
	"io"
	"log"

	"github.com/tidwall/gjson"
)

// First worker identifies valid JSONL lines - that is lines which
// contain well-formed JSON.
func getJSONLines(source io.Reader) <-chan string {
	channel := make(chan string, 1)
	go func() {
		defer close(channel)

		badLines := false
		scanner := bufio.NewScanner(source)

		for scanner.Scan() {
			if scanner.Err() != nil {
				panic(scanner.Err())
			}
			line := scanner.Text()
			if !gjson.Valid(line) {
				if !badLines {
					log.Printf("bad json: %s", line)
					badLines = true
				}
				continue
			}
			channel <- line
		}
	}()

	return channel
}

// Second worker consumes valid JSON lines and maps them into the desired
// ordered array list.
func parseJSONLines(lines <-chan string, fields []string) <-chan []*gjson.Result {
	channel := make(chan []*gjson.Result, 1)
	go func() {
		defer close(channel)
		badLines := false
		for line := range lines {
			result := gjson.Parse(line)
			if !result.IsObject() {
				if !badLines {
					log.Printf("malformed jsonl: %s", line)
					badLines = true
				}
				continue
			}
			jsonFields := result.Map()
			slice := make([]*gjson.Result, len(fields))
			invalid := false
			for idx, field := range fields {
				if jsonField, ok := jsonFields[field]; ok {
					slice[idx] = &jsonField
				} else {
					if !invalid {
						log.Printf("missing \"%s\" field in line: %s", field, line)
						invalid = true
					}
					break
				}
			}
			if !invalid {
				channel <- slice
			}
		}
	}()

	return channel
}

// ParseJSONLines will consume a JSONL source and translate JSON dictionaries
// containing the keys listed in `fields` into an array of `gjson.Result`s in
// the same order as `fields`.
//
// Thus a json line representing `{"for":1, "the":2, "win":3}` would return a row
// equivalent to `{3, 1}` if fields was `[]string{ "win", "for" }`
//
func ParseJSONLines(source io.Reader, fields []string) <-chan []*gjson.Result {
	return parseJSONLines(getJSONLines(source), fields)
}
