// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	CrossConnectRequiredOnlyResource = CrossConnectResourceDependencies + `
resource "oci_core_cross_connect" "test_cross_connect" {
	#Required
	compartment_id = "${var.compartment_id}"
	location_name = "${var.cross_connect_location_name}"
	port_speed_shape_name = "${var.cross_connect_port_speed_shape_name}"
}
`

	CrossConnectResourceConfig = CrossConnectResourceDependencies + `
resource "oci_core_cross_connect" "test_cross_connect" {
	#Required
	compartment_id = "${var.compartment_id}"
	location_name = "${var.cross_connect_location_name}"
	port_speed_shape_name = "${var.cross_connect_port_speed_shape_name}"

	#Optional
	cross_connect_group_id = "${oci_core_cross_connect_group.test_cross_connect_group.id}"
	display_name = "${var.cross_connect_display_name}"
	#far_cross_connect_or_cross_connect_group_id = "${oci_core_far_cross_connect_or_cross_connect_group.test_far_cross_connect_or_cross_connect_group.id}"
	#near_cross_connect_or_cross_connect_group_id = "${oci_core_near_cross_connect_or_cross_connect_group.test_near_cross_connect_or_cross_connect_group.id}"
	is_active = "${var.cross_connect_is_active}"
}
`
	CrossConnectPropertyVariables = `
variable "cross_connect_display_name" { default = "displayName" }
variable "cross_connect_location_name" { default = "SEA-R1-FAKE-LOCATION" }
variable "cross_connect_port_speed_shape_name" { default = "10 Gbps" }
variable "cross_connect_state" { default = "AVAILABLE" }
variable "cross_connect_is_active" { default = true }

`
	//CrossConnectResourceDependencies = CrossConnectGroupPropertyVariables + CrossConnectGroupResourceConfig + FarCrossConnectOrCrossConnectGroupPropertyVariables + FarCrossConnectOrCrossConnectGroupResourceConfig + NearCrossConnectOrCrossConnectGroupPropertyVariables + NearCrossConnectOrCrossConnectGroupResourceConfig
	CrossConnectResourceDependencies = CrossConnectGroupPropertyVariables + CrossConnectGroupResourceConfig
)

func TestCoreCrossConnectResource_basic(t *testing.T) {
	region := getRequiredEnvSetting("region")
	if strings.ToLower(region) != "r1" {
		t.Skip("FastConnect tests are not yet supported in production regions")
	}

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getRequiredEnvSetting("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_core_cross_connect.test_cross_connect"
	datasourceName := "data.oci_core_cross_connects.test_cross_connects"
	singularDatasourceName := "data.oci_core_cross_connect.test_cross_connect"

	var resId, resId2 string

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify create
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            config + CrossConnectPropertyVariables + compartmentIdVariableStr + CrossConnectRequiredOnlyResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					//resource.TestCheckResourceAttrSet(resourceName, "cross_connect_group_id"),
					//resource.TestCheckResourceAttrSet(resourceName, "far_cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "location_name", "SEA-R1-FAKE-LOCATION"),
					//resource.TestCheckResourceAttrSet(resourceName, "near_cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "port_speed_shape_name", "10 Gbps"),
					resource.TestCheckResourceAttr(resourceName, "state", "PENDING_CUSTOMER"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + CrossConnectResourceDependencies,
			},
			// verify create with optionals
			{
				Config: config + CrossConnectPropertyVariables + compartmentIdVariableStr + CrossConnectResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
					//resource.TestCheckResourceAttrSet(resourceName, "far_cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "location_name", "SEA-R1-FAKE-LOCATION"),
					//resource.TestCheckResourceAttrSet(resourceName, "near_cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "port_speed_shape_name", "10 Gbps"),
					resource.TestCheckResourceAttr(resourceName, "state", "PROVISIONED"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// verify updates to updatable parameters
			{
				Config: config + `
variable "cross_connect_display_name" { default = "displayName2" }
variable "cross_connect_location_name" { default = "SEA-R1-FAKE-LOCATION" }
variable "cross_connect_port_speed_shape_name" { default = "10 Gbps" }
variable "cross_connect_state" { default = "AVAILABLE" }
variable "cross_connect_is_active" { default = true }

                ` + compartmentIdVariableStr + CrossConnectResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
					//resource.TestCheckResourceAttrSet(resourceName, "far_cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "location_name", "SEA-R1-FAKE-LOCATION"),
					//resource.TestCheckResourceAttrSet(resourceName, "near_cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(resourceName, "port_speed_shape_name", "10 Gbps"),
					resource.TestCheckResourceAttr(resourceName, "state", "PROVISIONED"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify datasource
			{
				Config: config + `
variable "cross_connect_display_name" { default = "displayName2" }
variable "cross_connect_location_name" { default = "SEA-R1-FAKE-LOCATION" }
variable "cross_connect_port_speed_shape_name" { default = "10 Gbps" }
variable "cross_connect_state" { default = "AVAILABLE" }
variable "cross_connect_is_active" { default = true }

data "oci_core_cross_connects" "test_cross_connects" {
	#Required
	compartment_id = "${var.compartment_id}"

	#Optional
	cross_connect_group_id = "${oci_core_cross_connect_group.test_cross_connect_group.id}"
	display_name = "${var.cross_connect_display_name}"
	#state = "${var.cross_connect_state}"

    filter {
    	name = "id"
    	values = ["${oci_core_cross_connect.test_cross_connect.id}"]
    }
}
                ` + compartmentIdVariableStr + CrossConnectResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(datasourceName, "cross_connect_group_id"),
					resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
					//resource.TestCheckResourceAttrSet(datasourceName, "far_cross_connect_or_cross_connect_group_id"),
					//resource.TestCheckResourceAttrSet(datasourceName, "near_cross_connect_or_cross_connect_group_id"),
					//resource.TestCheckResourceAttr(datasourceName, "state", "AVAILABLE"),

					resource.TestCheckResourceAttr(datasourceName, "cross_connects.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "cross_connects.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(datasourceName, "cross_connects.0.cross_connect_group_id"),
					resource.TestCheckResourceAttr(datasourceName, "cross_connects.0.display_name", "displayName2"),
					//resource.TestCheckResourceAttrSet(datasourceName, "cross_connects.0.far_cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(datasourceName, "cross_connects.0.location_name", "SEA-R1-FAKE-LOCATION"),
					//resource.TestCheckResourceAttrSet(datasourceName, "cross_connects.0.near_cross_connect_or_cross_connect_group_id"),
					resource.TestCheckResourceAttr(datasourceName, "cross_connects.0.port_speed_shape_name", "10 Gbps"),
					resource.TestCheckResourceAttr(datasourceName, "cross_connects.0.state", "PROVISIONED"),
				),
			},
			// verify singular datasource
			{
				Config: config + `
variable "cross_connect_display_name" { default = "displayName2" }
variable "cross_connect_location_name" { default = "SEA-R1-FAKE-LOCATION" }
variable "cross_connect_port_speed_shape_name" { default = "10 Gbps" }
variable "cross_connect_state" { default = "AVAILABLE" }
variable "cross_connect_is_active" { default = true }

data "oci_core_cross_connect" "test_cross_connect" {
	#Required
	cross_connect_id = "${oci_core_cross_connect.test_cross_connect.id}"
}
                ` + compartmentIdVariableStr + CrossConnectResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(singularDatasourceName, "cross_connect_group_id"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "cross_connect_id"),
					//resource.TestCheckResourceAttrSet(singularDatasourceName, "far_cross_connect_or_cross_connect_group_id"),
					//resource.TestCheckResourceAttrSet(singularDatasourceName, "near_cross_connect_or_cross_connect_group_id"),

					resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
					resource.TestCheckResourceAttr(singularDatasourceName, "location_name", "SEA-R1-FAKE-LOCATION"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "port_name"),
					resource.TestCheckResourceAttr(singularDatasourceName, "port_speed_shape_name", "10 Gbps"),
					resource.TestCheckResourceAttr(singularDatasourceName, "state", "PROVISIONED"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
				),
			},
		},
	})
}
