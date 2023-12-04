package credentials

import (
	"os"
	"regexp"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// Decoder is the authorisation object for accessing https locations.
type Decoder struct {
	mut           *sync.Mutex
	authorisation map[string]token
}

// AuthInit generates a map of configuration objects for each of Github, Gitlab and AWS
// based off the available user environments and user input.
// This generates a decoder body which is used to extract bytes from websites.
func AuthInit(s3Profile string, keyparams ...string) (Decoder, error) {
	var err error
	auths := make(map[string]token)
	auths["s3_profile"], err = s3Auth(s3Profile, keyparams)
	auths["git_auth"] = gitlabAuth(keyparams)
	auths["github_auth"] = gitHubAuth(keyparams)

	var sm sync.Mutex
	d := Decoder{authorisation: auths, mut: &sm}

	return d, err
}

// token contains the tokens for any type of system
type token struct {
	tokenCode any
}

// s3AuthDetail contains the required information to access S3
type s3AuthDetail struct {
	credential *credentials.Credentials
	region     string
}

// make a struct of credentials and region
// make variac based on strings of keys
func s3Auth(profile string, keyParams []string) (token, error) {
	var pat token
	var s3Contents s3AuthDetail
	pat.tokenCode = &s3Contents

	//	fmt.Println("searching for aws default profile")
	if profile == "" {
		profile = "default"
	}
	creds := credentials.NewSharedCredentials(config.DefaultSharedCredentialsFilename(), profile)
	if _, err := creds.Get(); err == nil {
		s3Contents.credential = creds

		return pat, nil
	}

	creds = credentials.NewEnvCredentials()
	// fmt.Println("searching for aws environment variables")
	if _, err := creds.Get(); err == nil {
		s3Contents.credential = creds

		return pat, nil
	}
	if len(keyParams) < 2 {
		// fmt.Println("not enough AWS parameters passed")

		return pat, nil
	}

	// fmt.Println("Extracting aws values from user input")
	// Use these regexp to match the string
	reg := regexp.MustCompile(`(us|ap|ca|cn|eu|sa)-(central|(north|south)?(east|west)?)-\d$`)
	sk := regexp.MustCompile(`^[A-Za-z0-9/+=]{40}$`) // These are taken from https://aws.amazon.com/blogs/security/a-safer-way-to-distribute-aws-credentials-to-ec2/
	akid := regexp.MustCompile(`^[A-Z0-9]{20}$`)     // Without the lookarounds as it is not go compatible

	var secret string
	var id string
	for _, text := range keyParams {
		switch {
		case reg.MatchString(text):
			s3Contents.region = text
		case sk.MatchString(text):
			secret = text
		case akid.MatchString(text):
			id = text
		}
	}
	s3Contents.credential = credentials.NewStaticCredentials(id, secret, "")
	// Manually apply items

	return pat, nil
}

// make a struct of credentials and region
// make variac based on strings of keys
func gitlabAuth(keyParams []string) token {
	var pat token
	// Look for more names/ask bruce
	tokenNames := []string{"GITLAB_PAT"}

	for _, env := range tokenNames {
		if val := os.Getenv(env); val != "" {
			pat.tokenCode = val

			return pat
		}
	}

	if len(keyParams) > 0 {
		gitlabT := regexp.MustCompile(`^glpat-[0-9a-zA-Z\\-\\_]{20}$`)
		for _, v := range keyParams {
			if gitlabT.MatchString(v) {
				pat.tokenCode = v

				return pat
			}
		}
	}
	// Else assign an empty token to preven errors
	pat.tokenCode = ""

	return pat
}

// make a struct of credentials and region
// make variac based on strings of keys
func gitHubAuth(keyParams []string) token {
	var pat token
	// Look for more names/ask bruce
	tokenNames := []string{"GITHUB_PAT"}

	for _, env := range tokenNames {
		if val := os.Getenv(env); val != "" {
			pat.tokenCode = val

			return pat
		}
	}

	if len(keyParams) > 0 {
		githubT := regexp.MustCompile(`^ghp_[0-9a-zA-Z]{36}$`)
		for _, v := range keyParams {
			if githubT.MatchString(v) {
				pat.tokenCode = v

				return pat
			}
		}
	}
	// Else assign an empty token to preven errors
	pat.tokenCode = ""

	return pat
}
