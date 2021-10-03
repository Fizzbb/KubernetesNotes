# Publish and Use a github go module

## Publish a go module
In a folder named mymod1
1) ```go mod init github.com/fizzbb/kubernetesnotes/mymod1```
2) write code under the package name "mymod1"
3) commit and push code to github
4) **TAG A RELEASE** 
   - **Note 1**: go module version format vX.X.X, major, minor, patch,
   - **Note 2**: to tag a module in a repo's sub folder, create release "mymod1/v0.1.2" instead of "v0.1.2", otherwise if tag the upper level repo only, ```go mod tidy```, return **error** 
   ```go get: module github.com/fizzbb/kubernetesnotes@v0.1.1 found, but does not contain package github.com/fizzbb/kubernetesnotes/mymod1```
   - **Note 3**: only tag the top level repo (assume no go.mod)  is fine too. Just the tag won't be applicable for module in the subfoler. When import the module in subfolder, a version starts with ```v0.0.0``` will be automatically generated in the go.mod file

## Use a go module
1) ```import "github.com/fizzbb/kubernetesnotes/mymod1"```
2) use command to download a specific version```go get github.com/fizzbb/kubernetesnotes/mymod1@v0.1.2```, or write ```require github.com/fizzbb/kubernetesnotes/mymod1 v0.1.2``` in the go.mod file, and run ```go mod tidy```
3) If no specific version is required, could just run ```go mod tidy```, in the go.mod file, requires will be automatically updated, if the github module is not tagged a version will be automatically generated, e.g. ```require github.com/fizzbb/kubernetesnotes/mymod1 v0.0.0-20211003175247-3865fcd70cfa```
