package generator

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type Data struct {
	Famous   []string
	Before   []string
	After    []string
	Bullshit []string
}

const MAX_LENGTH int = 1000

var data Data

func init() {
	jsonFile, err := os.Open("generator/data.json")
	if err != nil {
		panic(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		panic(err)
	}
}

func shuffle(str []string) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret := make([]string, len(str))
	perm := r.Perm(len(str))
	for i, randIndex := range perm {
		ret[i] = str[randIndex]
	}
	return ret
}

func countSpecial(str string) int {
	chars := [...]string{" ", "，", "。", "？", ";", "！", ":"}
	length := 0
	for _, v := range chars {
		length += strings.Count(str, v)
	}
	return length
}

func canEnd(str string) bool {
	runeStr := []rune(str)
	if len(runeStr) < 2 {
		return false
	}

	if runeStr[len(runeStr)-1] == []rune("。")[0] {
		return true
	} else if runeStr[len(runeStr)-1] == []rune("？")[0] {
		return true
	}

	return false
}

func Generate(topic string, minLen int) string {
	if minLen > MAX_LENGTH {
		minLen = MAX_LENGTH
	}
	shuffledFamous := shuffle(data.Famous)
	shuffledBullshit := shuffle(data.Bullshit)

	rand.Seed(time.Now().UnixNano())
	var ret string
	var hasTopic bool
	indent := strings.Repeat("&nbsp;", 8)

	for utf8.RuneCountInString(ret) < minLen || !canEnd(ret) || !hasTopic {
		x := rand.Intn(100)
		if x < 5 && canEnd(ret) {
			// New paragraph
			ret += "<br><br>" + indent
			minLen += 10
		} else if x < 27 {
			// New famous sentence
			if len(shuffledFamous) == 0 {
				break
			}
			f := shuffledFamous[0]
			shuffledFamous = shuffledFamous[1:]
			before := data.Before[rand.Intn(len(data.Before))]
			after := data.After[rand.Intn(len(data.After))]
			f = strings.ReplaceAll(f, "a", before)
			f = strings.ReplaceAll(f, "b", after)
			minLen += countSpecial(f)
			ret += f
		} else {
			// New bullshit sentence
			if len(shuffledBullshit) == 0 {
				break
			}
			b := shuffledBullshit[0]
			shuffledBullshit = shuffledBullshit[1:]
			if strings.Contains(b, "x") {
				hasTopic = true
			}
			b = strings.ReplaceAll(b, "x", topic)
			minLen += countSpecial(b)
			ret += b
		}
	}
	ret = indent + ret
	return ret
}
