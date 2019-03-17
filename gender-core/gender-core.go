// Package gender is ***
package gender

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/kako-jun/gender/statik"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rakyll/statik/fs"
)

// ALocale for json
type ALocale struct {
	Translated string `json:"translated"`
	Gender     string `json:"gender"`
}

// Translation for json
type Translation struct {
	En string    `json:"en"`
	Ar []ALocale `json:"ar"`
	Fr []ALocale `json:"fr"`
	De []ALocale `json:"de"`
	Hi []ALocale `json:"hi"`
	It []ALocale `json:"it"`
	Pt []ALocale `json:"pt"`
	Ru []ALocale `json:"ru"`
	Es []ALocale `json:"es"`
}

// Gender is ***
type Gender struct{}

func (gender Gender) exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return true
}

func (gender Gender) fullGender(id string) (full string) {
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

func (gender Gender) createQuery(keyword string, exactFlag bool, closestFlag bool, arFlag bool, frFlag bool, deFlag bool, hiFlag bool, itFlag bool, ptFlag bool, ruFlag bool, esFlag bool) (query string) {
	where := ` WHERE en LIKE '` + keyword + `'`

	if exactFlag {
		if arFlag {
			where += ` OR ar_translated_1 LIKE '` + keyword + `'`
			if !closestFlag {
				where += ` OR ar_translated_2 LIKE '` + keyword + `'`
				where += ` OR ar_translated_3 LIKE '` + keyword + `'`
			}
		}

		if frFlag {
			where += ` OR fr_translated_1 LIKE '` + keyword + `'`
			if !closestFlag {
				where += ` OR fr_translated_2 LIKE '` + keyword + `'`
				where += ` OR fr_translated_3 LIKE '` + keyword + `'`
			}
		}

		if deFlag {
			where += ` OR de_translated_1 LIKE '` + keyword + `'`
			if !closestFlag {
				where += ` OR de_translated_2 LIKE '` + keyword + `'`
				where += ` OR de_translated_3 LIKE '` + keyword + `'`
			}
		}

		if hiFlag {
			where += ` OR hi_translated_1 LIKE '` + keyword + `'`
			if !closestFlag {
				where += ` OR hi_translated_2 LIKE '` + keyword + `'`
				where += ` OR hi_translated_3 LIKE '` + keyword + `'`
			}
		}

		if itFlag {
			where += ` OR it_translated_1 LIKE '` + keyword + `'`
			if !closestFlag {
				where += ` OR it_translated_2 LIKE '` + keyword + `'`
				where += ` OR it_translated_3 LIKE '` + keyword + `'`
			}
		}

		if ptFlag {
			where += ` OR pt_translated_1 LIKE '` + keyword + `'`
			if !closestFlag {
				where += ` OR pt_translated_2 LIKE '` + keyword + `'`
				where += ` OR pt_translated_3 LIKE '` + keyword + `'`
			}
		}

		if ruFlag {
			where += ` OR ru_translated_1 LIKE '` + keyword + `'`
			if !closestFlag {
				where += ` OR ru_translated_2 LIKE '` + keyword + `'`
				where += ` OR ru_translated_3 LIKE '` + keyword + `'`
			}
		}

		if esFlag {
			where += ` OR es_translated_1 LIKE '` + keyword + `'`
			if !closestFlag {
				where += ` OR es_translated_2 LIKE '` + keyword + `'`
				where += ` OR es_translated_3 LIKE '` + keyword + `'`
			}
		}
	} else {
		if arFlag {
			if closestFlag {
				where += ` OR ar_translated_1 LIKE '%` + keyword + `%'`
			} else {
				where += ` OR ar_search_keyword LIKE '%` + keyword + `%'`
			}
		}

		if frFlag {
			if closestFlag {
				where += ` OR fr_translated_1 LIKE '%` + keyword + `%'`
			} else {
				where += ` OR fr_search_keyword LIKE '%` + keyword + `%'`
			}
		}

		if deFlag {
			if closestFlag {
				where += ` OR de_translated_1 LIKE '%` + keyword + `%'`
			} else {
				where += ` OR de_search_keyword LIKE '%` + keyword + `%'`
			}
		}

		if hiFlag {
			if closestFlag {
				where += ` OR hi_translated_1 LIKE '%` + keyword + `%'`
			} else {
				where += ` OR hi_search_keyword LIKE '%` + keyword + `%'`
			}
		}

		if itFlag {
			if closestFlag {
				where += ` OR it_translated_1 LIKE '%` + keyword + `%'`
			} else {
				where += ` OR it_search_keyword LIKE '%` + keyword + `%'`
			}
		}

		if ptFlag {
			if closestFlag {
				where += ` OR pt_translated_1 LIKE '%` + keyword + `%'`
			} else {
				where += ` OR pt_search_keyword LIKE '%` + keyword + `%'`
			}
		}

		if ruFlag {
			if closestFlag {
				where += ` OR ru_translated_1 LIKE '%` + keyword + `%'`
			} else {
				where += ` OR ru_search_keyword LIKE '%` + keyword + `%'`
			}
		}

		if esFlag {
			if closestFlag {
				where += ` OR es_translated_1 LIKE '%` + keyword + `%'`
			} else {
				where += ` OR es_search_keyword LIKE '%` + keyword + `%'`
			}
		}
	}

	query = `SELECT * FROM view_translation` + where
	return
}

func (gender Gender) findKeyword(keyword string, exactFlag bool, closestFlag bool, arFlag bool, frFlag bool, deFlag bool, hiFlag bool, itFlag bool, ptFlag bool, ruFlag bool, esFlag bool) (results [][]string, err error) {
	// windows: C:\Users\{{user}}\AppData\Local\Temp
	// mac: /var/folders/t9/{{id}}/T/
	// linux: /tmp
	dbDirPath := os.TempDir() + "/gender"
	dbFilePath := dbDirPath + "/translation_1.0.0.db"

	if !gender.exists(dbFilePath) {
		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}

		f, err := statikFS.Open("/translation.db")
		if err != nil {
			log.Fatal(err)
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

		os.Mkdir(dbDirPath, 0755)

		file, err := os.Create(dbFilePath)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		file.Write(bytes)
	}

	// dbFilePath := "./translation.db"
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	query := gender.createQuery(keyword, exactFlag, closestFlag, arFlag, frFlag, deFlag, hiFlag, itFlag, ptFlag, ruFlag, esFlag)
	// fmt.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		cols := [57]string{}
		err = rows.Scan(&cols[0], &cols[1], &cols[2], &cols[3], &cols[4], &cols[5], &cols[6], &cols[7], // ar
			&cols[8], &cols[9], &cols[10], &cols[11], &cols[12], &cols[13], &cols[14], // fr
			&cols[15], &cols[16], &cols[17], &cols[18], &cols[19], &cols[20], &cols[21], // de
			&cols[22], &cols[23], &cols[24], &cols[25], &cols[26], &cols[27], &cols[28], // hi
			&cols[29], &cols[30], &cols[31], &cols[32], &cols[33], &cols[34], &cols[35], // it
			&cols[36], &cols[37], &cols[38], &cols[39], &cols[40], &cols[41], &cols[42], // pt
			&cols[43], &cols[44], &cols[45], &cols[46], &cols[47], &cols[48], &cols[49], // ru
			&cols[50], &cols[51], &cols[52], &cols[53], &cols[54], &cols[55], &cols[56]) // es
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, cols[:])
	}

	return results, err
}

func (gender Gender) createALocales(translated1 string, gender1 string, translated2 string, gender2 string, translated3 string, gender3 string) (aLocales []ALocale) {
	if translated1 != "" {
		var aLocale ALocale
		aLocale.Translated = translated1
		aLocale.Gender = gender.fullGender(gender1)
		aLocales = append(aLocales, aLocale)
	}

	if translated2 != "" {
		var aLocale ALocale
		aLocale.Translated = translated2
		aLocale.Gender = gender.fullGender(gender2)
		aLocales = append(aLocales, aLocale)
	}

	if translated3 != "" {
		var aLocale ALocale
		aLocale.Translated = translated3
		aLocale.Gender = gender.fullGender(gender3)
		aLocales = append(aLocales, aLocale)
	}

	return
}

func (gender Gender) translate(keyword string, exactFlag bool, closestFlag bool, arFlag bool, frFlag bool, deFlag bool, hiFlag bool, itFlag bool, ptFlag bool, ruFlag bool, esFlag bool) (translations []Translation, err error) {
	rows, err := gender.findKeyword(keyword, exactFlag, closestFlag, arFlag, frFlag, deFlag, hiFlag, itFlag, ptFlag, ruFlag, esFlag)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		en := row[0]

		arTranslated1 := row[1]
		arGender1 := row[2]
		arTranslated2 := ""
		arGender2 := ""
		arTranslated3 := ""
		arGender3 := ""
		if !closestFlag {
			arTranslated2 = row[3]
			arGender2 = row[4]
			arTranslated3 = row[5]
			arGender3 = row[6]
		}

		frTranslated1 := row[8]
		frGender1 := row[9]
		frTranslated2 := ""
		frGender2 := ""
		frTranslated3 := ""
		frGender3 := ""
		if !closestFlag {
			frTranslated2 = row[10]
			frGender2 = row[11]
			frTranslated3 = row[12]
			frGender3 = row[13]
		}

		deTranslated1 := row[15]
		deGender1 := row[16]
		deTranslated2 := ""
		deGender2 := ""
		deTranslated3 := ""
		deGender3 := ""
		if !closestFlag {
			deTranslated2 = row[17]
			deGender2 = row[18]
			deTranslated3 = row[19]
			deGender3 = row[20]
		}

		hiTranslated1 := row[22]
		hiGender1 := row[23]
		hiTranslated2 := ""
		hiGender2 := ""
		hiTranslated3 := ""
		hiGender3 := ""
		if !closestFlag {
			hiTranslated2 = row[24]
			hiGender2 = row[25]
			hiTranslated3 = row[26]
			hiGender3 = row[27]
		}

		itTranslated1 := row[29]
		itGender1 := row[30]
		itTranslated2 := ""
		itGender2 := ""
		itTranslated3 := ""
		itGender3 := ""
		if !closestFlag {
			itTranslated2 = row[31]
			itGender2 = row[32]
			itTranslated3 = row[33]
			itGender3 = row[34]
		}

		ptTranslated1 := row[36]
		ptGender1 := row[37]
		ptTranslated2 := ""
		ptGender2 := ""
		ptTranslated3 := ""
		ptGender3 := ""
		if !closestFlag {
			ptTranslated2 = row[38]
			ptGender2 = row[39]
			ptTranslated3 = row[40]
			ptGender3 = row[41]
		}

		ruTranslated1 := row[43]
		ruGender1 := row[44]
		ruTranslated2 := ""
		ruGender2 := ""
		ruTranslated3 := ""
		ruGender3 := ""
		if !closestFlag {
			ruTranslated2 = row[45]
			ruGender2 = row[46]
			ruTranslated3 = row[47]
			ruGender3 = row[48]
		}

		esTranslated1 := row[50]
		esGender1 := row[51]
		esTranslated2 := ""
		esGender2 := ""
		esTranslated3 := ""
		esGender3 := ""
		if !closestFlag {
			esTranslated2 = row[52]
			esGender2 = row[53]
			esTranslated3 = row[54]
			esGender3 = row[55]
		}

		var translation Translation
		translation.En = en

		ar := gender.createALocales(arTranslated1, arGender1, arTranslated2, arGender2, arTranslated3, arGender3)
		fr := gender.createALocales(frTranslated1, frGender1, frTranslated2, frGender2, frTranslated3, frGender3)
		de := gender.createALocales(deTranslated1, deGender1, deTranslated2, deGender2, deTranslated3, deGender3)
		hi := gender.createALocales(hiTranslated1, hiGender1, hiTranslated2, hiGender2, hiTranslated3, hiGender3)
		it := gender.createALocales(itTranslated1, itGender1, itTranslated2, itGender2, itTranslated3, itGender3)
		pt := gender.createALocales(ptTranslated1, ptGender1, ptTranslated2, ptGender2, ptTranslated3, ptGender3)
		ru := gender.createALocales(ruTranslated1, ruGender1, ruTranslated2, ruGender2, ruTranslated3, ruGender3)
		es := gender.createALocales(esTranslated1, esGender1, esTranslated2, esGender2, esTranslated3, esGender3)

		if arFlag {
			translation.Ar = ar
		}

		if frFlag {
			translation.Fr = fr
		}

		if deFlag {
			translation.De = de
		}

		if hiFlag {
			translation.Hi = hi
		}

		if itFlag {
			translation.It = it
		}

		if ptFlag {
			translation.Pt = pt
		}

		if ruFlag {
			translation.Ru = ru
		}

		if esFlag {
			translation.Es = es
		}

		translations = append(translations, translation)
	}

	return translations, err
}

func (gender Gender) createSimpleALocaleString(aLocales []ALocale) (languageString string) {
	for j, aLocale := range aLocales {
		if aLocale.Translated != "" {
			if j > 0 {
				languageString += ", "
			}

			languageString += aLocale.Translated + " (" + aLocale.Gender + ")"
		}
	}

	return
}

func (gender Gender) createALocaleString(language string, aLocales []ALocale) (languageString string) {
	languageString = language + ": "
	for j, aLocale := range aLocales {
		if aLocale.Translated != "" {
			if j > 0 {
				languageString += ", "
			}

			languageString += aLocale.Translated + " (" + aLocale.Gender + ")"
		}
	}

	return
}

func (gender Gender) start(keyword string, exactFlag bool, closestFlag bool, simpleFlag bool, jsonFlag bool, arFlag bool, frFlag bool, deFlag bool, hiFlag bool, itFlag bool, ptFlag bool, ruFlag bool, esFlag bool) (err error) {
	translations, err := gender.translate(keyword, exactFlag, closestFlag, arFlag, frFlag, deFlag, hiFlag, itFlag, ptFlag, ruFlag, esFlag)
	if err != nil {
		log.Fatal(err)
	}

	if simpleFlag {
		for i, translation := range translations {
			if i > 0 {
				fmt.Println("")
			}

			fmt.Println(translation.En)
			if arFlag {
				ar := gender.createSimpleALocaleString(translation.Ar)
				fmt.Println(ar)
			}

			if frFlag {
				fr := gender.createSimpleALocaleString(translation.Fr)
				fmt.Println(fr)
			}

			if deFlag {
				de := gender.createSimpleALocaleString(translation.De)
				fmt.Println(de)
			}

			if hiFlag {
				hi := gender.createSimpleALocaleString(translation.Hi)
				fmt.Println(hi)
			}

			if itFlag {
				it := gender.createSimpleALocaleString(translation.It)
				fmt.Println(it)
			}

			if ptFlag {
				pt := gender.createSimpleALocaleString(translation.Pt)
				fmt.Println(pt)
			}

			if ruFlag {
				ru := gender.createSimpleALocaleString(translation.Ru)
				fmt.Println(ru)
			}

			if esFlag {
				es := gender.createSimpleALocaleString(translation.Es)
				fmt.Println(es)
			}
		}
	} else if jsonFlag {
		jsonBytes, err := json.Marshal(&translations)
		if err != nil {
			log.Fatal(err)
		}

		jsonBuffer := new(bytes.Buffer)
		err = json.Indent(jsonBuffer, jsonBytes, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(jsonBuffer.String())
	} else {
		for i, translation := range translations {
			if i > 0 {
				fmt.Println("")
			}

			fmt.Println("English: " + translation.En)

			if arFlag {
				ar := gender.createALocaleString("Arabic", translation.Ar)
				fmt.Println(ar)
			}

			if frFlag {
				fr := gender.createALocaleString("French", translation.Fr)
				fmt.Println(fr)
			}

			if deFlag {
				de := gender.createALocaleString("German", translation.De)
				fmt.Println(de)
			}

			if hiFlag {
				hi := gender.createALocaleString("Hindi", translation.Hi)
				fmt.Println(hi)
			}

			if itFlag {
				it := gender.createALocaleString("Italian", translation.It)
				fmt.Println(it)
			}

			if ptFlag {
				pt := gender.createALocaleString("Portuguese", translation.Pt)
				fmt.Println(pt)
			}

			if ruFlag {
				ru := gender.createALocaleString("Russian", translation.Ru)
				fmt.Println(ru)
			}

			if esFlag {
				es := gender.createALocaleString("Spanish", translation.Es)
				fmt.Println(es)
			}
		}
	}

	return
}

// Exec is ***
func Exec(keyword string, exactFlag bool, closestFlag bool, simpleFlag bool, jsonFlag bool, arFlag bool, frFlag bool, deFlag bool, hiFlag bool, itFlag bool, ptFlag bool, ruFlag bool, esFlag bool) (errReturn error) {
	if simpleFlag {
		jsonFlag = false
	}

	if !arFlag && !frFlag && !deFlag && !hiFlag && !itFlag && !ptFlag && !ruFlag && !esFlag {
		arFlag = true
		frFlag = true
		deFlag = true
		hiFlag = true
		itFlag = true
		ptFlag = true
		ruFlag = true
		esFlag = true
	}

	gender := new(Gender)
	if err := gender.start(keyword, exactFlag, closestFlag, simpleFlag, jsonFlag, arFlag, frFlag, deFlag, hiFlag, itFlag, ptFlag, ruFlag, esFlag); err != nil {
		fmt.Println("error:", err)
		errReturn = errors.New("error")
		return
	}

	return
}
