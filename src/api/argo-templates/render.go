/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package argo_templates

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"text/template"

	"github.com/Cray-HPE/cray-nls/src/api/models"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/Masterminds/sprig/v3"
)

//go:embed base/*
var argoWorkflowTemplateFS embed.FS

var validator utils.Validator = utils.NewValidator()

func GetWorkflowTemplate() ([][]byte, error) {
	list, err := fs.Glob(argoWorkflowTemplateFS, "**/*.template.argo.yaml")
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("template: pattern matches no files")
	}
	var filenames []string
	filenames = append(filenames, list...)

	var res [][]byte
	for _, filename := range filenames {
		tmpRes, err := fs.ReadFile(argoWorkflowTemplateFS, filename)
		if err != nil {
			return nil, err
		}
		res = append(res, tmpRes)
	}

	return res, nil
}

func GetWorkerRebuildWorkflow(workerRebuildWorkflowFS fs.FS, createRebuildWorkflowRequest models.CreateRebuildWorkflowRequest) ([]byte, error) {
	err := validator.ValidateWorkerHostnames(createRebuildWorkflowRequest.Hosts)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("worker.rebuild.yaml")

	return GetWorkflow(tmpl, workerRebuildWorkflowFS, createRebuildWorkflowRequest)
}

func GetStorageRebuildWorkflow(storageRebuildWorkflowFS fs.FS, createRebuildWorkflowRequest models.CreateRebuildWorkflowRequest) ([]byte, error) {
	err := validator.ValidateStorageHostnames(createRebuildWorkflowRequest.Hosts)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("storage.rebuild.yaml")

	return GetWorkflow(tmpl, storageRebuildWorkflowFS, createRebuildWorkflowRequest)
}

func GetWorkflow(tmpl *template.Template, workflowFS fs.FS, createRebuildWorkflowRequest models.CreateRebuildWorkflowRequest) ([]byte, error) {
	// add useful helm templating func: include
	var funcMap template.FuncMap = map[string]interface{}{}
	funcMap["include"] = func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	// add sprig templating func
	tmpl, err := tmpl.Funcs(sprig.TxtFuncMap()).Funcs(funcMap).ParseFS(workflowFS, "**/*.yaml")
	if err != nil {
		return nil, err
	}

	var tmpRes bytes.Buffer
	err = tmpl.Execute(&tmpRes, map[string]interface{}{
		"TargetNcns":     createRebuildWorkflowRequest.Hosts,
		"DryRun":         createRebuildWorkflowRequest.DryRun,
		"SwitchPassword": createRebuildWorkflowRequest.SwitchPassword,
		"WipeOsd":        createRebuildWorkflowRequest.WipeOsd,
	})
	if err != nil {
		return nil, err
	}
	return tmpRes.Bytes(), nil
}
