/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package impl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Jeffail/gabs"
	"github.com/go-openapi/loads"
	jsoniter "github.com/json-iterator/go"
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	yaml2 "gopkg.in/yaml.v2"
)

// Directories to be created during init
var dirs = []string{
	utils.InitProjectDefinitions,
	utils.InitProjectImage,
	utils.InitProjectDocs,
	utils.InitProjectSequences,
	utils.InitProjectSequencesFault,
	utils.InitProjectSequencesIn,
	utils.InitProjectSequencesOut,
	utils.InitProjectClientCertificates,
	utils.InitProjectClientCertificates,
	utils.InitProjectInterceptors,
	utils.InitProjectLibs,
}

// InitAPIProject function is used to initlialize an API Project
func InitAPIProject(initCmdOutputDir, initCmdInitialState, initCmdSwaggerPath, initCmdApiDefinitionPath string, isAWSAPI bool) error {
	var dir string
	swaggerSavePath := filepath.Join(initCmdOutputDir, filepath.FromSlash("Definitions/swagger.yaml"))

	if initCmdOutputDir != "" {
		err := os.MkdirAll(initCmdOutputDir, os.ModePerm)
		if err != nil {
			return err
		}
		p, err := filepath.Abs(initCmdOutputDir)
		if err != nil {
			return err
		}
		dir = p
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		dir = pwd
	}
	fmt.Println("Initializing a new WSO2 API Manager project in", dir)

	definitionFile, err := loadDefaultSpecFromDisk()

	// Get the API DTO specific details to process
	def := &definitionFile.Data
	if err != nil {
		return err
	}

	// initCmdInitialState has already validated before creating the 'dir'
	if initCmdInitialState != "" {
		def.LifeCycleStatus = initCmdInitialState
	}

	err = createDirectories(initCmdOutputDir)
	if err != nil {
		return err
	}

	// Use the swagger definition to populate the API definition and save the swagger file separately inside the project
	if initCmdSwaggerPath != "" {
		// Load the swagger file from the provided path
		doc, err := loadSwagger(initCmdSwaggerPath)
		if err != nil {
			return err
		}
		// We use swagger2 loader. It works fine for now
		// Since we don't use 3.0 specific details its ok
		// otherwise please use v2.openAPI3 loaders
		err = v2.Swagger2Populate(def, doc)
		if err != nil {
			return err
		}

		// Convert and write the swagger definition as yaml
		yamlSwagger, err := utils.JsonToYaml(doc.Raw())
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(swaggerSavePath, yamlSwagger, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		// Create an empty swagger definition
		utils.Logln(utils.LogPrefixInfo + "Writing " + swaggerSavePath)
		swaggerDoc, _ := box.Get("/init/swagger-default.yaml")
		err = ioutil.WriteFile(swaggerSavePath, swaggerDoc, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Use the API definition if provided
	if initCmdApiDefinitionPath != "" {
		// Read the API definition file
		utils.Logln(utils.LogPrefixInfo + "Reading API Definition from " + initCmdApiDefinitionPath)
		content, err := ioutil.ReadFile(initCmdApiDefinitionPath)
		if err != nil {
			return err
		}

		apiDef := &v2.APIDefinitionFile{}

		// Substitute env variables to the API definition
		utils.Logln(utils.LogPrefixInfo + "Substituting environment variables")
		data, err := utils.EnvSubstitute(string(content))
		if err != nil {
			return err
		}
		content = []byte(data)

		// Read from yaml definition
		err = yaml2.Unmarshal(content, &apiDef)
		if err != nil {
			return err
		}

		// Marshal original definition
		originalDefBytes, err := jsoniter.Marshal(definitionFile)
		if err != nil {
			return err
		}
		// Marshal new definition
		newDefBytes, err := jsoniter.Marshal(apiDef)
		if err != nil {
			return err
		}

		// Merge two definitions
		finalDefBytes, err := utils.MergeJSON(originalDefBytes, newDefBytes)
		if err != nil {
			return err
		}
		tmpDef := &v2.APIDefinitionFile{}
		err = json.Unmarshal(finalDefBytes, &tmpDef)
		if err != nil {
			return err
		}
		definitionFile.Data = tmpDef.Data
	}

	apiData, err := yaml2.Marshal(definitionFile)
	if err != nil {
		return err
	}

	// Write the API definition to the project directory
	apiJSONPath := filepath.Join(initCmdOutputDir, filepath.FromSlash("api.yaml"))
	utils.Logln(utils.LogPrefixInfo + "Writing " + apiJSONPath)
	err = ioutil.WriteFile(apiJSONPath, apiData, os.ModePerm)
	if err != nil {
		return err
	}

	// Populate the deployment environments configuration and write it to the project directory
	apimProjDeploymentEnvironmentsFilePath := filepath.Join(initCmdOutputDir, "deployment_environments.yaml")
	utils.Logln(utils.LogPrefixInfo + "Writing " + apimProjDeploymentEnvironmentsFilePath)
	deploymentEnvironments, _ := box.Get("/init/default_deployment_environments.yaml")
	err = ioutil.WriteFile(apimProjDeploymentEnvironmentsFilePath, deploymentEnvironments, os.ModePerm)
	if err != nil {
		return err
	}

	// Write the README.md to the project directory
	apimProjReadmeFilePath := filepath.Join(initCmdOutputDir, "README.md")
	utils.Logln(utils.LogPrefixInfo + "Writing " + apimProjReadmeFilePath)
	readme, _ := box.Get("/init/README.md")
	err = ioutil.WriteFile(apimProjReadmeFilePath, readme, os.ModePerm)
	if err != nil {
		return err
	}

	// Create the metaData struct using details from definition
	metaData := utils.MetaData{
		Name:    definitionFile.Data.Name,
		Version: definitionFile.Data.Version,
	}
	marshaledData, err := jsoniter.Marshal(metaData)
	if err != nil {
		return err
	}
	jsonMetaData, err := gabs.ParseJSON(marshaledData)
	metaDataContent, err := utils.JsonToYaml(jsonMetaData.Bytes())
	if err != nil {
		return err
	}

	// Write the api_meta.yaml file to the project directory
	apiMetaDataPath := filepath.Join(initCmdOutputDir, filepath.FromSlash(utils.MetaFileAPI))
	utils.Logln(utils.LogPrefixInfo + "Writing " + apiMetaDataPath)
	err = ioutil.WriteFile(apiMetaDataPath, metaDataContent, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Println("Project initialized")
	fmt.Println("Open README file to learn more")
	return nil
}

// loadDefaultSpecFromDisk loads the API definition stored in HOME/.wso2apictl/default_api.yaml
func loadDefaultSpecFromDisk() (*v2.APIDefinitionFile, error) {
	defaultData, err := ioutil.ReadFile(utils.DefaultAPISpecFilePath)
	if err != nil {
		return nil, err
	}
	def := &v2.APIDefinitionFile{}
	err = yaml.Unmarshal(defaultData, &def)
	if err != nil {
		return nil, err
	}
	return def, nil
}

// createDirectories will create dirs in current working directory
func createDirectories(name string) error {
	for _, dir := range dirs {
		dirPath := filepath.Join(name, filepath.FromSlash(dir))
		utils.Logln(utils.LogPrefixInfo + "Creating directory " + dirPath)
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadSwagger will Load the swagger definition from swaggerDoc
// Swagger2.0/OpenAPI3.0 specs are supported
func loadSwagger(swaggerDoc string) (*loads.Document, error) {
	utils.Logln(utils.LogPrefixInfo + "Loading swagger from " + swaggerDoc)
	return loads.Spec(swaggerDoc)
}
