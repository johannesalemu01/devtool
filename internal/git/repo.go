package git

import (
	"os/exec"
	"strings"
)

func DetectRepo() (string,string,error){
out,err:= exec.Command("git","remote","get-url","origin").Output()
 if err !=nil{
	return "","",err
 }

 url:= strings.TrimSpace(string(out))
 url= strings.Replace(url,"https://github.com/","",1)
 url= strings.Replace(url,".git","",1)

 parts:= strings.Split(url,"/")

 return parts[0], parts[1], nil

}