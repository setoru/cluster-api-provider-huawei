package ecs

import (
	"bytes"
	"text/template"

	"gopkg.in/yaml.v2"
)

// CloudConfig holds the configuration for the cloud provider.
type CloudConfig struct {
	Region    string
	AccessKey string
	SecretKey string
	VPCID     string
	SubnetID  string
}

// CloudInitConfig represents the structure of a cloud-init configuration.
type CloudInitConfig struct {
	WriteFiles []*WriteFile `yaml:"write_files,omitempty"`
	RunCmd     []string     `yaml:"runcmd,omitempty"`
}

// WriteFile represents a file to be written by cloud-init.
type WriteFile struct {
	Path        string `yaml:"path"`
	Owner       string `yaml:"owner"`
	Permissions string `yaml:"permissions"`
	Content     string `yaml:"content"`
}

func (c *CloudConfig) genCloudProviderSecretTask() (writeFile *WriteFile, runCmd []string, err error) {
	contentTemplate := "[Global]\n  region={{.Region}}\n  access-key={{.AccessKey}}\n  secret-key={{.SecretKey}}\n\n[Vpc]\n  id={{.VPCID}}\n  subnet-id={{.SubnetID}}\n"

	var contentBuffer bytes.Buffer
	tmpl, err := template.New("content").Parse(contentTemplate)
	if err != nil {
		return nil, nil, err
	}
	if err := tmpl.Execute(&contentBuffer, c); err != nil {
		return nil, nil, err
	}

	writeFile = &WriteFile{
		Path:        "/etc/kubernetes/cloud-config",
		Owner:       "root:root",
		Permissions: "0644",
		Content:     contentBuffer.String(),
	}

	runCmd = []string{
		// TODO: remove sleep if we can find a better way to wait for the cluster to be ready
		"sleep 10",
		"export KUBECONFIG=/etc/kubernetes/super-admin.conf",
		"if ! kubectl -n kube-system get secret cloud-config; then kubectl -n kube-system create secret generic cloud-config --from-file=/etc/kubernetes/cloud-config; fi",
		"rm -rf /etc/kubernetes/cloud-config",
	}

	return writeFile, runCmd, nil
}

// appendCloudConfig appends the cloud provider secret to the cloud-init configuration.
func (c *CloudConfig) appendCloudConfig(cloudConfYaml []byte) ([]byte, error) {
	configHeader := "## template: jinja\n#cloud-config\n"

	// Parse the existing cloud-init YAML
	var cloudInitConfig CloudInitConfig
	if err := yaml.Unmarshal(cloudConfYaml, &cloudInitConfig); err != nil {
		return nil, err
	}

	// task1: generate the cloud provider secret
	writeFile, runCmd, err := c.genCloudProviderSecretTask()
	if err != nil {
		return nil, err
	}
	cloudInitConfig.WriteFiles = append(cloudInitConfig.WriteFiles, writeFile)
	cloudInitConfig.RunCmd = append(cloudInitConfig.RunCmd, runCmd...)

	// Marshal the updated cloud-init configuration back to YAML
	updatedCloudConfYaml, err := yaml.Marshal(&cloudInitConfig)
	if err != nil {
		return nil, err
	}

	// Combine configHeader and updatedCloudConfYaml
	var finalYaml bytes.Buffer
	finalYaml.WriteString(configHeader)
	finalYaml.Write(updatedCloudConfYaml)
	return finalYaml.Bytes(), nil
}
