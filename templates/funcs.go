package templates

import (
	"os"
	"text/template"
)

func execute(tpl *template.Template, fileName string, data interface{}) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	err = tpl.Execute(file, data)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func VirtualHostHomePage(fileName string, data interface{}) error {
	tpl, err := template.New("virtual_host_home_page").Parse(virtualHostHomePage)
	if err != nil {
		return err
	}

	return execute(tpl, fileName, data)
}

func VirtualHostProxy(fileName string, data interface{}) error {
	tpl, err := template.New("virtual_host_proxy").Parse(virtualHostProxy)
	if err != nil {
		return err
	}

	return execute(tpl, fileName, data)
}

func HtmlHomePage(fileName string, data interface{}) error {
	tpl, err := template.New("virtual_host_proxy").Parse(htmlHomePage)
	if err != nil {
		return err
	}

	return execute(tpl, fileName, data)
}
