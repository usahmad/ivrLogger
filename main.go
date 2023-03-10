package main

import (
	db2 "IvrParser/db"
	"IvrParser/db/models"
	"bufio"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func checkEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		os.Exit(1)
	}
	envs := []string{
		"DB_USERNAME",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"LOG_DIRECTORY",
	}
	for _, env := range envs {
		_, exists := os.LookupEnv(env)
		if !exists {
			fmt.Println(fmt.Sprintf("%v does not exist in .env", env))
			os.Exit(1)
		}
	}
}

func scanDir() []string {
	var directory, _ = os.LookupEnv("LOG_DIRECTORY")

	files, err := filepath.Glob(filepath.Join(directory, "*"))
	if err != nil {
		panic(err)
	}

	return files
}

func main() {
	fmt.Println("IVR parser by Us.@hmad started blЭт")
	checkEnv()
	defer timeTrack(time.Now(), "Execution")
	fmt.Println("Connecting to DB...")
	db := db2.InitDb()
	types := getTypes(db)
	if db == nil {
		fmt.Println("ERROR CONNECTING TO DB")
		os.Exit(1)
	}
	for _, file := range scanDir() {
		fmt.Println("Reading file: " + file)
		data := readFile(types, file)
		if len(data) > 0 {
			fmt.Println("Inserting Data of: " + file)
			err := models.CreateBulk(db, data)
			check(err)
			fmt.Println(fmt.Sprintf("Removing file %v", file))
			err = os.Remove(file)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(fmt.Sprintf("Removing Empty File %v", file))
			err := os.Remove(file)
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Done")
}

func getTypes(db *gorm.DB) map[int]string {
	var items []models.IvrDetail
	err := models.GetAll(db, &items)
	check(err)
	details := make(map[int]string)
	for _, val := range items {
		details[val.ID] = val.Description
	}
	return details
}

func makeDate(fileName string) string {
	str := strings.Split(fileName, "-")
	if len(str) != 2 {
		return ""
	}
	if len(str[1]) != 8 {
		return ""
	}
	year := str[1][0:4]
	month := str[1][4:6]
	day := str[1][6:8]
	return fmt.Sprintf("%v-%v-%v", day, month, year)
}

func readFile(types map[int]string, fileName string) []models.IVR {
	date := makeDate(fileName)
	if date == "" {
		return []models.IVR{}
	}
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	var data []models.IVR
	amounts := make(map[string]int)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "[s@ivr") {
			re := regexp.MustCompile("@ivr-([0-9]+):([0-9]+)] ([a-zA-Z]+)\\(\"SIP\\/([a-zA-Z]+)")
			match := re.FindStringSubmatch(scanner.Text())
			if len(match) != 0 {
				if err == nil {
					amounts[fmt.Sprintf("%v_%v", match[4], match[1])] += 1
				}
			}
		}
	}

	for ivrData, amount := range amounts {
		ivrItem := strings.Split(ivrData, "_")
		key, err := strconv.Atoi(ivrItem[1])
		if err != nil {
			continue
		}
		value, _ := types[key]
		data = append(data, models.IVR{
			Ivr:       value,
			Sip:       ivrItem[0],
			Amount:    amount,
			GroupDate: date,
		})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
