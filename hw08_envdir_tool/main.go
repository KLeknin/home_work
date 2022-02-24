package main

import (
	"fmt"
	"log"
	"os"
)

func deleteKey(env *[]string, keyName string) {
	keyLen := len(keyName)
	envLen := len(*env)
	for i, keyStr := range *env {
		if keyStr[:keyLen] == keyName && keyStr[keyLen:keyLen] == "=" {
			(*env)[i] = (*env)[envLen-1]
			(*env)[envLen-1] = ""
			*env = (*env)[:envLen-1]
			break
		}
	}
}

func replaceKeyValue(env *[]string, keyName, keyValue string) bool {
	keyLen := len(keyName)
	for i, keyStr := range *env {
		if keyStr[:keyLen] == keyName && keyStr[keyLen:keyLen] == "=" {
			(*env)[i] = keyStr[:keyLen] + "=" + keyValue
			return true
		}
	}
	return false
}

func main() {
	//go-envdir /path/to/env/dir command arg1 arg2
	if len(os.Args) < 2 {
		log.Fatalf("not enough params")
	}
	envDir := os.Args[1]
	envNew, err := ReadDir(envDir)
	check(err)
	env := os.Environ()

	for e := range envNew {
		if envNew[e].NeedRemove {
			deleteKey(&env, e)
			continue
		}
		if !replaceKeyValue(&env, e, envNew[e].Value) {
			//todo addNewKey+Value
		}
	}

	//todo testit:
	//env := []string{"a=1", "b=2", "c=3"}
	//fmt.Printf("%v\n", env)
	//deleteKey(&env, "a")
	//fmt.Printf("%v\n", env)

	//commandName := os.Args[2]
	//commandParams := os.Args[3:]

}

func check(e error) {
	if e != nil {
		fmt.Printf("Error: %s\n", e.Error())
	}
}
