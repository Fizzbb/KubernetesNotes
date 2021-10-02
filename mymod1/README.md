# Publish and Use a github go module

## Publish a go module
In a folder named mymod1
1) ```go mod init github.com/fizzbb/kubernetesnotes/mymod1```
2) write code under the package name "mymod1"
3) commit and push code to github
4) tag a release

## Use a go module
1) ```import "github.com/fizzbb/kubernetesnotes/mymod1```
2) use command to download ```go get github.com/fizzbb/kubernetesnotes/mymod1```
3) in the go.mod file, requires will be automatically updated with the downloaded version (tagged release version may not be immediately available)
