package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	configv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	// TODO(alberto): move to "quay.io/openshift/origin-kubemark-machine-controllers:v4.0.0" once available
	clusterAPIControllerKubemark = "docker.io/gofed/kubemark-machine-controllers:v1.0"
	clusterAPIControllerNoOp     = "no-op"
	kubemarkPlatform             = configv1.PlatformType("kubemark")
)

type Provider string

// OperatorConfig contains configuration for MAO
type OperatorConfig struct {
	TargetNamespace      string `json:"targetNamespace"`
	Controllers          Controllers
	BaremetalControllers BaremetalControllers
}

type Controllers struct {
	Provider           string
	NodeLink           string
	MachineHealthCheck string
}

type BaremetalControllers struct {
	BaremetalOperator         string
	Ironic                    string
	IronicInspector           string
	IronicIpaDownloader       string
	IronicMachineOsDownloader string
	IronicStaticIpManager     string
}

// Images allows build systems to inject images for MAO components
type Images struct {
	MachineAPIOperator            string `json:"machineAPIOperator"`
	ClusterAPIControllerAWS       string `json:"clusterAPIControllerAWS"`
	ClusterAPIControllerOpenStack string `json:"clusterAPIControllerOpenStack"`
	ClusterAPIControllerLibvirt   string `json:"clusterAPIControllerLibvirt"`
	ClusterAPIControllerBareMetal string `json:"clusterAPIControllerBareMetal"`
	ClusterAPIControllerAzure     string `json:"clusterAPIControllerAzure"`
	ClusterAPIControllerGCP       string `json:"clusterAPIControllerGCP"`
	ClusterAPIControllerOvirt     string `json:"clusterAPIControllerOvirt"`
	ClusterAPIControllerVSphere   string `json:"clusterAPIControllerVSphere"`
	// Images required for the metal3 pod
	BaremetalOperator            string `json:"baremetalOperator"`
	BaremetalIronic              string `json:"baremetalIronic"`
	BaremetalIronicInspector     string `json:"baremetalIronicInspector"`
	BaremetalIpaDownloader       string `json:"baremetalIpaDownloader"`
	BaremetalMachineOsDownloader string `json:"baremetalMachineOsDownloader"`
	BaremetalStaticIpManager     string `json:"baremetalStaticIpManager"`
}

func getProviderFromInfrastructure(infra *configv1.Infrastructure) (configv1.PlatformType, error) {
	if infra.Status.Platform == "" {
		return "", fmt.Errorf("no platform provider found on install config")
	}
	return infra.Status.Platform, nil
}

func getProviderFromConfigMap(configmap *corev1.ConfigMap) (configv1.PlatformType, error) {

	provider := configmap.Data["machine-api-provider"]

	switch provider {
	case string(configv1.AWSPlatformType):
		return configv1.AWSPlatformType, nil
	case string(configv1.LibvirtPlatformType):
		return configv1.LibvirtPlatformType, nil
	case string(configv1.OpenStackPlatformType):
		return configv1.OpenStackPlatformType, nil
	case string(configv1.AzurePlatformType):
		return configv1.AzurePlatformType, nil
	case string(configv1.GCPPlatformType):
		return configv1.GCPPlatformType, nil
	case string(configv1.BareMetalPlatformType):
		return configv1.BareMetalPlatformType, nil
	case string(configv1.OvirtPlatformType):
		return configv1.OvirtPlatformType, nil
	case string(configv1.VSpherePlatformType):
		return configv1.BareMetalPlatformType, nil
	case string(kubemarkPlatform):
		return clusterAPIControllerKubemark, nil
	default:
		return clusterAPIControllerNoOp, fmt.Errorf("no platform provider found in config map")
	}
}

func getImagesFromJSONFile(filePath string) (*Images, error) {
	data, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	var i Images
	if err := json.Unmarshal(data, &i); err != nil {
		return nil, err
	}
	return &i, nil
}

func getProviderControllerFromImages(platform configv1.PlatformType, images Images) (string, error) {
	switch platform {
	case configv1.AWSPlatformType:
		return images.ClusterAPIControllerAWS, nil
	case configv1.LibvirtPlatformType:
		return images.ClusterAPIControllerLibvirt, nil
	case configv1.OpenStackPlatformType:
		return images.ClusterAPIControllerOpenStack, nil
	case configv1.AzurePlatformType:
		return images.ClusterAPIControllerAzure, nil
	case configv1.GCPPlatformType:
		return images.ClusterAPIControllerGCP, nil
	case configv1.BareMetalPlatformType:
		return images.ClusterAPIControllerBareMetal, nil
	case configv1.OvirtPlatformType:
		return images.ClusterAPIControllerOvirt, nil
	case configv1.VSpherePlatformType:
		return images.ClusterAPIControllerVSphere, nil
	case kubemarkPlatform:
		return clusterAPIControllerKubemark, nil
	default:
		return clusterAPIControllerNoOp, nil
	}
}

// This function returns images required to bring up the Baremetal Pod.
func newBaremetalControllers(images Images, usingBareMetal bool) BaremetalControllers {
	if !usingBareMetal {
		return BaremetalControllers{}
	}
	return BaremetalControllers{
		BaremetalOperator:         images.BaremetalOperator,
		Ironic:                    images.BaremetalIronic,
		IronicInspector:           images.BaremetalIronicInspector,
		IronicIpaDownloader:       images.BaremetalIpaDownloader,
		IronicMachineOsDownloader: images.BaremetalMachineOsDownloader,
		IronicStaticIpManager:     images.BaremetalStaticIpManager,
	}
}

func getMachineAPIOperatorFromImages(images Images) (string, error) {
	if images.MachineAPIOperator == "" {
		return "", fmt.Errorf("failed gettingMachineAPIOperator image. It is empty")
	}
	return images.MachineAPIOperator, nil
}
