package main

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

// go run .\cmd\conf\main.go
func main() {

	path := "cmd/conf/configs/a.yaml"
	c := config.New(config.WithSource(file.NewSource(path)))
	if err := c.Load(); err != nil {
		panic(err)
	}
	name, err := c.Value("svc.name").String()
	if err != nil {
		panic(err)
	}
	fmt.Println("直接读取配置文件:", name)

	toStruct(c)

}

func toStruct(c config.Config) {
	var v struct {
		Svc struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"svc"`
	}
	if err := c.Scan(&v); err != nil {
		panic(err)
	}
	fmt.Println("直接读取配置文件到struct:", v.Svc.Name)
}
