// Copyright 2018 Caicloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"fmt"

	s2iconfigmap "github.com/caicloud/ciao/pkg/s2i/configmap"
	pytorchv1beta2 "github.com/kubeflow/pytorch-operator/pkg/apis/pytorch/v1beta2"
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1beta2"
	tfv1beta2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1beta2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/ciao/pkg/types"
)

const (
	baseImageTF      = "tensorflow/tensorflow:1.10.1-py3"
	baseImagePyTorch = "pytorch/pytorch:v0.2"
)

// CM is the type for CM generator.
type CM struct {
	Namespace string
}

// NewCM returns a new CM generator.
func NewCM(namespace string) *CM {
	return &CM{
		Namespace: namespace,
	}
}

// GenerateTFJob generates a new TFJob.
func (c CM) GenerateTFJob(parameter *types.Parameter) (*tfv1beta2.TFJob, error) {
	psCount := int32(parameter.PSCount)
	workerCount := int32(parameter.WorkerCount)
	cleanPodPolicy := common.CleanPodPolicy(parameter.CleanPolicy)

	psResource, err := parameter.Resource.PSLimits()
	if err != nil {
		return nil, err
	}
	workerResource, err := parameter.Resource.WorkerLimits()
	if err != nil {
		return nil, err
	}

	mountPath := fmt.Sprintf("/%s", parameter.Image)
	filename := fmt.Sprintf("/%s/%s", parameter.Image, s2iconfigmap.FileName)

	return &tfv1beta2.TFJob{
		TypeMeta: metav1.TypeMeta{
			Kind: tfv1beta2.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      parameter.GenerateName,
			Namespace: c.Namespace,
		},
		Spec: tfv1beta2.TFJobSpec{
			CleanPodPolicy: &cleanPodPolicy,
			TFReplicaSpecs: map[tfv1beta2.TFReplicaType]*common.ReplicaSpec{
				tfv1beta2.TFReplicaTypePS: {
					Replicas: &psCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  defaultContainerNameTF,
									Image: baseImageTF,
									Command: []string{
										"python",
										filename,
									},
									Resources: v1.ResourceRequirements{
										Limits: psResource,
									},
									VolumeMounts: []v1.VolumeMount{
										{
											Name:      parameter.Image,
											MountPath: mountPath,
										},
									},
								},
							},
							Volumes: []v1.Volume{
								{
									Name: parameter.Image,
									VolumeSource: v1.VolumeSource{
										ConfigMap: &v1.ConfigMapVolumeSource{
											LocalObjectReference: v1.LocalObjectReference{
												Name: parameter.Image,
											},
										},
									},
								},
							},
						},
					},
				},
				tfv1beta2.TFReplicaTypeWorker: {
					Replicas: &workerCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  defaultContainerNameTF,
									Image: baseImageTF,
									Command: []string{
										"python",
										filename,
									},
									Resources: v1.ResourceRequirements{
										Limits: workerResource,
									},
									VolumeMounts: []v1.VolumeMount{
										{
											Name:      parameter.Image,
											MountPath: mountPath,
										},
									},
								},
							},
							Volumes: []v1.Volume{
								{
									Name: parameter.Image,
									VolumeSource: v1.VolumeSource{
										ConfigMap: &v1.ConfigMapVolumeSource{
											LocalObjectReference: v1.LocalObjectReference{
												Name: parameter.Image,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

// GeneratePyTorchJob generates a new PyTorchJob.
func (c CM) GeneratePyTorchJob(parameter *types.Parameter) (*pytorchv1beta2.PyTorchJob, error) {
	masterCount := int32(parameter.MasterCount)
	workerCount := int32(parameter.WorkerCount)
	cleanPodPolicy := common.CleanPodPolicy(parameter.CleanPolicy)

	masterResource, err := parameter.Resource.MasterLimits()
	if err != nil {
		return nil, err
	}
	workerResource, err := parameter.Resource.WorkerLimits()
	if err != nil {
		return nil, err
	}

	mountPath := fmt.Sprintf("/%s", parameter.Image)
	filename := fmt.Sprintf("/%s/%s", parameter.Image, s2iconfigmap.FileName)

	return &pytorchv1beta2.PyTorchJob{
		TypeMeta: metav1.TypeMeta{
			Kind: pytorchv1beta2.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      parameter.GenerateName,
			Namespace: c.Namespace,
		},
		Spec: pytorchv1beta2.PyTorchJobSpec{
			CleanPodPolicy: &cleanPodPolicy,
			PyTorchReplicaSpecs: map[pytorchv1beta2.PyTorchReplicaType]*common.ReplicaSpec{
				pytorchv1beta2.PyTorchReplicaTypeMaster: {
					Replicas: &masterCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  defaultContainerNamePyTorch,
									Image: baseImagePyTorch,
									Command: []string{
										"python",
										filename,
									},
									Resources: v1.ResourceRequirements{
										Limits: masterResource,
									},
									VolumeMounts: []v1.VolumeMount{
										{
											Name:      parameter.Image,
											MountPath: mountPath,
										},
									},
								},
							},
							Volumes: []v1.Volume{
								{
									Name: parameter.Image,
									VolumeSource: v1.VolumeSource{
										ConfigMap: &v1.ConfigMapVolumeSource{
											LocalObjectReference: v1.LocalObjectReference{
												Name: parameter.Image,
											},
										},
									},
								},
							},
						},
					},
				},
				pytorchv1beta2.PyTorchReplicaTypeWorker: {
					Replicas: &workerCount,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  defaultContainerNamePyTorch,
									Image: baseImagePyTorch,
									Command: []string{
										"python",
										filename,
									},
									Resources: v1.ResourceRequirements{
										Limits: workerResource,
									},
									VolumeMounts: []v1.VolumeMount{
										{
											Name:      parameter.Image,
											MountPath: mountPath,
										},
									},
								},
							},
							Volumes: []v1.Volume{
								{
									Name: parameter.Image,
									VolumeSource: v1.VolumeSource{
										ConfigMap: &v1.ConfigMapVolumeSource{
											LocalObjectReference: v1.LocalObjectReference{
												Name: parameter.Image,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}
