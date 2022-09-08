package object

import (
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
)

func SaveYMLFile(yamlObject interface{}, savePath string, perm fs.FileMode) (err error) {

	data, err := yaml.Marshal(&yamlObject)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(savePath, data, perm)
	if err != nil {
		return err
	}

	return err
}

func OpenYMLFile(yamlFile string, yamlObject interface{}) (err error) {

	yamlFileData, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFileData, yamlObject)

	return err
}
