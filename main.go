package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	expandFlag = flag.Bool("expand", false, "expand $SHELL_VAR in string values")
	helpFlag   = flag.Bool("h", false, "display this usage")
)

func usage() {
	fmt.Printf("usage: launchwith [-expand] CONFIG CMD ARGS...\n")
	fmt.Printf("  -expand : expand $SHELL_VARS in yaml string values\n")
}

func errExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func envPrefix(pre string, key interface{}) string {
	keyStr := strings.ToUpper(fmt.Sprint(key))
	if pre == "" {
		return keyStr
	}
	return pre + "_" + keyStr
}

func yaml2Env(prefix string, d interface{}) []string {
	switch d1 := d.(type) {
	case map[interface{}]interface{}:
		ret := make([]string, 0, len(d1))
		for k, v := range d1 {
			ret = append(ret, yaml2Env(envPrefix(prefix, k), v)...)
		}
		return ret
	case []interface{}:
		ret := make([]string, 0, len(d1))
		for i, v := range d1 {
			ret = append(ret, yaml2Env(envPrefix(prefix, i), v)...)
		}
		return ret
	case string:
		if prefix == "" {
			errExit(fmt.Errorf("yaml literals are not allowed at top level: %v", d1))
		}
		if *expandFlag {
			return []string{prefix + "=" + os.ExpandEnv(d1)}
		} else {
			return []string{prefix + "=" + d1}
		}
	case int:
		if prefix == "" {
			errExit(fmt.Errorf("yaml literals are not allowed at top level: %v", d1))
		}
		return []string{fmt.Sprintf("%s=%d", prefix, d1)}
	case float64:
		if prefix == "" {
			errExit(fmt.Errorf("yaml literals are not allowed at top level: %v", d1))
		}
		return []string{fmt.Sprintf("%s=%f", prefix, d1)}
	case bool:
		if prefix == "" {
			errExit(fmt.Errorf("yaml literals are not allowed at top level: %v", d1))
		}
		if d1 {
			return []string{prefix + "=1"}
		} else {
			return []string{}
		}
	}
	errExit(fmt.Errorf("unrecognized type: %T", d))
	return nil // unreachable
}

func main() {
	flag.Parse()
	if *helpFlag {
		usage()
		os.Exit(0)
	}
	args := flag.Args()
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}

	conf := func() interface{} {
		confName := args[0]
		confFile, err := os.Open(confName)
		defer confFile.Close()
		errExit(err)
		yamlU := yaml.NewDecoder(confFile)
		var dest interface{}
		err = yamlU.Decode(&dest)
		errExit(err)
		return dest
	}()
	fmt.Printf("YAML\n%v\n", conf)

	env := yaml2Env("", conf)
	fmt.Printf("ENV\n")
	for _, v := range env {
		fmt.Println(v)
	}

	cmd := args[1]
	if len(args) > 2 {
		args = args[2:]
	} else {
		args = args[:0]
	}
	fmt.Printf("CMD: %s %v\n", cmd, args)

	ecmd := exec.Command(cmd, args...)
	ecmd.Env = append(os.Environ(), env...)
	ecmd.Stdin = os.Stdin
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	errExit(ecmd.Run())
}
