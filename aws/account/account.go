package account

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"github.com/otiai10/copy"
	"os/exec"
	"strings"
)

type ReplaceInput struct {
	old string
	new string
}



func GitClone(url string, branch string) string {

	// Get the credentials to download the repository
	username := os.Getenv("GITHUB_USERNAME")
	token := os.Getenv("GITHUB_TOKEN")

	// It will create a random directory name
	dir, err := ioutil.TempDir(".", "")
	if err != nil {
		log.Fatal(err)
	}

	// Clone the given repository to the given directory
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: username, // anything except an empty string
			Password: token,
		},
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Progress:      os.Stdout,
	})
	if err != nil {
		log.Println(err)
	}

	return dir
}

func CopyEnvironment(dir, environment string) (error){
	err := copy.Copy(dir + "/master/us-east-1/dev", dir +  "/master/us-east-1/" + environment)
	if err != nil {
		log.Error(err)
	}
	return err
}

func ReplaceInputs(filename string, replaceInput []ReplaceInput) error {

	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		for _,v := range replaceInput{
			if strings.Contains(line,v.old){
				lines[i] = v.new
			}
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

func TerragruntPlan(workingDirectory string) error {

	cmd := exec.Command(  "terragrunt", "plan-all", "terragrunt-source-update", "terragrunt-include-external-dependencies")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Dir = workingDirectory
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Println("out:", outb.String(), "err:", errb.String())
	return err
}

func TerragruntApply(workingDirectory string) error {

	cmd := exec.Command(  "terragrunt", "plan-all", "terragrunt-source-update", "terragrunt-include-external-dependencies")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Dir = workingDirectory
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Println("out:", outb.String(), "err:", errb.String())
	return err
}


func CreateAccount(name string, environment string) error {
	terraformLiveDir := GitClone("https://github.com/mnsanfilippo/terraform-live.git", "main")

	err := CopyEnvironment(terraformLiveDir,environment)
	if err != nil {
		log.Error(err)
		return err
	}

	filename := terraformLiveDir +  "/master/us-east-1/" + environment + "/accounts/terragrunt.hcl"
	replaceInput := []ReplaceInput{
		{
			old: "account_name",
			new: "  account_name = " + "\"" +  name + "-" + environment + "\"",
		},
		{
			old: "account_email",
			new: "  account_email = " + "\"mnsanfilippo+" +  name + "+" +  environment + "@gmail.com\"",
		},
	}

	err = ReplaceInputs(filename,replaceInput)
	if err != nil {
		log.Error(err)
		return err
	}

	workingDirectory := terraformLiveDir +  "/master/us-east-1/" + environment + "/accounts/"
	err = TerragruntPlan(workingDirectory)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}