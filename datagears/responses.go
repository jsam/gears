package datagears

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type Desc struct {
	ID          string `json:"dgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Type        string `json:"type"`
}
type Registration struct {
	Desc             Desc             `json:"desc"`
	ID               string           `json:"id"`
	PrivateData      PrivateData      `json:"privateData"`
	Reader           string           `json:"reader"`
	RegistrationData RegistrationData `json:"registrationData"`
}
type DepsList struct {
	BasePath string   `json:"basePath"`
	Name     string   `json:"name"`
	Wheels   []string `json:"wheels"`
}
type PrivateData struct {
	DepsList  []DepsList `json:"depsList"`
	SessionID string     `json:"sessionId"`
}
type RegistrationData struct {
	Args         map[string]interface{} `json:"args"`
	LastError    interface{}            `json:"lastError"`
	Mode         string                 `json:"mode"`
	NumAborted   int                    `json:"numAborted"`
	NumFailures  int                    `json:"numFailures"`
	NumSuccess   int                    `json:"numSuccess"`
	NumTriggered int                    `json:"numTriggered"`
	Status       string                 `json:"status"`
}

type DumpRegistrationSerializer struct {
	cmd *redis.Cmd
}

func (serializer *DumpRegistrationSerializer) ToMap(data []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	prevKey := ""
	for index, value := range data {
		if valStr, ok := value.(string); ok && index%2 == 0 {
			result[valStr] = nil
			prevKey = valStr
		} else {
			if prevKey == "PD" {
				var privateData map[string]interface{}
				value := strings.Replace(value.(string), "'", "\"", -1)
				err := json.Unmarshal([]byte(value), &privateData)
				if err != nil {
					// TODO: log it
				}

				delete(result, "PD")
				result["privateData"] = privateData
			} else if prevKey == "RegistrationData" {
				delete(result, "RegistrationData")
				result["registrationData"] = serializer.ToMap(value.([]interface{}))
			} else if prevKey == "args" {
				result["args"] = serializer.ToMap(value.([]interface{}))
			} else if prevKey == "desc" {
				var d Desc
				err := json.Unmarshal([]byte(value.(string)), &d)
				if err != nil {
					// TODO:
				}
				result["desc"] = d
			} else {
				result[prevKey] = value
			}
		}
	}

	return result
}

func (serializer *DumpRegistrationSerializer) Parse() ([]Registration, error) {
	result, err := serializer.cmd.Result()
	if err != nil {
		return nil, err
	}

	if resultVal, ok := result.([]interface{}); ok {
		registrations := make([]Registration, 0)
		for _, regData := range resultVal {
			data := serializer.ToMap(regData.([]interface{}))

			reg := Registration{}
			err := mapstructure.Decode(data, &reg)
			if err != nil {
				// TODO: log
			}

			registrations = append(registrations, reg)
		}
		return registrations, nil
	}

	// TODO: return error
	return nil, nil
}
