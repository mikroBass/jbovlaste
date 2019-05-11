/*
This is a command line program that parses the lojban language dictionary
that anyone can download as an "xml-export.html" file from this URI :
http://jbovlaste.lojban.org/export/xml-export.html?lang=en

This XML content could satisfy a DTD like that :

<!DOCTYPE dictionary [

<!ELEMENT dictionary (direction)>
<!ELEMENT direction (valsi*)>
<!ELEMENT valsi (selmaho?, user, definition, definitionid, notes?, glossword*)>
<!ELEMENT selmaho (#PCDATA)>
<!ELEMENT user (username, realname)>
<!ELEMENT username, (#PCDATA)>
<!ELEMENT realname, (#PCDATA)>
<!ELEMENT definition (#PCDATA)>
<!ELEMENT definitionid (#PCDATA)>
<!ELEMENT notes (#PCDATA)>
<!ELEMENT glossword EMPTY>

<!ATTLIST direction from CDATA #FIXED "lojban">
<!ATTLIST direction to CDATA #REQUIRED>
<!ATTLIST valsi unofficial (true|false) #IMPLIED>
<!ATTLIST valsi word ID #REQUIRED>
<!ATTLIST valsi type CDATA #REQUIRED>
<!ATTLIST glossword word CDATA #REQUIRED>
<!ATTLIST glossword sense CDATA #IMPLIED>

]>

Once compiled you can run this program from a commande line terminal
bu using this syntax :
jbovlaste [-t <type>][-c <clue>][-u <user>][-o][-w][-n] <arg>

where <type>, <clue>, <user> and <arg> are regular expressions
for low case strings, <arg> being a requested argument

[-t] limit the results to entries matching a regexp in words types
[-c] limit the results to entries matching a regexp in definitions or notes
[-u] limit the results to entries matching a regexp in user names
[-o] limit the results to official data (overwrite user flag)
[-w] limit the results to wild data tagged as unofficial
[-n] display additional notes to words descriptions

(for regexp syntax see https://en.wikipedia.org/wiki/Regular_expression)

Author : François Jarriges

License :

		This program is free software: you can redistribute it and/or modify
  	it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"encoding/xml"
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"regexp"
)
type GlossWord struct {
	Word string `xml:"word,attr"`
	Sense string `xml:"sense,attr"`
}
type Valsi struct {
	Unofficial string `xml:"unofficial,attr"`
	Word string `xml:"word,attr"`
	Type string `xml:"type,attr"`
	Selmaho string `xml:"selmaho"`
	UserName string `xml:"user>username"`
	RealName string `xml:"user>realname"`
	Definition string `xml:"definition"`
	DefinitionId string `xml:"definitionid"`
	Notes string `xml:"notes"`
	GlossWords []GlossWord `xml:"glossword"`
}
type Result struct {
	XMLName xml.Name `xml:"dictionary"`
	ValSlice []Valsi `xml:"direction>valsi"`
}

func main() {

	const usageTypeFlag =
		"limit the results to entries matching a regexp in words types\n"
	const usageClueFlag =
		"limit the results to entries matching a regexp in definitions or notes\n"
	const usageUserFlag =
			"limit the results to entries matching a regexp in user names\n"
	const usageOfficialFlag =
		"limit the results to official data (overwrite user flag)\n"
	const usageWildFlag =
		"limit the results to wild data tagged as unofficial\n"
	const usageNotesFlag =
		"display additional notes to words descriptions\n"


	const errorMessage = "\njbovlaste: command failed\n"
	const syntaxMessage =
		"\nusage: jbovlaste [-t <type>][-c <clue>][-u <user>][-o][-w][-n] <arg>\n"+
		"\nwhere <type>, <clue>, <user> and <arg> are regular expressions\n"+
		"for low case strings, <arg> being a requested argument\n"+
		"(for regexp syntax see https://en.wikipedia.org/wiki/Regular_expression)\n"

	var typeRegexpStr, clueRegexpStr, userRegexpStr string
	var officialFlag, wildFlag, notesFlag bool

	flag.StringVar(&typeRegexpStr, "t", "", usageTypeFlag)
	flag.StringVar(&clueRegexpStr, "c", "", usageClueFlag)
	flag.StringVar(&userRegexpStr, "u", "", usageUserFlag)
	flag.BoolVar(&officialFlag, "o", false, usageOfficialFlag)
	flag.BoolVar(&wildFlag, "w", false, usageWildFlag)
	flag.BoolVar(&notesFlag, "n", false, usageNotesFlag)

	flag.Parse()
	args := flag.Args()

	argRegexp := getArgument(args, errorMessage+syntaxMessage)
	if argRegexp == nil {
		flag.PrintDefaults()
		return
	}
	typeRegexp := getRegexp (typeRegexpStr, errorMessage+syntaxMessage)
	clueRegexp := getRegexp (clueRegexpStr, errorMessage+syntaxMessage)

	var userRegexp *regexp.Regexp
	if officialFlag {
		userRegexp = regexp.MustCompile(`official`)
	} else {
		userRegexp = getRegexp (userRegexpStr, errorMessage+syntaxMessage)
	}

	fileName := "/src/jbovlaste/xml-export.html"
	xmlData, err := ioutil.ReadFile(os.Getenv("GOPATH")+fileName)
	checkError(err, errorMessage)

	parsedData := new(Result)

	err = xml.Unmarshal(xmlData, parsedData)
	checkError(err, errorMessage)

	selectedData := selectData(parsedData,
			argRegexp, typeRegexp, userRegexp, clueRegexp, wildFlag)

	displayData(fileName, selectedData, officialFlag, notesFlag)
}

func selectData(parsed *Result,
	argX, typeX, userX, clueX *regexp.Regexp,
	wild bool) (selected []Valsi) {

	for _, valsi := range (*parsed).ValSlice {

		if argX.MatchString(strings.ToLower(valsi.Word)) &&
			(!wild || valsi.Unofficial == "true") &&
			(typeX == nil ||
				typeX.MatchString(strings.ToLower(valsi.Type))) &&
			(userX == nil ||
				userX.MatchString(strings.ToLower(valsi.UserName)) ||
				userX.MatchString(strings.ToLower(valsi.RealName))) {

			if (clueX == nil ||
					 clueX.MatchString(strings.ToLower(valsi.Definition)) ||
					 clueX.MatchString(strings.ToLower(valsi.Notes))) {

				selected = append(selected, valsi)

			} else {
				glossary := valsi.GlossWords

				for _, gloss := range glossary {
					if clueX.MatchString(strings.ToLower(gloss.Word)) ||
						clueX.MatchString(strings.ToLower(gloss.Sense)) {

						selected = append(selected, valsi)
					}
				}
			}
		}
	}
	return selected
}

func displayData(file string, selected []Valsi, official, notes bool) {
	fmt.Printf("Parsing XML lojban dictionary (%v) : "+
		"found %v occurences\n", file, len(selected))

	for _, valsi := range selected {
		fmt.Printf("\n[ %v ]\ttype : %v %v",
			valsi.Word, valsi.Type, valsi.Selmaho)
		if valsi.Unofficial == "true" {
			fmt.Print("\t(unofficial)")
		}
		if official ||
			strings.Contains(strings.ToLower(valsi.UserName), "official") {
			fmt.Print("\t(official)\n")
		} else if strings.Contains(valsi.RealName, valsi.UserName) {
			fmt.Printf("\t(user: %v)\n", valsi.RealName)
		} else {
			fmt.Printf("\t(user: %v / %v)\n",
				valsi.UserName, valsi.RealName)
		}
		fmt.Printf("\ndefinition \t(ID=%v) :\n\t",
			valsi.DefinitionId)
		if txt := valsi.Definition; len(txt) > 0 {
			fmt.Println(formatTxt(txt))
		}
		for g, gloss := range valsi.GlossWords {
			if g == 0 { fmt.Printf("\nglossary :\n") }
			if len(gloss.Sense) >0 {
				fmt.Printf("\t%v\t->\t%v\n", gloss.Word, gloss.Sense)
			} else {
				fmt.Printf("\t%v\n",	gloss.Word)
			}
		}
		if txt := valsi.Notes; notes && len(txt) > 0 {
			fmt.Print("\nnotes :\n", formatTxt(txt), "\n")
		}
		fmt.Println("\n-----------------------------------------------------------")
	}
}

func formatTxt(txt string) string {
	txt = regexp.MustCompile(`\s+`).ReplaceAllString(txt, " ")
	txt = regexp.MustCompile(`([Cc]f\.)\s+`).ReplaceAllString(txt, "$1")
	txt = regexp.MustCompile(`(etc\.)\s+`).ReplaceAllString(txt, "$1")
	txt = regexp.MustCompile(`(e\.g\.)\s+`).ReplaceAllString(txt, "$1")
	txt = regexp.MustCompile(`(\.)\s+`).ReplaceAllString(txt, "$1\n")
	return txt
}

func getArgument(args []string, msg string) *(regexp.Regexp) {
	if len(args) == 1 {
		return getRegexp(args[0], msg)
	} else {
		fmt.Println(msg)
		return nil
	}
}

func getRegexp(str , msg string) *(regexp.Regexp) {
	if str != "" {
		exp, err := regexp.Compile(str)
		checkError(err, msg)
		return exp
	} else {
		return nil
	}
}

func checkError(err error,  msg string) {
	if err != nil {
		fmt.Print(msg)
		panic(err)
	}
}
