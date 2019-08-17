package nginx

import (
	"github.com/philippe-vandermoere/docker-proxy-go/types/proxy"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

const templateVirtualHostProxyFileName = "/template/virtualHost/proxy.tpl"
const templateVirtualHostHomePageFileName = "/template/virtualHost/homePage.tpl"
const virtualHostHomePageFileName = "default.conf"
const templateHomePageFileName = "/template/homePage.tpl"
const homePageFilename = "index.html"

func Proxy(proxyCollection typeProxy.Collection) error {
	virtualHostDirectory := os.Getenv("VIRTUAL_HOST_DIRECTORY")
	if _, err := os.Stat(virtualHostDirectory); err != nil {
		return err
	}

	homePageDirectory := os.Getenv("HOMEPAGE_DIRECTORY")
	if _, err := os.Stat(homePageDirectory); err != nil {
		return err
	}

	_, filename, _, _ := runtime.Caller(1)
	projectPath := path.Join(path.Dir(filename), "..")

	for _, proxy := range proxyCollection {
		if err := computeTemplate(
			projectPath+templateVirtualHostProxyFileName,
			virtualHostDirectory+"/"+proxy.Domain+".conf",
			proxy,
		); err != nil {
			return err
		}

		log.Info("Created virtual host for domain '" + proxy.Domain + "'.")
	}

	if err := cleanLegacyVirtualHost(proxyCollection, virtualHostDirectory); err != nil {
		return err
	}

	if err := computeTemplate(
		projectPath+templateHomePageFileName,
		homePageDirectory+"/"+homePageFilename,
		proxyCollection,
	); err != nil {
		return err
	}

	data := make(map[string]string)
	data["documentRoot"] = homePageDirectory
	data["index"] = homePageFilename

	if err := computeTemplate(
		projectPath+templateVirtualHostHomePageFileName,
		virtualHostDirectory+"/"+virtualHostHomePageFileName,
		data,
	); err != nil {
		return err
	}

	return nil
}

func cleanLegacyVirtualHost(proxyCollection typeProxy.Collection, virtualHostDirectory string) error {
	files, err := ioutil.ReadDir(virtualHostDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		_, ok := proxyCollection[strings.Trim(file.Name(), ".conf")]
		if !ok && file.Name() != virtualHostHomePageFileName {
			if err := os.Remove(virtualHostDirectory + "/" + file.Name()); err != nil {
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
