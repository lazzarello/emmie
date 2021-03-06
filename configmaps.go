/*
Copyright (c) 2015, UPMC Enterprises
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name UPMC Enterprises nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL UPMC ENTERPRISES BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
*/

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/pkg/fields"
	"k8s.io/client-go/1.4/pkg/labels"

	"github.com/gorilla/mux"
)

func getConfigMapRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ConfigMapName := vars["ConfigMapName"]
	namespace := vars["namespace"]

	ConfigMap, err := getConfigMap(ConfigMapName, namespace)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(ConfigMap); err != nil {
			panic(err)
		}
	}
}

func getConfigMapsRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]
	namespace := vars["namespace"]

	ConfigMap, err := listConfigMaps(namespace, key, value)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		if len(ConfigMap.Items) > 0 {
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(ConfigMap); err != nil {
				panic(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func listConfigMapsByNamespace(namespace string) (*v1.ConfigMapList, error) {
	list, err := client.Core().ConfigMaps(namespace).List(api.ListOptions{})

	if err != nil {
		log.Println("[listConfigMapsByNamespace] error listing ConfigMaps", err)
		return nil, err
	}

	if len(list.Items) == 0 {
		log.Println("[listConfigMapsByNamespace] No ConfigMaps could be found for namespace!", namespace)
	}

	return list, nil
}

func listConfigMaps(namespace, labelKey, labelValue string) (*v1.ConfigMapList, error) {
	selector := labels.Set{labelKey: labelValue}.AsSelector()
	listOptions := api.ListOptions{FieldSelector: fields.Everything(), LabelSelector: selector}
	list, err := client.Core().ConfigMaps(namespace).List(listOptions)

	if err != nil {
		log.Println("[listConfigMaps] Error listing ConfigMaps", err)
		return nil, err
	}

	if len(list.Items) == 0 {
		log.Println("[listConfigMaps] No ConfigMaps could be found for namespace:", namespace, " labelKey: ", labelKey, " labelValue: ", labelValue)
	}

	return list, nil
}

func getConfigMap(ConfigMapName, namespace string) (*v1.ConfigMap, error) {
	svc, err := client.ConfigMaps(namespace).Get(ConfigMapName)

	if err != nil {
		log.Println("[getConfigMap] Error getting ConfigMap!", err)
		return nil, err
	}

	return svc, nil
}

func createConfigMap(namespace string, ConfigMap *v1.ConfigMap) error {
	_, err := client.Core().ConfigMaps(namespace).Create(ConfigMap)

	if err != nil {
		log.Println("[createConfigMap] Error creating ConfigMap:", err)
	}
	return err
}

func deleteConfigMap(namespace, name string) error {
	// TODO: nil on the DeleteOptions?
	err := client.Core().ConfigMaps(namespace).Delete(name, nil)

	if err != nil {
		log.Println("[deleteConfigMap] Error deleting ConfigMap:", err)
	}
	return err
}
