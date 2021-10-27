package main

import (
    "fmt"
    "os"
    "path/filepath"
    "os/exec"
    "net/http"
    "encoding/json"
    "bytes"
    flag "github.com/ogier/pflag"
    "strings"
    "io/ioutil"
    "nagacli/langs"
)

// flags
var (
   app  string
   funcName string
   payload string
   gw	string
   registry string
   lang string
   dir string
)

func main() {
	flag.Parse()
	
	 // if user does not supply flags, print usage
	 // we can clean this up later by putting this into its own function
	 /* if flag.NFlag() == 0 {
	     fmt.Printf("Usage: %s [options]\n", os.Args[0])
	     fmt.Println("Options:")
	     flag.PrintDefaults()
	     os.Exit(1)
	  }*/
	  //fmt.Printf("Numbers %d\n", flag.NFlag())	
  	  //fmt.Printf("Command %s\n", flag.Arg(0))
  	  
  	gw = "localhost:8080"
	registry = ""
	if os.Getenv("NAGA_GW")!="" {
	      gw = os.Getenv("NAGA_GW")
	}
	if os.Getenv("NAGA_REGISTRY")!="" {
	    // e.g. some.machine:5000
	    registry = os.Getenv("NAGA_REGISTRY")
	}  	 
  	var err error
	switch flag.Arg(0) {
		case "deploy":
			err = deploy()
		case "upload":
			err = upload()
		case "init":
			err = initialize()
		case "build":
			err = build()
		case "invoke":
			err = invoke()
		default:
			fmt.Printf("'%s' is not a valid command! Commands are: build, deploy, init, version, create\n", flag.Arg(0))
			fmt.Printf("Usage: %s [command] [options]\n", os.Args[0])
			fmt.Println("Options:")
			flag.PrintDefaults()
			os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error while running command! %s", err.Error())			
		os.Exit(1)
	}
  	  
  	os.Exit(0)
}

func deploy() error {
	
	fmt.Printf("GW address %s\n",gw)
	fmt.Printf("Registry address %s\n",registry)
	// Func name is the same as dir name
	// TODO: validate it here and in init()
	dir, err := os.Getwd()
	if err != nil {
	   fmt.Println(err.Error())
           return err
	}
	if (app=="") {
		fmt.Print("No app name specified, 'app' will be used...\n");
		app = "app"
	}
	funcName:= filepath.Base(dir)
	imageName := strings.ToLower(funcName)
	// Add registry, if applicable
	if (registry != "") {
		imageName = registry + "/" + imageName
	}
	fmt.Printf("In deploy of app %s\n",app)
	fmt.Printf("In deploy of func/image name %s %s\n",funcName,imageName)
	// TODO: take version from pom and advance it
	// currently only latest
	fmt.Printf("Building docker image\n")
	//
	cmd := exec.Command("docker", "build", "-t", imageName, ".")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Print(string(stdout))
	// Push the code
	if (registry != "") {
		fmt.Print("\nPushing to registry...\n")
		cmd = exec.Command("docker", "push", imageName)
		stdout, err = cmd.Output()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Print(string(stdout))
	}
	// update gw
	fmt.Printf("Updating Naga GW on %s about the new image\n", gw)
	message := map[string]interface{}{
		"app": app,
		"func": funcName,
		"version": "1",
		"imageName": imageName,
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Post("http://"+gw+"/deploy/add", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Println(err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println(result)
	fmt.Println(result["data"])
	return nil
}

func upload() error {
	fmt.Printf("GW address %s\n",gw)
	fmt.Printf("Registry address %s\n",registry)
	fmt.Printf("uploading app %s func %s dir %s\n",app,funcName,dir)
	
	imageName := strings.ToLower(funcName)
	// Add registry, if applicable
	if (registry != "") {
		imageName = registry + "/" + imageName
	}
	fmt.Printf("In deploy of func/image name %s %s\n",funcName,imageName)
	//
	fmt.Printf("Building docker image\n")
	//
	cmd := exec.Command("docker", "build", "-t", imageName, ".")
	cmd.Dir = dir
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Print(string(stdout))
	// Push the code
	if (registry != "") {
		fmt.Print("\nPushing to registry...\n")
		cmd = exec.Command("docker", "push", imageName)
		stdout, err = cmd.Output()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Print(string(stdout))
	}
	// update gw
	fmt.Printf("Updating Naga GW on %s about the new image\n", gw)
	message := map[string]interface{}{
		"app": app,
		"func": funcName,
		"version": "1",
		"imageName": imageName,
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Post("http://"+gw+"/deploy/add", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Println(err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println(result)
	fmt.Println(result["data"])
	return nil
}

func build() error {
	fmt.Printf("in build")
	return nil
}

func initialize() error {
	fmt.Printf("Lang: %s, App: %s Func: %s \n",lang,app,funcName)
	if (lang!="java") {
		fmt.Print("Unknown lang, java will be used...\n");
	}
	if (app=="") {
		fmt.Print("No app name specified, 'app' will be used...\n");
		app = "app"
	}
	// Create dir if does not exist
	_, err := os.Stat(funcName)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(funcName, 0755)
		errDir = os.MkdirAll(funcName+"/src/main/java/com/naga/"+funcName, 0755)
		errDir = os.MkdirAll(funcName+"/src/main/resources", 0755)
		errDir = os.MkdirAll(funcName+"/src/test/java", 0755)
		errDir = os.MkdirAll(funcName+"/src/test/resources", 0755)
		if errDir != nil {
			return err
		}
	 
	} else {
		fmt.Print("Function folder already exists. Please remove it first (if that's what you want...)\n");
		return nil
	}
	// create files in the right place (Dockerfile, Main.java, pom.xml)
	createFile(funcName+"/Dockerfile",fmt.Sprintf(langs.GetDocker(),funcName,funcName))
	createFile(funcName+"/pom.xml",fmt.Sprintf(langs.GetPOM(),funcName))
	createFile(funcName+"/src/main/java/com/naga/"+funcName+"/Main.java",fmt.Sprintf(langs.GetJava(),funcName))
	return nil
}

func createFile(filename string, content string) error {
    f, err := os.Create(filename)
    if err != nil {
        fmt.Println(err)
        return err
    }
    _, err = f.WriteString(content)
    if err != nil {
        fmt.Println(err)
        f.Close()
        return err
    }
    fmt.Println(filename+" created OK")
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return err
    }
    return nil
}

func invoke() error {
	fmt.Printf("App: %s, Func: %s \n",app,funcName)
	resp, err := http.Post("http://"+gw+"/invoke/"+app+"/"+funcName, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println(err)
	}	
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))	
	return nil
}
	
func init() {
	 flag.StringVarP(&app, "app", "a", "", "The app name")
	 flag.StringVarP(&dir, "dir", "d", "", "The root dir with the lib and dockerfile")
	 flag.StringVarP(&funcName, "fun", "f", "", "The function name")
	 flag.StringVarP(&payload, "payload", "p", "", "The payload")
	 flag.StringVarP(&lang, "lang", "l", "", "The language (currently only java is supported")
}