package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Performs necessary upgrade operations if needed
func upgradeConfig() error {
	// read a config file into an interface map, so we can manipulate values without losing any
	configFile := filepath.Join(config.ourBinaryDir, config.ourConfigFilename)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("config file %s does not exist, nothing to upgrade", configFile)
		return nil
	}
	diskConfig := map[string]interface{}{}
	body, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("Couldn't read config file '%s': %s", configFile, err)
		return err
	}

	err = yaml.Unmarshal(body, &diskConfig)
	if err != nil {
		log.Printf("Couldn't parse config file '%s': %s", configFile, err)
		return err
	}

	schemaVersionInterface, ok := diskConfig["schema_version"]
	trace("schemaVersionInterface = %v, ok = %v", schemaVersionInterface, ok)
	if !ok {
		// no schema version, set it to 0
		schemaVersionInterface = 0
	}

	schemaVersion, ok := schemaVersionInterface.(int)
	if !ok {
		err = fmt.Errorf("configuration file contains non-integer schema_version, abort")
		log.Println(err)
		return err
	}

	if schemaVersion == currentSchemaVersion {
		// do nothing
		return nil
	}

	return upgradeConfigSchema(schemaVersion, &diskConfig)
}

// Upgrade from oldVersion to newVersion
func upgradeConfigSchema(oldVersion int, diskConfig *map[string]interface{}) error {
	switch oldVersion {
	case 0:
		err := upgradeSchema0to1(diskConfig)
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("configuration file contains unknown schema_version, abort")
		log.Println(err)
		return err
	}

	configFile := filepath.Join(config.ourBinaryDir, config.ourConfigFilename)
	body, err := yaml.Marshal(diskConfig)
	if err != nil {
		log.Printf("Couldn't generate YAML file: %s", err)
		return err
	}

	err = safeWriteFile(configFile, body)
	if err != nil {
		log.Printf("Couldn't save YAML config: %s", err)
		return err
	}

	return nil
}

func upgradeSchema0to1(diskConfig *map[string]interface{}) error {
	trace("Called")

	// The first schema upgrade:
	// Added "ID" field to "filter" -- we need to populate this field now
	// Added "config.ourDataDir" -- where we will now store filters contents
	for i := range config.Filters {
		filter := &config.Filters[i] // otherwise we will be operating on a copy

		// Set the filter ID
		log.Printf("Seting ID=%d for filter %s", NextFilterId, filter.URL)
		filter.ID = NextFilterId
		NextFilterId++

		// Forcibly update the filter
		_, err := filter.update(true)
		if err != nil {
			log.Fatal(err)
		}

		// Saving it to the filters dir now
		err = filter.save()
		if err != nil {
			log.Fatal(err)
		}
	}

	// No more "dnsfilter.txt", filters are now loaded from config.ourDataDir/filters/
	dnsFilterPath := filepath.Join(config.ourBinaryDir, "dnsfilter.txt")
	_, err := os.Stat(dnsFilterPath)
	if !os.IsNotExist(err) {
		log.Printf("Deleting %s as we don't need it anymore", dnsFilterPath)
		err = os.Remove(dnsFilterPath)
		if err != nil {
			log.Printf("Cannot remove %s due to %s", dnsFilterPath, err)
			// not fatal, move on
		}
	}

	(*diskConfig)["schema_version"] = 1

	return nil
}
