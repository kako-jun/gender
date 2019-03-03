// sex is ***
package sex

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/kako-jun/sex/statik"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rakyll/statik/fs"
)

// for json
type ALocale struct {
	Translated string `json:"translated"`
	Gender     string `json:"gender"`
}

type Translation struct {
	En string    `json:"en"`
	Fr []ALocale `json:"fr"`
}

// Sex is ***
type Sex struct{}

func (sex Sex) exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return true
}

func (sex Sex) full_gender(id string) (full string) {
	switch id {
	case "f":
		full = "feminine"
	case "m":
		full = "masculine"
	case "n":
		full = "neuter"
	}

	return full
}

func (sex Sex) findKeyword(keyword string) (results [][]string, err error) {

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	f, err := statikFS.Open("/translation.db")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	stats, statsErr := f.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(f)
	_, err = bufr.Read(bytes)

	// // windows: C:\Users\{{user}}\AppData\Local\Temp
	// // mac: /var/folders/t9/{{id}}/T/
	// // linux: /tmp
	// db_file_path := os.TempDir() + "/sex/translation.db"
	// file, err := os.Create(db_file_path)
	// if err != nil {
	// 	panic(err)
	// }

	// defer file.Close()

	// // overwrite force.
	// file.Write(bytes)

	var db_file_path string = "./translation.db"
	db, err := sql.Open("sqlite3", db_file_path)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rows, err := db.Query(`SELECT * FROM view_translation t WHERE t.search_keyword LIKE '%` + keyword + `%'`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		cols := [8]string{}
		err = rows.Scan(&cols[0], &cols[1], &cols[2], &cols[3], &cols[4], &cols[5], &cols[6], &cols[7])
		if err != nil {
			panic(err)
		}

		results = append(results, cols[:])
	}

	return results, err
}

func (sex Sex) translate(keyword string) (translations []Translation, err error) {
	rows, err := sex.findKeyword(keyword)
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		en := row[0]
		fr_1_translated := row[1]
		fr_1_gender := row[2]
		fr_2_translated := row[3]
		fr_2_gender := row[4]
		fr_3_translated := row[5]
		fr_3_gender := row[6]
		// search_keyword := row[7]
		// fmt.Println(en, fr_1_translated, fr_1_gender, fr_2_translated, fr_2_gender, fr_3_translated, fr_3_gender, search_keyword)

		var translation Translation
		translation.En = en
		var fr []ALocale

		{
			var aLocale ALocale
			aLocale.Translated = fr_1_translated
			aLocale.Gender = sex.full_gender(fr_1_gender)
			fr = append(fr, aLocale)
		}

		{
			var aLocale ALocale
			aLocale.Translated = fr_2_translated
			aLocale.Gender = sex.full_gender(fr_2_gender)
			fr = append(fr, aLocale)
		}

		{
			var aLocale ALocale
			aLocale.Translated = fr_3_translated
			aLocale.Gender = sex.full_gender(fr_3_gender)
			fr = append(fr, aLocale)
		}

		translation.Fr = fr
		translations = append(translations, translation)
	}

	return translations, err
}

func (sex Sex) Start(keyword string, args []string) (err error) {
	translations, err := sex.translate(keyword)
	if err != nil {
		panic(err)
	}

	for i, translation := range translations {
		if i > 0 {
			fmt.Println("")
		}

		fmt.Println("English: " + translation.En)
		french := "French: "
		for j, aLocale := range translation.Fr {
			if j > 0 {
				french += ", "
			}

			french += aLocale.Translated + " (" + aLocale.Gender + ")"
		}

		french += "\n"
		fmt.Printf(french)
	}

	jsonBytes, err := json.Marshal(&translations)
	if err != nil {
		panic(err)
	}

	out := new(bytes.Buffer)
	json.Indent(out, jsonBytes, "", "  ")
	fmt.Println(out.String())
	return
}

// Exec is ***
func Exec(keyword string, args []string) (errReturn error) {
	sex := new(Sex)
	if err := sex.Start(keyword, args); err != nil {
		fmt.Println("error:", err)
		errReturn = errors.New("error")
		return
	}

	return
}
