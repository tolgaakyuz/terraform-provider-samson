package samson

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	_samson "github.com/tolgaakyuz/samson-go"
)

func resourceSamsonProjects() *schema.Resource {
	return &schema.Resource{
		Create: resourceSamsonProjectCreate,
		Read:   resourceSamsonProjectRead,
		Update: resourceSamsonProjectUpdate,
		Delete: resourceSamsonProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repository_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_variable": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"scope_type_and_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceSamsonProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*_samson.Samson)

	project := &_samson.Project{
		Name:          _samson.String(d.Get("name").(string)),
		RepositoryURL: _samson.String(d.Get("repository_url").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		project.Description = _samson.String(v.(string))
	}

	if v, ok := d.GetOk("environment_variable"); ok {
		environmentVaribles := make([]*_samson.EnvironmentVariable, len(v.([]interface{})))
		for i, d := range v.([]interface{}) {
			environmentVaribles[i] = expandProjectEnvironmentVariable(d)
		}
		project.EnvironmentVariableAttributes = environmentVaribles
	}

	log.Printf("[INFO] Creating new project: %#v", project)

	project, _, err := client.Projects.Upsert(project)
	if err != nil {
		return fmt.Errorf("Error creating samson project error: %s", err)
	}

	d.SetId(strconv.FormatInt(int64(*project.ID), 10))
	return resourceSamsonProjectRead(d, meta)
}

func resourceSamsonProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*_samson.Samson)

	id, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	project, _, err := client.Projects.Get(int(id))

	if err != nil {
		return fmt.Errorf("Error reading samson project %s: %s", d.Id(), err)
	}

	log.Printf("[INFO] samson project read: %#v", project)

	if err := d.Set("name", *project.Name); err != nil {
		return err
	}
	if err := d.Set("repository_url", *project.RepositoryURL); err != nil {
		return err
	}
	if project.Description != nil {
		if err := d.Set("description", *project.Description); err != nil {
			return err
		}
	}
	if len(project.EnvironmentVariableAttributes) > 0 {
		environmentVaribles := make([]interface{}, len(project.EnvironmentVariableAttributes))

		for i, ev := range project.EnvironmentVariableAttributes {
			environmentVaribles[i] = flattenProjectEnvironmentVariable(ev)
		}

		if err := d.Set("environment_variable", environmentVaribles); err != nil {
			return err
		}
	}

	return nil
}

func resourceSamsonProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*_samson.Samson)

	id, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	project := _samson.Project{
		ID:            _samson.Int(int(id)),
		Name:          _samson.String(d.Get("name").(string)),
		RepositoryURL: _samson.String(d.Get("repository_url").(string)),
	}

	if d.HasChange("description") {
		project.Description = _samson.String(d.Get("description").(string))
	}

	log.Printf("[INFO] Updating samson project: %#v", project)

	_, _, updErr := client.Projects.Upsert(&project)
	if updErr != nil {
		return fmt.Errorf("Error updating samson project: %s", updErr)
	}

	log.Printf("[INFO] Updated samson project %#v", project)

	return resourceSamsonProjectRead(d, meta)
}

func resourceSamsonProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*_samson.Samson)

	id, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting project: %d", id)

	_, err = client.Projects.Delete(int(id))
	if err != nil {
		return fmt.Errorf("Error deleting project: %s", err)
	}

	return nil
}

// Expanders

func expandProjectEnvironmentVariable(in interface{}) *_samson.EnvironmentVariable {
	ev := _samson.EnvironmentVariable{}

	m := in.(map[string]interface{})

	if v, ok := m["name"].(string); ok {
		ev.Name = _samson.String(v)
	}
	if v, ok := m["value"].(string); ok {
		ev.Value = _samson.String(v)
	}
	if v, ok := m["scope_type_and_id"].(string); ok {
		ev.ScopeTypeAndID = _samson.String(v)
	}

	return &ev
}

// Flatteners

func flattenProjectEnvironmentVariable(ev *_samson.EnvironmentVariable) interface{} {
	m := make(map[string]interface{}, 0)

	if ev.Name != nil {
		m["name"] = *ev.Name
	}
	if ev.Value != nil {
		m["value"] = *ev.Value
	}
	if ev.ScopeTypeAndID != nil {
		m["scope_type_and_id"] = *ev.ScopeTypeAndID
	}

	return m
}
