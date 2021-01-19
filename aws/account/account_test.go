package account

import (
	log "github.com/sirupsen/logrus"
	"testing"
)
var accountName = "bootstrap"
var environment = "test"
var terraformLiveDir string



func init(){
	terraformLiveDir = GitClone("https://github.com/mnsanfilippo/terraform-live.git", "main")
}

func TestCopyEnvironment(t *testing.T){
	err := CopyEnvironment(terraformLiveDir,environment)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
}

func TestReplaceInputs(t *testing.T){
	filename := terraformLiveDir +  "/master/us-east-1/" + environment + "/accounts/terragrunt.hcl"

	replaceInput := []ReplaceInput{
		{
			old: "account_name",
			new: "  account_name = " + "\"" + environment + "\"",
		},
		{
			old: "account_email",
			new: "  account_email = " + "\"mnsanfilippo+" + environment + "@gmail.com\"",
		},
	}

	err := ReplaceInputs(filename,replaceInput)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
}

func TestTerragruntPlan(t *testing.T){

	workingDirectory := terraformLiveDir +  "/master/us-east-1/" + environment + "/accounts/"
	err := TerragruntPlan(workingDirectory)
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
}

//func TestTerragruntApply(t *testing.T){
//
//	workingDirectory := terraformLiveDir +  "/master/us-east-1/" + environment + "/accounts/"
//	err := TerragruntApply(workingDirectory)
//	if err != nil {
//		log.Fatal(err)
//		t.Fail()
//	}
//}


func TestCreateAccount(t *testing.T){
	err := CreateAccount(accountName,environment)
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
}

