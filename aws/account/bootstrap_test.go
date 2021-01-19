package account
import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	copy "github.com/otiai10/copy"
)
var accountName = "bootstrap"
var environment = "test"
var terraformLiveDir string

func gitClone(url string, branch string) string {

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

func init(){
	terraformLiveDir = gitClone("git@github.com:mnsanfilippo/terraform-live.git", "main")
}

func TestCopyEnvironment(t *testing.T){
	err := copy.Copy(terraformLiveDir + "/master/us-east-1/dev", terraformLiveDir +  "/master/us-east-1/" + environment)
	if err != nil {
		log.Error(err)
		t.Fail()
	}
}
