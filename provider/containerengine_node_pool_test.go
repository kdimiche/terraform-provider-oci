// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	NodePoolRequiredOnlyResource = NodePoolResourceDependencies + `
resource "oci_containerengine_node_pool" "test_node_pool" {
	#Required
	cluster_id = "${oci_containerengine_cluster.test_cluster.id}"
	compartment_id = "${var.compartment_id}"
	kubernetes_version = "${var.node_pool_kubernetes_version}"
	name = "${var.node_pool_name}"
	node_image_name = "${var.node_pool_node_image_name}"
	node_shape = "${var.node_pool_node_shape}"
	subnet_ids = ["${oci_core_subnet.nodePool_Subnet_1.id}"] 
}
`

	NodePoolResourceConfig = NodePoolResourceDependencies + `
resource "oci_containerengine_node_pool" "test_node_pool" {
	#Required
	cluster_id = "${oci_containerengine_cluster.test_cluster.id}"
	compartment_id = "${var.compartment_id}"
	kubernetes_version = "${var.node_pool_kubernetes_version}"
	name = "${var.node_pool_name}"
	node_image_name = "${var.node_pool_node_image_name}"
	node_shape = "${var.node_pool_node_shape}"
	subnet_ids = ["${oci_core_subnet.nodePool_Subnet_1.id}"] 

	#Optional
	initial_node_labels {

		#Optional
		key = "${var.node_pool_initial_node_labels_key}"
		value = "${var.node_pool_initial_node_labels_value}"
	}
	quantity_per_subnet = "${var.node_pool_quantity_per_subnet}"
	ssh_public_key = "${var.node_pool_ssh_public_key}"
}
`
	NodePoolPropertyVariables = `
variable "node_pool_initial_node_labels_key" { default = "key" }
variable "node_pool_initial_node_labels_value" { default = "value" }
variable "node_pool_kubernetes_version" { default = "v1.8.11" }
variable "node_pool_name" { default = "name" }
variable "node_pool_node_image_name" { default = "Oracle-Linux-7.4" }
variable "node_pool_node_shape" { default = "VM.Standard1.8" }
variable "node_pool_quantity_per_subnet" { default = 2 }
variable "node_pool_ssh_public_key" { default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample" }

`
	NodePoolResourceDependencies = ClusterPropertyVariables + ClusterResourceConfig + `
resource "oci_core_subnet" "nodePool_Subnet_1" {
	#Required
	availability_domain = "${lookup(data.oci_identity_availability_domains.test_availability_domains.availability_domains[0],"name")}"
	cidr_block = "10.0.22.0/24"
	compartment_id = "${var.compartment_id}"
	vcn_id = "${oci_core_vcn.test_vcn.id}"
    security_list_ids = ["${oci_core_vcn.test_vcn.default_security_list_id}"] # Provider code tries to maintain compatibility with old versions.
	display_name = "tfSubNet1ForNodePool"
}
resource "oci_core_subnet" "nodePool_Subnet_2" {
	#Required
	availability_domain = "${lookup(data.oci_identity_availability_domains.test_availability_domains.availability_domains[0],"name")}"
	cidr_block = "10.0.23.0/24"
	compartment_id = "${var.compartment_id}"
	vcn_id = "${oci_core_vcn.test_vcn.id}"
    security_list_ids = ["${oci_core_vcn.test_vcn.default_security_list_id}"] # Provider code tries to maintain compatibility with old versions.
	display_name = "tfSubNet2ForNodePool"
}`
)

func TestContainerengineNodePoolResource_basic(t *testing.T) {
	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getRequiredEnvSetting("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_containerengine_node_pool.test_node_pool"
	datasourceName := "data.oci_containerengine_node_pools.test_node_pools"

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
				Config:            config + NodePoolPropertyVariables + compartmentIdVariableStr + NodePoolRequiredOnlyResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_version", "v1.8.11"),
					resource.TestCheckResourceAttr(resourceName, "name", "name"),
					resource.TestCheckResourceAttr(resourceName, "node_image_name", "Oracle-Linux-7.4"),
					resource.TestCheckResourceAttr(resourceName, "node_shape", "VM.Standard1.8"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "1"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + NodePoolResourceDependencies,
			},
			// verify create with optionals
			{
				Config: config + NodePoolPropertyVariables + compartmentIdVariableStr + NodePoolResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.0.key", "key"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.0.value", "value"),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_version", "v1.8.11"),
					resource.TestCheckResourceAttr(resourceName, "name", "name"),
					resource.TestCheckResourceAttr(resourceName, "node_image_name", "Oracle-Linux-7.4"),
					resource.TestCheckResourceAttr(resourceName, "node_shape", "VM.Standard1.8"),
					resource.TestCheckResourceAttr(resourceName, "quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(resourceName, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "1"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// verify updates to updatable parameters
			{
				Config: config + `
variable "node_pool_initial_node_labels_key" { default = "key2" }
variable "node_pool_initial_node_labels_value" { default = "value2" }
variable "node_pool_kubernetes_version" { default = "v1.8.11" }
variable "node_pool_name" { default = "name2" }
variable "node_pool_node_image_name" { default = "Oracle-Linux-7.4" }
variable "node_pool_node_shape" { default = "VM.Standard1.8" }
variable "node_pool_quantity_per_subnet" { default = "2" }
variable "node_pool_ssh_public_key" { default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample" }
variable "node_pool_subnet_ids" { default = [] }

                ` + compartmentIdVariableStr + NodePoolResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "cluster_id"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(resourceName, "initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttr(resourceName, "kubernetes_version", "v1.8.11"),
					resource.TestCheckResourceAttr(resourceName, "name", "name2"),
					resource.TestCheckResourceAttr(resourceName, "node_image_name", "Oracle-Linux-7.4"),
					resource.TestCheckResourceAttr(resourceName, "node_shape", "VM.Standard1.8"),
					resource.TestCheckResourceAttr(resourceName, "quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(resourceName, "ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "1"),

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
variable "node_pool_initial_node_labels_key" { default = "key2" }
variable "node_pool_initial_node_labels_value" { default = "value2" }
variable "node_pool_kubernetes_version" { default = "v1.8.11" }
variable "node_pool_name" { default = "name2" }
variable "node_pool_node_image_name" { default = "Oracle-Linux-7.4" }
variable "node_pool_node_shape" { default = "VM.Standard1.8" }
variable "node_pool_quantity_per_subnet" { default = "2" }
variable "node_pool_ssh_public_key" { default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample" }
variable "node_pool_subnet_ids" { default = [] }

data "oci_containerengine_node_pools" "test_node_pools" {
	#Required
	compartment_id = "${var.compartment_id}"

	#Optional
	cluster_id = "${oci_containerengine_cluster.test_cluster.id}"
	name = "${var.node_pool_name}"

    filter {
    	name = "id"
    	values = ["${oci_containerengine_node_pool.test_node_pool.id}"]
    }
}
                ` + compartmentIdVariableStr + NodePoolResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "cluster_id"),
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "name", "name2"),

					resource.TestCheckResourceAttr(datasourceName, "node_pools.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "node_pools.0.cluster_id"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.initial_node_labels.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.initial_node_labels.0.key", "key2"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.initial_node_labels.0.value", "value2"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.kubernetes_version", "v1.8.11"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.name", "name2"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.node_image_name", "Oracle-Linux-7.4"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.node_shape", "VM.Standard1.8"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.quantity_per_subnet", "2"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.ssh_public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDOuBJgh6lTmQvQJ4BA3RCJdSmxRtmiXAQEEIP68/G4gF3XuZdKEYTFeputacmRq9yO5ZnNXgO9akdUgePpf8+CfFtveQxmN5xo3HVCDKxu/70lbMgeu7+wJzrMOlzj+a4zNq2j0Ww2VWMsisJ6eV3bJTnO/9VLGCOC8M9noaOlcKcLgIYy4aDM724MxFX2lgn7o6rVADHRxkvLEXPVqYT4syvYw+8OVSnNgE4MJLxaw8/2K0qp19YlQyiriIXfQpci3ThxwLjymYRPj+kjU1xIxv6qbFQzHR7ds0pSWp1U06cIoKPfCazU9hGWW8yIe/vzfTbWrt2DK6pLwBn/G0x3 sample"),
					resource.TestCheckResourceAttr(datasourceName, "node_pools.0.subnet_ids.#", "1"),
				),
			},
		},
	})
}
