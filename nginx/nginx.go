package nginx

import (
	"github.com/philippe-vandermoere/docker-proxy-go/templates"
	typeProxy "github.com/philippe-vandermoere/docker-proxy-go/types/proxy"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

const virtualHostHomePageFileName = "default.conf"
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

	for _, proxy := range proxyCollection {
		if err := templates.VirtualHostProxy(virtualHostDirectory+"/"+proxy.Domain+".conf", proxy); err != nil {
			return err
		}

		log.Info("Created virtual host for domain '" + proxy.Domain + "'.")
	}

	if err := cleanLegacyVirtualHost(proxyCollection, virtualHostDirectory); err != nil {
		return err
	}

	if err := templates.HtmlHomePage(homePageDirectory+"/"+homePageFilename, proxyCollection); err != nil {
		return err
	}

	data := make(map[string]string)
	data["documentRoot"] = homePageDirectory
	data["index"] = homePageFilename

	if err := templates.VirtualHostHomePage(virtualHostDirectory+"/"+virtualHostHomePageFileName, data); err != nil {
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
		_, ok := proxyCollection[strings.TrimSuffix(file.Name(), ".conf")]
		if !ok && file.Name() != virtualHostHomePageFileName {
			if err := os.Remove(virtualHostDirectory + "/" + file.Name()); err != nil {
				return err
			}

			log.Info("Remove useless virtualHost '" + file.Name() + "'.")
		}
	}

	return nil
}
