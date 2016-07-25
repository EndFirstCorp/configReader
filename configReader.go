package configReader

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func ReadFile(path string, outStruct interface{}) error {
	filename := getAbsPath(path)
	values := make(map[string]string)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`(?m)^\s*(\w*)\s*=\s*(["'](?:.|\n)*?["']|.*?)\s*$`)
	matches := re.FindAllStringSubmatch(strings.Replace(string(data), "\r", "", -1), -1)
	for _, group := range matches {
		values[group[1]] = group[2]
	}
	replaceVariables(values)
	mapToStruct(values, outStruct)
	return nil
}

func mapToStruct(valueMap map[string]string, outPtr interface{}) {
	outStruct := reflect.ValueOf(outPtr).Elem()
	for key, value := range valueMap {
		publicName := strings.ToUpper(key[0:1]) + key[1:len(key)]
		field := outStruct.FieldByName(publicName)
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.Atoi(value); err == nil {
				field.SetInt(int64(intVal))
			}
		}
	}
}

func getAbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	absPath := filepath.Join(dir, path)

	// if not in the path of the currently running program, try the working dir
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		dir, _ = os.Getwd()
		absPath = filepath.Join(dir, path)
	}

	return absPath
}

func trimQuotes(value string) string {
	if strings.HasPrefix(value, "\"") || strings.HasPrefix(value, "'") {
		value = value[1:len(value)]
		if strings.HasSuffix(value, "\"") || strings.HasSuffix(value, "'") {
			value = value[0 : len(value)-1]
		}
	}
	return value
}

func replaceVariables(values map[string]string) {
	for key, value := range values {
		// skip if string starts with ' or if there's nothing to replace
		if strings.HasPrefix(value, "'") || !strings.Contains(value, "$") {
			values[key] = trimQuotes(value)
			continue
		}

		for k1, v1 := range values {
			if k1 == key {
				continue
			}
			re := regexp.MustCompile(`\$\{*` + k1 + `\}*`)
			value = re.ReplaceAllString(value, v1)
		}
		values[key] = trimQuotes(value)
	}
}
