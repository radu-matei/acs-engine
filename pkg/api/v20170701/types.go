package v20170701

import (
	"fmt"
	"strings"
)

// ResourcePurchasePlan defines resource plan as required by ARM
// for billing purposes.
type ResourcePurchasePlan struct {
	Name          string `json:"name,omitempty"`
	Product       string `json:"product,omitempty"`
	PromotionCode string `json:"promotionCode,omitempty"`
	Publisher     string `json:"publisher,omitempty"`
}

// ContainerService complies with the ARM model of
// resource definition in a JSON template.
type ContainerService struct {
	ID       string                `json:"id,omitempty"`
	Location string                `json:"location,omitempty"`
	Name     string                `json:"name,omitempty"`
	Plan     *ResourcePurchasePlan `json:"plan,omitempty"`
	Tags     map[string]string     `json:"tags,omitempty"`
	Type     string                `json:"type,omitempty"`

	Properties *Properties `json:"properties"`
}

// Properties represents the ACS cluster definition
type Properties struct {
	ProvisioningState       ProvisioningState        `json:"provisioningState,omitempty"`
	OrchestratorProfile     *OrchestratorProfile     `json:"orchestratorProfile,omitempty"`
	MasterProfile           *MasterProfile           `json:"masterProfile,omitempty"`
	AgentPoolProfiles       []*AgentPoolProfile      `json:"agentPoolProfiles,omitempty"`
	LinuxProfile            *LinuxProfile            `json:"linuxProfile,omitempty"`
	WindowsProfile          *WindowsProfile          `json:"windowsProfile,omitempty"`
	ServicePrincipalProfile *ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
	CustomProfile           *CustomProfile           `json:"customProfile,omitempty"`
}

// ServicePrincipalProfile contains the client and secret used by the cluster for Azure Resource CRUD
// The 'Secret' parameter could be either a plain text, or referenced to a secret in a keyvault.
// In the latter case, the format of the parameter's value should be
// "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>/secrets/<NAME>[/<VERSION>]"
// where:
//    <SUB_ID> is the subscription ID of the keyvault
//    <RG_NAME> is the resource group of the keyvault
//    <KV_NAME> is the name of the keyvault
//    <NAME> is the name of the secret.
//    <VERSION> (optional) is the version of the secret (default: the latest version)
type ServicePrincipalProfile struct {
	ClientID string `json:"clientId,omitempty"`
	Secret   string `json:"secret,omitempty"`
}

// CustomProfile specifies custom properties that are used for
// cluster instantiation.  Should not be used by most users.
type CustomProfile struct {
	Orchestrator string `json:"orchestrator,omitempty"`
}

// LinuxProfile represents the Linux configuration passed to the cluster
type LinuxProfile struct {
	AdminUsername string `json:"adminUsername"`

	SSH struct {
		PublicKeys []struct {
			KeyData string `json:"keyData"`
		} `json:"publicKeys"`
	} `json:"ssh"`
}

// WindowsProfile represents the Windows configuration passed to the cluster
type WindowsProfile struct {
	AdminUsername string `json:"adminUsername,omitempty"`
	AdminPassword string `json:"adminPassword,omitempty"`
}

// ProvisioningState represents the current state of container service resource.
type ProvisioningState string

const (
	// Creating means ContainerService resource is being created.
	Creating ProvisioningState = "Creating"
	// Updating means an existing ContainerService resource is being updated
	Updating ProvisioningState = "Updating"
	// Failed means resource is in failed state
	Failed ProvisioningState = "Failed"
	// Succeeded means resource created succeeded during last create/update
	Succeeded ProvisioningState = "Succeeded"
	// Deleting means resource is in the process of being deleted
	Deleting ProvisioningState = "Deleting"
	// Migrating means resource is being migrated from one subscription or
	// resource group to another
	Migrating ProvisioningState = "Migrating"
)

// OrchestratorProfile contains Orchestrator properties
type OrchestratorProfile struct {
	OrchestratorType    OrchestratorType    `json:"orchestratorType"`
	OrchestratorVersion OrchestratorVersion `json:"orchestratorVersion"`
}

// MasterProfile represents the definition of master cluster
type MasterProfile struct {
	Count                    int    `json:"count"`
	DNSPrefix                string `json:"dnsPrefix"`
	VMSize                   string `json:"vmSize"`
	OSDiskSizeGB             int    `json:"osDiskSizeGB,omitempty"`
	VnetSubnetID             string `json:"vnetSubnetID,omitempty"`
	FirstConsecutiveStaticIP string `json:"firstConsecutiveStaticIP,omitempty"`
	StorageProfile           string `json:"storageProfile,omitempty"`

	// subnet is internal
	subnet string
	// Master LB public endpoint/FQDN with port
	// The format will be FQDN:2376
	// Not used during PUT, returned as part of GET
	FQDN string `json:"fqdn,omitempty"`
}

// AgentPoolProfile represents configuration of VMs running agent
// daemons that register with the master and offer resources to
// host applications in containers.
type AgentPoolProfile struct {
	Name           string `json:"name"`
	Count          int    `json:"count"`
	VMSize         string `json:"vmSize"`
	OSDiskSizeGB   int    `json:"osDiskSizeGB,omitempty"`
	DNSPrefix      string `json:"dnsPrefix"`
	FQDN           string `json:"fqdn"`
	Ports          []int  `json:"ports,omitempty"`
	StorageProfile string `json:"storageProfile"`
	VnetSubnetID   string `json:"vnetSubnetID,omitempty"`
	// OSType is the operating system type for agents
	// Set as nullable to support backward compat because
	// this property was added later.
	// If the value is null or not set, it defaulted to Linux.
	OSType OSType `json:"osType,omitempty"`

	// subnet is internal
	subnet string
}

// OrchestratorType defines orchestrators supported by ACS
type OrchestratorType string

// OrchestratorVersion defines the version for orchestratorType
type OrchestratorVersion string

// UnmarshalText decodes OrchestratorType text, do a case insensitive comparison with
// the defined OrchestratorType constant and set to it if they equal
func (o *OrchestratorType) UnmarshalText(text []byte) error {
	s := string(text)
	switch {
	case strings.EqualFold(s, string(DCOS)):
		*o = DCOS
	case strings.EqualFold(s, string(Kubernetes)):
		*o = Kubernetes
	case strings.EqualFold(s, string(Swarm)):
		*o = Swarm
	case strings.EqualFold(s, string(DockerCE)):
		*o = DockerCE
	default:
		return fmt.Errorf("OrchestratorType has unknown orchestrator: %s", s)
	}

	return nil
}

// OSType represents OS types of agents
type OSType string

// HasWindows returns true if the cluster contains windows
func (a *Properties) HasWindows() bool {
	for _, agentPoolProfile := range a.AgentPoolProfiles {
		if agentPoolProfile.OSType == Windows {
			return true
		}
	}
	return false
}

// IsCustomVNET returns true if the customer brought their own VNET
func (m *MasterProfile) IsCustomVNET() bool {
	return len(m.VnetSubnetID) > 0
}

// GetSubnet returns the read-only subnet for the master
func (m *MasterProfile) GetSubnet() string {
	return m.subnet
}

// SetSubnet sets the read-only subnet for the master
func (m *MasterProfile) SetSubnet(subnet string) {
	m.subnet = subnet
}

// IsManagedDisks returns true if the master specified managed disks
func (m *MasterProfile) IsManagedDisks() bool {
	return m.StorageProfile == ManagedDisks
}

// IsStorageAccount returns true if the master specified storage account
func (m *MasterProfile) IsStorageAccount() bool {
	return m.StorageProfile == StorageAccount
}

// IsCustomVNET returns true if the customer brought their own VNET
func (a *AgentPoolProfile) IsCustomVNET() bool {
	return len(a.VnetSubnetID) > 0
}

// IsWindows returns true if the agent pool is windows
func (a *AgentPoolProfile) IsWindows() bool {
	return a.OSType == Windows
}

// IsLinux returns true if the agent pool is linux
func (a *AgentPoolProfile) IsLinux() bool {
	return a.OSType == Linux
}

// IsManagedDisks returns true if the customer specified managed disks
func (a *AgentPoolProfile) IsManagedDisks() bool {
	return a.StorageProfile == ManagedDisks
}

// IsStorageAccount returns true if the customer specified storage account
func (a *AgentPoolProfile) IsStorageAccount() bool {
	return a.StorageProfile == StorageAccount
}

// GetSubnet returns the read-only subnet for the agent pool
func (a *AgentPoolProfile) GetSubnet() string {
	return a.subnet
}

// SetSubnet sets the read-only subnet for the agent pool
func (a *AgentPoolProfile) SetSubnet(subnet string) {
	a.subnet = subnet
}

// IsSwarmMode returns true if this template is for Docker CE orchestrator
func (o *OrchestratorProfile) IsSwarmMode() bool {
	return o.OrchestratorType == DockerCE
}
