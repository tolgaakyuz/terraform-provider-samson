package samson

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	_samson "github.com/tolgaakyuz/samson-go"
)

func resourceSamsonCommands() *schema.Resource {
	return &schema.Resource{
		Create: resourceSamsonCommandCreate,
		Read:   resourceSamsonCommandRead,
		Update: resourceSamsonCommandUpdate,
		Delete: resourceSamsonCommandDelete,

		Schema: map[string]*schema.Schema{
			"command": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSamsonCommandCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*_samson.Samson)

	command := &_samson.Command{
		Command: _samson.String(d.Get("command").(string)),
	}

	if v, ok := d.GetOk("project_id"); ok {
		command.ProjectID = _samson.String(v.(string))
	}

	log.Printf("[INFO] Creating new command: %#v", command)

	command, _, err := client.Commands.Upsert(command)
	if err != nil {
		return fmt.Errorf("Error creating samson command error: %s", err)
	}

	d.SetId(strconv.FormatInt(int64(*command.ID), 10))
	return resourceSamsonCommandRead(d, meta)
}

func resourceSamsonCommandRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*_samson.Samson)

	id, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	command, _, err := client.Commands.Get(int(id))

	if err != nil {
		return fmt.Errorf("Error reading samson command %s: %s", d.Id(), err)
	}

	log.Printf("[INFO] samson command read: %#v", command)

	if err := d.Set("command", *command.Command); err != nil {
		return err
	}
	if command.ProjectID != nil {
		if err := d.Set("project_id", *command.ProjectID); err != nil {
			return err
		}
	}

	return nil
}

func resourceSamsonCommandUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*_samson.Samson)

	id, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	command := _samson.Command{
		ID:      _samson.Int(int(id)),
		Command: _samson.String(d.Get("command").(string)),
	}

	if d.HasChange("project_id") {
		command.ProjectID = _samson.String(d.Get("project_id").(string))
	}

	log.Printf("[INFO] Updating samson command: %#v", command)

	_, _, updErr := client.Commands.Upsert(&command)
	if updErr != nil {
		return fmt.Errorf("Error updating samson command: %s", updErr)
	}

	log.Printf("[INFO] Updated samson command %#v", command)

	return resourceSamsonCommandRead(d, meta)
}

func resourceSamsonCommandDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*_samson.Samson)

	id, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting command: %d", id)

	_, err = client.Commands.Delete(int(id))
	if err != nil {
		return fmt.Errorf("Error deleting command: %s", err)
	}

	return nil
}
