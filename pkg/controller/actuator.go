// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package controller

import (
	"fmt"
	"sync/atomic"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/network"
	gardenerkubernetes "github.com/gardener/gardener/pkg/client/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type actuator struct {
	restConfig *rest.Config
	client     client.Client

	chartRendererFactory extensionscontroller.ChartRendererFactory
	chartApplier         gardenerkubernetes.ChartApplier

	atomicShootWebhookConfig *atomic.Value
	webhookServerPort        int
}

// NewActuator creates a new Actuator that updates the status of the handled Network resources.
func NewActuator(chartRendererFactory extensionscontroller.ChartRendererFactory, shootWebhookConfig *atomic.Value, webhookServerPort int) network.Actuator {
	return &actuator{
		chartRendererFactory:     chartRendererFactory,
		atomicShootWebhookConfig: shootWebhookConfig,
		webhookServerPort:        webhookServerPort,
	}
}

func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

func (a *actuator) InjectConfig(config *rest.Config) error {
	a.restConfig = config

	var err error
	a.chartApplier, err = gardenerkubernetes.NewChartApplierForConfig(config)
	if err != nil {
		return fmt.Errorf("could not create ChartApplier: %w", err)
	}
	return nil
}
