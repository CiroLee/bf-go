# bf-go    
> cli **[bf](https://github.com/CiroLee/bf)** implemented by go  

## debug    
```shell
# clone repo
git clone https://github.com/CiroLee/bf-go

# install pkg
cd yourpath/bf-fo
get mod download

# dev
go run xxx/main.go -from en -to zh hello

# output 
# - from: en to: zh
# - 你好

```

ps:     
- different from **[bf](https://github.com/CiroLee/bf)**, command args must be placed first   
- the repo doesn't give config.go file. config is used to save appid and key which are for baidu translate api. you can use your onwn config. the structure follows:      

```shell
# config/config.go

package config

type secrets struct {
	Appid string
	Key   string
}

var Secrets = secrets{
	Appid: "xxxxx",
	Key:   "xxxx",
}
```
