/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func BoolPointer(b bool) *bool {
	return &b
}

func removeMetadata(data map[string]interface{}) {
	for key, value := range data {
		if key == "metadata" {
			delete(data, "metadata")
		} else if nestedMap, ok := value.(map[string]interface{}); ok {
			removeMetadata(nestedMap)
		}
	}
}

// reconcilePVCSize checks if PVCs matching the given name prefix in the namespace
// have the correct size, and if not, expands them to match the desired size
// from the volumeClaimTemplate.
func reconcilePVCSize(ctx context.Context, c client.Client, namespace string, pvcPrefix string, vct *corev1.PersistentVolumeClaimTemplate) error {
	log := log.FromContext(ctx)

	desiredSize, ok := vct.Spec.Resources.Requests[corev1.ResourceStorage]
	if !ok {
		return nil
	}

	pvcList := &corev1.PersistentVolumeClaimList{}
	if err := c.List(ctx, pvcList, client.InNamespace(namespace)); err != nil {
		return fmt.Errorf("unable to list PVCs in namespace %s: %w", namespace, err)
	}

	for i := range pvcList.Items {
		pvc := &pvcList.Items[i]
		if !strings.HasPrefix(pvc.Name, pvcPrefix) {
			continue
		}

		currentSize := pvc.Spec.Resources.Requests[corev1.ResourceStorage]
		if currentSize.Cmp(desiredSize) < 0 {
			log.V(1).Info("Expanding PVC", "pvc", pvc.Name, "namespace", namespace,
				"currentSize", currentSize.String(), "desiredSize", desiredSize.String())
			patch := client.MergeFrom(pvc.DeepCopy())
			pvc.Spec.Resources.Requests[corev1.ResourceStorage] = desiredSize
			if err := c.Patch(ctx, pvc, patch); err != nil {
				return fmt.Errorf("unable to expand PVC %s/%s: %w", namespace, pvc.Name, err)
			}
		}
	}

	return nil
}
