package main

import (
	"flag"
	"github.com/Unknwon/goconfig"
	"log"
	"proxy/service"
	"strconv"
)

var (
	cfgPath = flag.String("config", "", "配置文件")

	sport = flag.Uint("sport", 9090, "代理端口")
	protocol = flag.String("protocol", "http", "代理协议")
	auth = flag.String("auth", "", "代理使用授权")

	adapter = flag.String("adapter", "mysql", "数据库类型")
	host = flag.String("host", "127.0.0.1", "数据库地址")
	port = flag.Uint("port", 3306, "数据库端口")
	user = flag.String("user", "root", "数据库用户名")
	pass = flag.String("pass", "", "数据库密码")
	dbname = flag.String("dbname", "test", "数据库名称")
	charset = flag.String("charset", "utf8mb4", "数据库编码")
)

func main(){
	flag.Parse()

	if *cfgPath != "" {
		config, err := goconfig.LoadConfigFile(*cfgPath)
		if err != nil {
			log.Println(err)
			return
		}

		if section, err := config.GetSection("proxy"); err == nil {
			if v, ok := section["protocol"]; ok {
				*protocol = v
			}
			if v, ok := section["auth"]; ok {
				*auth = v
			}
			if v, ok := section["port"]; ok {
				if p, err := strconv.ParseInt(v, 10, 0); err == nil {
					*sport = uint(p)
				}
			}
		}

		if section, err := config.GetSection("database"); err == nil {
			if v, ok := section["adapter"]; ok {
				*adapter = v
			}
			if v, ok := section["host"]; ok {
				*host = v
			}
			if v, ok := section["user"]; ok {
				*user = v
			}
			if v, ok := section["pass"]; ok {
				*pass = v
			}
			if v, ok := section["dbname"]; ok {
				*dbname = v
			}
			if v, ok := section["charset"]; ok {
				*charset = v
			}
			if v, ok := section["port"]; ok {
				if p, err := strconv.ParseInt(v, 10, 0); err == nil {
					*port = uint(p)
				}
			}
		}
	}

	server := service.CreateServer(sport, protocol, pass)
	if err := server.Run(); err != nil {
		log.Println(err)
	}
}
