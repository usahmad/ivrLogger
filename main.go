package main

import (
	db2 "IvrParser/db"
	"IvrParser/db/models"
	"bufio"
	"fmt"
	"github.com/joho/godotenv"
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
	time.Sleep(time.Second * 3)
}

func checkEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		os.Exit(1)
	}
	envs := []string{
		"DB_USERNAME",
		"DB_SCHEMA",
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
	types := getTypes()
	fmt.Println("Connecting to DB...")
	//db := db2.InitDb("")
	for _, file := range scanDir() {
		fmt.Println("Reading file: " + file)
		data := readFile(types, file)
		if len(data) > 0 {
			//fmt.Println("Inserting Data of: " + file)
			//err := models.CreateBulk(db, data)
			//check(err)
			//err = os.Remove(file)
			//if err != nil {
			//	panic(err)
			//}
		} else {
			fmt.Println(fmt.Sprintf("Empty File %v", file))
			err := os.Remove(file)
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Done")
}

func getTypes() map[int]string {
	var items []models.IvrDetail
	db := db2.InitDb("asterisk")
	err := models.GetAll(db, &items)
	check(err)
	details := make(map[int]string)
	for _, val := range items {
		details[val.ID] = val.Description
	}
	return details
}

func readFile(types map[int]string, fileName string) []models.IVR {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	var data []models.IVR
	var amounts map[string]int
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "[s@ivr") {
			re := regexp.MustCompile("@ivr-([0-9]+):([0-9]+)] ([a-zA-Z]+)\\(\"SIP\\/([a-zA-Z]+)")
			match := re.FindStringSubmatch(scanner.Text())
			if len(match) != 0 {
				fmt.Println(match)

				ivr1, err := strconv.Atoi(match[1])
				if err == nil {
					value, _ := types[ivr1]
					fmt.Println(fmt.Sprintf("%v_%v", match[4], value))
					amounts[fmt.Sprintf("%v_%v", match[4], value)] += 1
				}
			}
		}
	}

	//for ivrData, amount := range amounts {
	//	ivrItem := strings.Split(ivrData, "_")
	//	data = append(data, models.IVR{
	//		Ivr:       ivrItem[0],
	//		Sip:       ivrItem[1],
	//		Amount:    amount,
	//		GroupDate: fmt.Sprintf("%s-%s-%d", "01", "01", 1),
	//	})
	//}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	for _, item := range data {
		fmt.Println(item)
	}

	return data
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
