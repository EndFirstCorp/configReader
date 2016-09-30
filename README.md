# Config Reader
Simple Unix-style config file reader. Handles single-quotes, double-quotes, multi-line, $variable, #comment, etc

## Getting Started
    go get https://github.com/robarchibald/configReader

 1. Create a Go struct that will hold your configuration data (typically mirroring the config parameters you want to access from the config file)
 2. Call `configReader.ReadFile(pathToConfigFile, pointerToYourConfigurationStruct)`
 3. Use your values from the struct (e.g. `fmt.Println(pointerToYourConfigurationStruct.Value1)`)


**Example**
```
type MyConfigStructure struct {
	Value1 string
	Value2 int
}
    
config := &MyConfigStructure{}
configReader.ReadFile(pathToFile, config)
fmt.Println(config.Value1)
```
For additional examples, see configReader_test.go