# mmReader

mmReader is used to extract file data from websites as bytes, what you choose to do with that data is up to you, that may be require additional login information to access.
This is designed to handle any API calls that may occur automatically for when you need to access across a variety of domains.

Current supported domains are:
- Gitlab
- Github
- S3

# Using mmReader
First use go get to install the latest version of the library.

```$ go get gitlab.com/mmTristan/mmreader```

Then import it into your go program.

```go
import (
    "github.com/mrmxf/mmReader"
)
```

## set up the configuration body to be used in the program
The configuration body searches for keys in a set order until a value is returned. The order is

1. Specified Profile (S3 Only)
2. Enviroment variable
3. Manually parsed keys

For example, if an enviroment variable is found then the manually parsed keys will not be checked or added.

### Manually parse keys

```go
decoder, error := mmreader.AuthInit("s3profilename", key1, key2 key3)
```

### Environment variables
Each key is searched for in the enviroment variables before the manually entered keys are checked. The following profiles have specific variable names that are searched for.

- Gitlab enviroment variables: Gitlab tokens are searched under the $GITLAB_PAT environment variable.
- Github enviroment variables: Gitlab tokens are searched under the $GITHUB_PAT environment variable.
- AWS environment variables: The aws key is searched under $AWS_ACCESS_KEY_ID, the region is searched for under $AWS_DEFAULT_REGION and the secret key is searched for under $AWS_SECRET_ACCESS_KEY .



The keys are automatically added assigned to the relevant domain, e.g. github, s3. Only one of each domain profile can be set up per decoder object.

## get your file

Go get your the bytes of the file from the website you want with the following function.

```go
decoder, error := mmreader.AuthInit("s3profilename", key1, key2 key3)
// Then use the decoder to access the website with any of the added keys
fileBytes, errHttp := decoder.Decode("example.com/pathto/important/file.json")
```


## List of acceptable domain styles

### Gitlab
- https://gitlab.com/api/v4/projects/{project-id-number}/repository/files/path%2Fto%2Ffile.json?ref={branch-name}
- https://gitlab.com/{usernmae}/{repo}/{pathtofile}

### Github
- https://api.github.com/repos/{username}/{repo}/contents/path%2Fto%2Ffile.json
- https://github.com/{username}/{repo}/{filepath}/{filepath}/{file.txt}

### S3
- s3://{bucketname}/{filepath}/{filepath}/{file.txt}
- http://s3.amazonaws.com/{bucketname}/{filepath}/{filepath}/{file.txt}

### http
Any style is acceptable for http but it is expected to follow this layout.

- https://example.com/path/to/name/of/file.json
