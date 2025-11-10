package pkg

import "fmt"

type Schedule struct {
	Days []Day `toml:"days"`
}

type Day struct {
	Name    string   `toml:"name"`
	Lessons []Lesson `toml:"lessons"`
}

type Lesson struct {
	Time     string `toml:"time"`
	Subject  string `toml:"subject"`
	Teacher  string `toml:"teacher"`
	Room     string `toml:"room"`
	Building string `toml:"building"`
	Type     string `toml:"type"`
}

func PrintSchedule(sch Schedule) {
	for i, day := range sch.Days {
		fmt.Printf("\n%d. %s (%d lessons):\n", i+1, day.Name, len(day.Lessons))
		for j, lesson := range day.Lessons {
			fmt.Printf("   %d)%s\n", j+1, lesson.Time)
			fmt.Printf("       %s\n", lesson.Subject)
			fmt.Printf("       %s\n", lesson.Teacher)
			if lesson.Room != "" {
				fmt.Printf("      %s, %s\n", lesson.Room, lesson.Building)
			}
			fmt.Printf("      %s\n", lesson.Type)
		}
	}
}
