package main

import (
	"os"
	"log"
	"strings"
	"regexp"
	"fmt"
	"html/template"
)

type Section[T interface{}] struct {
	Title string
	Attributes map[string]string
	Entries []T
}

func printTree(sections []Section[Section[string]]) {
	for _, section := range sections {
		fmt.Println(section.Title)
		for key, val := range section.Attributes {
			fmt.Println("\t"+key+" "+val)
		}

		for _, entry := range section.Entries {
			if (entry.Title != "") {
				fmt.Println("\t"+entry.Title)
				for key, val := range entry.Attributes {
					fmt.Println("\t\t"+key+" "+val)
				}
				for _, detail := range entry.Entries {
					fmt.Println("\t\t* "+detail[0:55]+"...")
				}
			} else {
				for key, val := range entry.Attributes {
					fmt.Println("\t"+key+" "+val)
				}
				for _, detail := range entry.Entries {
					fmt.Println("\t* "+detail[0:55]+"...")
				}
			}
		}
	}
}

func splitLine(line string) (string, string) {
	split := strings.SplitN(line, " ", 2)
	if (len(split) == 1) {
		return split[0], ""
	} else {
		return split[0], split[1]
	}
}

func main() {
	content, err := os.ReadFile("README.md")
	if err != nil { log.Fatal(err) }

	sections := make([]Section[Section[string]], 1)

	for _, line := range strings.Split(string(content), "\n") {
		pre, val := splitLine(line)

		switch pre {
			case "##":
				sections = append(sections, Section[Section[string]]{
					Title: val,
				})
				break

			case "###":
				section := &sections[len(sections)-1]

				if (section.Entries == nil) {
					section.Entries = []Section[string]{Section[string]{
						Title: val,
					}}
				} else {
					section.Entries = append(section.Entries, Section[string]{
						Title: val,
					})
				}
				break

			case "*":
				section := &sections[len(sections)-1]

				if (section.Entries == nil) {
					section.Entries = []Section[string]{Section[string]{
						Entries: []string{val},
					}}
				} else {
					entry := &section.Entries[len(section.Entries)-1]
					entry.Entries = append(entry.Entries, val)
				}
				break

			case "[//]:":
				re := regexp.MustCompile(`\((.*?)\)`)
				matches := re.FindStringSubmatch(val)
				attr := strings.SplitN(matches[1], " ", 2)
				if (len(attr) == 2) {
					section := &sections[len(sections)-1]

					if (section.Entries == nil) {
						if (section.Attributes == nil) {
							 section.Attributes = map[string]string{
							 	attr[0]: attr[1],
							 }
						} else {
							section.Attributes[attr[0]] = attr[1]
						}
					} else {
						entry := &section.Entries[len(section.Entries)-1]
						if (entry.Attributes == nil) {
							 entry.Attributes = map[string]string{
							 	attr[0]: attr[1],
							 }
						} else {
							entry.Attributes[attr[0]] = attr[1]
						}
					}
				}
				break
		}
	}

	// printTree(sections)

	template, err := template.ParseFiles("template.html")
	if err != nil { log.Fatal(err) }

	file, err := os.Create("resume.html")
	if err != nil { log.Fatal(err) }

	err = template.Execute(file, sections)
	if err != nil { log.Fatal(err) }

	err = file.Close()
	if err != nil { log.Fatal(err) }
}
