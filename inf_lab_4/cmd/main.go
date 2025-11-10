package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	pjwl "gitlab.se.ifmo.ru/s503298/inf_lab_4/internal/parserJSONWithLib"
	pr "gitlab.se.ifmo.ru/s503298/inf_lab_4/internal/parserRON"
	st "gitlab.se.ifmo.ru/s503298/inf_lab_4/internal/serializerTOML"
	stwl "gitlab.se.ifmo.ru/s503298/inf_lab_4/internal/serializerTOMLWithLib"
	sx "gitlab.se.ifmo.ru/s503298/inf_lab_4/internal/serializerXML"
	s "gitlab.se.ifmo.ru/s503298/inf_lab_4/pkg"
)

func main() {
	//HandMadeParseAndSerializer("schedule.ron", "schedule.toml")
	// ParserAndSerializerWithLib("schedule.json", "schedule_lib.toml")
	// ParserAndSerializerXML("schedule.ron", "schedule.xml")
	СomparisonNoLibsAndWithLibs()
}

func HandMadeParseAndSerializer(inPath, outPath string) {
	data, err := os.ReadFile(inPath)

	if err != nil {
		fmt.Printf("Error read file: %v\n", err)
		return
	}

	fmt.Println("---START PARSING RON---")

	schedule, err := pr.ParseRON(string(data))
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	fmt.Println("Parsing successful")

	fmt.Printf("\n %d days:\n", len(schedule.Days))
	s.PrintSchedule(*schedule)

	fmt.Println("\n+++END PARSING+++")

	outFile, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	if err := st.WriteTOML(outFile, *schedule); err != nil {
		panic(err)
	}

	fmt.Println("\nWrote TOML to\n", outPath)
}

func ParserAndSerializerWithLib(inPath, outPath string) {
	data, err := os.ReadFile(inPath)
	if err != nil {
		fmt.Printf("Error read gile: %v\n", err)
		return
	}

	fmt.Println("---START PARSING JSON WITH LIB---")
	schedule, err := pjwl.ParseJSONWithLib(data)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	fmt.Printf("\n %d days:\n", len(schedule.Days))
	s.PrintSchedule(*schedule)

	fmt.Println("\n Сериализация в TOML с библиотекой BurntSushi/toml")
	tomlText, err := stwl.SerializeTOMLWithLib(schedule)
	if err != nil {
		fmt.Printf("ERROR serializer: %v\n", err)
		return
	}

	err = os.WriteFile(outPath, []byte(tomlText), 0644)
	if err != nil {
		fmt.Printf("Write file ERROR: %v\n", err)
		return
	}
}

func ParserAndSerializerXML(inPath, outPath string) {
	data, err := os.ReadFile(inPath)

	if err != nil {
		fmt.Printf("Error read file: %v\n", err)
		return
	}

	fmt.Println("---START PARSING RON---")

	schedule, err := pr.ParseRON(string(data))
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	fmt.Println("Parsing successful")

	fmt.Printf("\n %d days:\n", len(schedule.Days))
	s.PrintSchedule(*schedule)

	fmt.Println("\n+++END PARSING+++")

	outFile, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	if err := sx.WriteXML(outFile, *schedule); err != nil {
		panic(err)
	}

	fmt.Println("\nWrote XML to\n", outPath)
}

func СomparisonNoLibsAndWithLibs() {
	data, err := os.ReadFile("schedule.ron")
	if err != nil {
		panic(err)
	}

	startHandMade := time.Now()
	for i := 0; i < 100; i++ {
		schedule, err := pr.ParseRON(string(data))
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		if err := st.WriteTOML(&buf, *schedule); err != nil {
			panic(err)
		}
	}
	elapsedHandMade := time.Since(startHandMade)

	data, err = os.ReadFile("schedule.json")
	if err != nil {
		panic(err)
	}

	startWithLibs := time.Now()
	for i := 0; i < 100; i++ {
		schedule, err := pjwl.ParseJSONWithLib(data)
		if err != nil {
			panic(err)
		}
		_, err = stwl.SerializeTOMLWithLib(schedule)
		if err != nil {
			panic(err)
		}
	}
	elapsedWithLibs := time.Since(startWithLibs)

	fmt.Println("Результаты (100 повторений):")
	fmt.Printf("Hand Made: %v\n", elapsedHandMade)
	fmt.Printf("With Libs : %v\n", elapsedWithLibs)
}
