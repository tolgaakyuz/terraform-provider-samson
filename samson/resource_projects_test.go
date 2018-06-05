package samson

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	_samson "github.com/tolgaakyuz/samson-go"
)

func TestAccSamsonProject_Basic(t *testing.T) {
	var project _samson.Project

	name := acctest.RandString(10)
	description := acctest.RandString(10)
	repositoryURL := fmt.Sprintf("http://%s.com", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckSamsonProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSamsonProjectConfig_basic(name, description, repositoryURL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamsonProjectExists("samson_project.foobar", &project),
					testAccCheckSamsonProjectName(&project, name),
					testAccCheckSamsonProjectDescription(&project, description),
					resource.TestCheckResourceAttr(
						"samson_project.foobar", "name", name),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SAMSON_TOKEN"); v == "" {
		t.Fatal("SAMSON_TOKEN must be set for acceptance tests")
	}
}

func testAccCheckSamsonProjectExists(n string, project *_samson.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project ID is set")
		}

		testAccProvider := Provider().(*schema.Provider)
		client := testAccProvider.Meta().(*_samson.Samson)

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		foundPorject, _, err := client.Projects.Get(int(id))

		if err != nil {
			return err
		}

		if foundPorject.ID == nil || *foundPorject.ID != int(id) {
			return fmt.Errorf("Project not found")
		}

		*project = *foundPorject

		return nil
	}
}

func testAccCheckSamsonProjectName(project *_samson.Project, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if project.Name == nil || *project.Name != name {
			return fmt.Errorf("Bad name: %s", *project.Name)
		}

		return nil
	}
}

func testAccCheckSamsonProjectDescription(project *_samson.Project, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if project.Description == nil || *project.Description != description {
			return fmt.Errorf("Bad description: %s", *project.Description)
		}

		return nil
	}
}

func testAccCheckSamsonProjectDestroy(s *terraform.State) error {
	testAccProvider := Provider().(*schema.Provider)
	client := testAccProvider.Meta().(*_samson.Samson)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "samson_project" {
			continue
		}

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		_, _, err = client.Projects.Get(int(id))

		if err == nil {
			return fmt.Errorf("Project still exists")
		}
	}

	return nil
}

func testAccCheckSamsonProjectConfig_basic(name, description, repositoryURL string) string {
	return fmt.Sprintf(`
provider "samson" {
  token = "123"
}
resource "samson_project" "foobar" {
    name = "%s"
    description = "%s"
	repository_url = "%s"
    attributes {
      runbook_url = "https://www.youtube.com/watch?v=oHg5SJYRHA0"
    }
	environment_variable {
	  name = "env_var_name"
	  value = "env_var_value"
	  scope_type_and_id = "1"
	}
	environment_variable {
	  name = "env_var_name_2"
	  value = "env_var_value_2"
	  scope_type_and_id = "2"
	}
}`, name, description, repositoryURL)
}
