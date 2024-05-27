#!/usr/bin/env bash
kubectl --context=kind-bekk apply -f nginx.yaml
kubectl --context=kind-bekk apply -f https://raw.githubusercontent.com/stakater/Reloader/master/deployments/kubernetes/reloader.yaml
