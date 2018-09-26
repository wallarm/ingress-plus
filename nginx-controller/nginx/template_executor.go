package nginx

import (
	"bytes"
	"path"
	"text/template"
)

// TemplateExecutor executes NGINX configuration templates
type TemplateExecutor struct {
	HealthStatus             bool
	mainTemplate             *template.Template
	ingressTemplate          *template.Template
	wallarmTarantoolTemplate *template.Template
}

// NewTemplateExecutor creates a TemplateExecutor
func NewTemplateExecutor(mainTemplatePath string, ingressTemplatePath string, wallarmTarantoolTemplatePath string, healthStatus bool) (*TemplateExecutor, error) {
	// template name must be the base name of the template file https://golang.org/pkg/text/template/#Template.ParseFiles
	nginxTemplate, err := template.New(path.Base(mainTemplatePath)).ParseFiles(mainTemplatePath)
	if err != nil {
		return nil, err
	}

	ingressTemplate, err := template.New(path.Base(ingressTemplatePath)).ParseFiles(ingressTemplatePath)
	if err != nil {
		return nil, err
	}

	wallarmTarantoolTemplate, err := template.New(path.Base(wallarmTarantoolTemplatePath)).ParseFiles(wallarmTarantoolTemplatePath)
	if err != nil {
		return nil, err
	}

	return &TemplateExecutor{
		HealthStatus:             healthStatus,
		mainTemplate:             nginxTemplate,
		ingressTemplate:          ingressTemplate,
		wallarmTarantoolTemplate: wallarmTarantoolTemplate,
	}, nil
}

// UpdateMainTemplate updates the main NGINX template
func (te *TemplateExecutor) UpdateMainTemplate(templateString *string) error {
	newTemplate, err := template.New("nginxTemplate").Parse(*templateString)
	if err != nil {
		return err
	}
	te.mainTemplate = newTemplate

	return nil
}

// UpdateIngressTemplate updates the ingress template
func (te *TemplateExecutor) UpdateIngressTemplate(templateString *string) error {
	newTemplate, err := template.New("ingressTemplate").Parse(*templateString)
	if err != nil {
		return err
	}
	te.ingressTemplate = newTemplate

	return nil
}

// UpdateWallarmTarantoolTemplate updates the Wallarm Tarantool Service template
func (te *TemplateExecutor) UpdateWallarmTarantoolTemplate(templateString *string) error {
	newTemplate, err := template.New("wallarmTarantoolTemplate").Parse(*templateString)
	if err != nil {
		return err
	}
	te.wallarmTarantoolTemplate = newTemplate

	return nil
}

// ExecuteMainConfigTemplate generates the content of the main NGINX configuration file
func (te *TemplateExecutor) ExecuteMainConfigTemplate(cfg *NginxMainConfig) ([]byte, error) {
	cfg.HealthStatus = te.HealthStatus

	var configBuffer bytes.Buffer
	err := te.mainTemplate.Execute(&configBuffer, cfg)

	return configBuffer.Bytes(), err
}

// ExecuteIngressConfigTemplate generates the content of a NGINX configuration file for an Ingress resource
func (te *TemplateExecutor) ExecuteIngressConfigTemplate(cfg *IngressNginxConfig) ([]byte, error) {
	var configBuffer bytes.Buffer
	err := te.ingressTemplate.Execute(&configBuffer, cfg)

	return configBuffer.Bytes(), err
}

// ExecuteIngressConfigTemplate generates the content of a NGINX configuration file for an Ingress resource
func (te *TemplateExecutor) ExecuteWallarmTarantoolTemplate(cfg *WallarmTarantoolConfig) ([]byte, error) {
	var configBuffer bytes.Buffer
	err := te.wallarmTarantoolTemplate.Execute(&configBuffer, cfg)

	return configBuffer.Bytes(), err
}
