package configReader

import (
	"path/filepath"
	"testing"
)

type ConfigTest struct {
	ConfigStandard                      string
	ConfigWithSpace                     string
	ConfigWithDoubleQuote               string
	ConfigWithSingleQuote               string
	ConfigMultiLine                     string
	ConfigMultiLineWithoutSpace         string
	ConfigStandardFromVar               string
	ConfigStandardFromVarInDoubleQuote  string
	ConfigStandardFromVarInSingleQuote  string
	ConfigNotFoundFromVar               string
	ConfigStandardFromVarInsideBrackets string
	ConfigWithLotsOfTrailingWhitespace  string
	ConfigWithInnerSingleQuote          string
	ConfigWithEscapedDoubleQuote        string
	ConfigWithEscapedSingleQuote        string
	ConfigWithoutEndQuote               string
	ConfigIsBogusDueToPrior             string
	ConfigWithoutMoreQuotes             string
	ConfigInt                           int
}

type ConfigInvalid struct {
	NonMatching     string
	ConfigStandard  int
	ConfigWithSpace string
}

func TestTrimQuotes(t *testing.T) {
	if trimQuotes("\"hello there") != "hello there" && trimQuotes("'howdy") != "howdy" && trimQuotes("howdy'") != "howdy'" &&
		trimQuotes("howdy\"") != "howdy\"" && trimQuotes("\"hello world\"") != "hello world" && trimQuotes("'hello world'") != "hello world" {
		t.Fatal()
	}
}

func TestReadBogusFile(t *testing.T) {
	err := ReadFile("bogus.conf", &ConfigTest{})
	if err == nil {
		t.Fatal("expected failure on bogus filename")
	}
}

func TestReadFile(t *testing.T) {
	path, _ := filepath.Abs("configTest.conf")
	config := &ConfigTest{}
	ReadFile(path, config)
	expect(t, config.ConfigStandard, "value")
	expect(t, config.ConfigWithSpace, "configWithSpace")
	expect(t, config.ConfigWithDoubleQuote, "configWithDoubleQuote")
	expect(t, config.ConfigWithSingleQuote, "configWithSingleQuote")
	expect(t, config.ConfigMultiLine, "config\n  Multi Line")
	expect(t, config.ConfigMultiLineWithoutSpace, "config\nMultiLineWithoutSpace")
	expect(t, config.ConfigStandardFromVar, "value")
	expect(t, config.ConfigStandardFromVarInDoubleQuote, "value")
	expect(t, config.ConfigStandardFromVarInSingleQuote, "$configStandard")
	expect(t, config.ConfigNotFoundFromVar, "$notFound")
	expect(t, config.ConfigStandardFromVarInsideBrackets, "value")
	expect(t, config.ConfigWithLotsOfTrailingWhitespace, "configWithLotsOfTrailingWhitespace")
	expect(t, config.ConfigWithInnerSingleQuote, "configWithInnerSingleQuote'")
	expect(t, config.ConfigWithEscapedDoubleQuote, "configWithEscapedDoubleQuote\\\"")
	expect(t, config.ConfigWithEscapedSingleQuote, "configWithEscapedSingleQuote\\'")
	expect(t, config.ConfigWithoutEndQuote, "configWithoutEndQuote")
	expect(t, config.ConfigIsBogusDueToPrior, "bummer")
	expect(t, config.ConfigWithoutMoreQuotes, "shouldTriggerError")
	if config.ConfigInt != 12345 {
		t.Fatalf("expected config |%i| = |%i|", config.ConfigInt, 12345)
	}
}

func TestReadFileDoesntBlowOnInvalidFields(t *testing.T) {
	path, _ := filepath.Abs("configTest.conf")
	config := &ConfigInvalid{}
	ReadFile(path, config)
	expect(t, config.ConfigWithSpace, "configWithSpace")
}

func expect(t *testing.T, actual string, expected string) {
	if expected != actual {
		t.Fatalf("expected config |%s| = |%s|", expected, actual)
	}
}
