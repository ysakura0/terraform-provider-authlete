package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestClientResourceService_create_import(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: stateSimpleClientState,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authlete_client.client1", "client_id_alias", "terraform_client"),
					resource.TestCheckResourceAttrSet("authlete_client.client1", "client_id"),
					resource.TestCheckResourceAttrSet("authlete_client.client1", "client_secret"),
				),
			},
			{
				ResourceName:            "authlete_client.client1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_id", "client_secret"},
			},
		},
	})
}

func TestClientResourceService_dynamic_services(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: stateDynamicServiceState,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authlete_client.client1", "client_id_alias", "terraform_client"),
					resource.TestCheckResourceAttrSet("authlete_client.client1", "client_id"),
					resource.TestCheckResourceAttrSet("authlete_client.client1", "client_secret"),
				),
			},
		},
	})
}

const stateSimpleClientState = `
provider "authlete" {
}


resource "authlete_client" "client1" {
	developer = "test"
	client_id_alias = "terraform_client"
    client_id_alias_enabled = false
	client_type = "CONFIDENTIAL"
	redirect_uris = [ "https://www.authlete.com/cb" ]
    response_types = [ "CODE" ]
	grant_types = [ "AUTHORIZATION_CODE", "REFRESH_TOKEN" ]
	client_name = "Authlete client"
    requestable_scopes = ["openid", "profile"]
}

`

const stateDynamicServiceState = `
provider "authlete" {
}

resource "authlete_service" "prod" {
  issuer = "https://test.com"
  service_name = "Service for client test"
  supported_grant_types = ["AUTHORIZATION_CODE", "REFRESH_TOKEN"]
  supported_response_types = ["CODE"]
supported_scopes {
	name = "scope1"
    default_entry = false
  }
supported_scopes {
	name = "scope2"
    default_entry = false
  }
}

resource "authlete_client" "client1" {
	apikey = authlete_service.prod.id
	apisecret = authlete_service.prod.api_secret
	developer = "test"
	client_id_alias = "terraform_client"
    client_id_alias_enabled = false
	client_type = "CONFIDENTIAL"
	redirect_uris = [ "https://www.authlete.com/cb" ]
    response_types = [ "CODE" ]
	grant_types = [ "AUTHORIZATION_CODE", "REFRESH_TOKEN" ]
	client_name = "Authlete client"
    requestable_scopes = ["scope1", "scope2"]
}

`