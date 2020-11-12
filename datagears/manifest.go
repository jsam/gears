package datagears

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/teris-io/shortid"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

// TODO: rework searching for this with prefix based (e.g. datagears.yaml, datagears.yml)
var manifestNamePrefix = "datagears.yml"

var ErrNotFound = fmt.Errorf("not found")

//DGRemote remote definitions.
type DGRemote struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	// TODO:
	//User     string `json:"user"`
	//Password string `json:"password"`
	Database int `yaml:"database"`
}

//DGGear gear definitions.
type DGGear struct {
	Entrypoint string `yaml:"entrypoint"`

	Description  string            `yaml:"description"`
	Type         string            `yaml:"type"`
	Version      string            `yaml:"version"`
	Blocking     bool              `yaml:"blocking"`
	Values       map[string]string `yaml:"values"`
	Secrets      []string          `yaml:"secrets"`
	Requirements []string          `yaml:"requirements"`

	scriptData   string
	outputData   string
	name         string
	secretValues map[string]string
}

//IsBlocking
func (gear *DGGear) IsBlocking() bool {
	return gear.Blocking
}

func (gear *DGGear) Script() string {
	return gear.outputData
}

func (gear *DGGear) SetName(name string) {
	gear.name = name
}

func (gear *DGGear) ReadSecrets() map[string]string {
	secretValues := make(map[string]string)

	for _, secret := range gear.Secrets {
		secretValues[secret] = os.Getenv(secret)
	}

	return secretValues
}

//Build assembles all needed information for the gear.
func (gear *DGGear) Build() error {
	scriptData, err := ioutil.ReadFile(gear.Entrypoint)
	if err != nil {
		return err
	}

	tmpl, err := template.New("gearTmpl").Parse(string(scriptData))
	if err != nil {
		return err
	}

	if gear.Values == nil {
		gear.Values = make(map[string]string)
	}

	values := gear.Values
	for secretKey, secretValue := range gear.ReadSecrets() {
		if secretValue == "" {
			return fmt.Errorf("secret %s contains no value", secretKey)
		}

		values[secretKey] = secretValue
	}

	var output []byte
	var buf = bytes.NewBuffer(output)
	err = tmpl.Execute(buf, values)
	if err != nil {
		return err
	}
	rended := buf.String()
	lines := strings.Split(rended, "\n")
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		// TODO:
	}

	descPayload := make(map[string]interface{})
	descPayload["dgId"], _ = sid.Generate()
	descPayload["name"] = gear.name
	descPayload["description"] = gear.Description
	descPayload["version"] = gear.Version
	descPayload["type"] = gear.Type
	b, err := json.Marshal(descPayload)

	editedLines := make([]string, 0)
	for _, line := range lines {
		if strings.Contains(line, "GearsBuilder") || strings.Contains(line, "GB") {
			openLine := line[0 : len(line)-1]
			newLine := fmt.Sprintf("%s, desc='%s')", openLine, string(b))
			editedLines = append(editedLines, newLine)
		} else {
			editedLines = append(editedLines, line)
		}
	}

	gear.outputData = strings.Join(editedLines, "\n")
	return nil
}

//DGManifest manifest definition for data gears.
type DGManifest struct {
	Version int                 `yaml:"version"`
	Remotes map[string]DGRemote `yaml:"remotes"`
	Gears   map[string]DGGear   `yaml:"gears"`
}

//GetRemote
func (manifest *DGManifest) GetRemote(remoteName string) (*DGRemote, error) {
	if remote, ok := manifest.Remotes[remoteName]; ok {
		return &remote, nil
	}
	return nil, ErrNotFound
}

//GetGear
func (manifest *DGManifest) GetGear(gearName string) (*DGGear, error) {
	if gear, ok := manifest.Gears[gearName]; ok {
		gear.SetName(gearName)
		return &gear, nil
	}
	return nil, ErrNotFound
}

//NewDGManifest reads the project manifest and returns the DGManifest object.
func NewDGManifest(manifestPath string) *DGManifest {
	manifest := &DGManifest{}

	var path string
	if manifestPath != "" {
		path = manifestPath
	} else {
		path = manifestNamePrefix
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, manifest)
	if err != nil {
		panic(err)
	}

	return manifest
}
