package main

import (
	"os"
	"log"
	"strings"
	"regexp"
	"fmt"
)

type Section[T interface{}] struct {
	title string
	attributes map[string]string
	entries []T
}

func printTree(sections []Section[Section[string]]) {
	for _, section := range sections {
		fmt.Println(section.title)
		for key, val := range section.attributes {
			fmt.Println("\t"+key+" "+val)
		}

		for _, entry := range section.entries {
			if (entry.title != "") {
				fmt.Println("\t"+entry.title)
				for key, val := range entry.attributes {
					fmt.Println("\t\t"+key+" "+val)
				}
				for _, detail := range entry.entries {
					fmt.Println("\t\t* "+detail[0:55]+"...")
				}
			} else {
				for key, val := range entry.attributes {
					fmt.Println("\t"+key+" "+val)
				}
				for _, detail := range entry.entries {
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
					title: val,
				})
				break

			case "###":
				section := &sections[len(sections)-1]

				if (section.entries == nil) {
					section.entries = []Section[string]{Section[string]{
						title: val,
					}}
				} else {
					section.entries = append(section.entries, Section[string]{
						title: val,
					})
				}
				break

			case "*":
				section := &sections[len(sections)-1]

				if (section.entries == nil) {
					section.entries = []Section[string]{Section[string]{
						entries: []string{val},
					}}
				} else {
					entry := &section.entries[len(section.entries)-1]
					entry.entries = append(entry.entries, val)
				}
				break

			case "[//]:":
				re := regexp.MustCompile(`\((.*?)\)`)
				matches := re.FindStringSubmatch(val)
				attr := strings.SplitN(matches[1], " ", 2)
				if (len(attr) == 2) {
					section := &sections[len(sections)-1]

					if (section.entries == nil) {
						if (section.attributes == nil) {
							 section.attributes = map[string]string{
							 	attr[0]: attr[1],
							 }
						} else {
							section.attributes[attr[0]] = attr[1]
						}
					} else {
						entry := &section.entries[len(section.entries)-1]
						if (entry.attributes == nil) {
							 entry.attributes = map[string]string{
							 	attr[0]: attr[1],
							 }
						} else {
							entry.attributes[attr[0]] = attr[1]
						}
					}
				}
				break
		}
	}

	printTree(sections)
}
