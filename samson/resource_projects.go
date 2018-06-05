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

	log.Printf("[TOLGA] Creating new project: %#v", project)
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
