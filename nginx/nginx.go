package nginx

import (
	"github.com/philippe-vandermoere/docker-proxy-go/types/proxy"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const virtualHostDirectory = "/etc/nginx/conf.d"
const homePageFilename = "/var/www/index.html"
const templateProxyVirtualHostFileName = "/template/proxy.tpl"
const templateHomePageFileName = "/template/homePage.tpl"

func Proxy(proxyCollection typeProxy.Collection) error {
	err := os.MkdirAll(virtualHostDirectory, 755)
	if err != nil {
		return err
	}

	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, proxy := range proxyCollection {
		err := computeTemplate(
			projectPath+templateProxyVirtualHostFileName,
			virtualHostDirectory+"/"+proxy.Domain+".conf",
			proxy,
		)
		if err != nil {
			return err
		}

		log.Info("Created virtual host for domain '" + proxy.Domain + "'.")
	}
	err = cleanLegacyVirtualHost(proxyCollection)
	if err != nil {
		return err
	}

	err = computeTemplate(
		projectPath+templateHomePageFileName,
		homePageFilename,
		proxyCollection,
	)
	if err != nil {
		return err
	}

	return nil
}

func cleanLegacyVirtualHost(proxyCollection typeProxy.Collection) error {
	files, err := ioutil.ReadDir(virtualHostDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		_, ok := proxyCollection[strings.Trim(file.Name(), ".conf")]
		if !ok && file.Name() != "default.conf" {
			err := os.Remove(virtualHostDirectory + "/" + file.Name())
			if err != nil {
				return err
			}

			log.Info("Remove useless virtualHost '" + file.Name() + "'.")
		}
	}

	return nil
}

func computeTemplate(templateFileName string, fileName string, data interface{}) error {
	template, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	err = template.Execute(file, data)
	if err != nil {
		return err
	}

	file.Close()

	return nil
}
