![kubethanos](https://github.com/berkay-dincer/kubethanos/kubethanos.png)

# kubethanos
kubethanos kills half of your pods randomly to engineer chaos in your preferred environment, gives you the opportunity to see how your system behaves under failures. 

## Usage

See the `kubethanos.yaml` file for an example run. Here are the list of valid parameters:

```
--namespaces=!kubesystem,foo-bar // A namespace or a set of namespaces to restrict thanoskube
--included-pod-names=<regex_to_include_pod_names>
--excluded-pod-names=<regex_to_include_pod_names>
--master // The address of the Kubernetes cluster to target, if none looks under $HOME/.kube
--kubeconfig // Path to a kubeconfig file
--interval // Interval between killing pods
--dry-run // If true, print out the pod names without actually killing them.
--debug // Enable debug logging.
```

## Other similar projects

* [chaoskube](https://github.com/linki/chaoskube)
* [kube-monkey](https://github.com/asobti/kube-monkey)
* [PowerfulSeal](https://github.com/bloomberg/powerfulseal)
* [fabric8's chaos monkey](https://fabric8.io/guide/chaosMonkey.html)
* [k8aos](https://github.com/AlexsJones/k8aos)

## Acknowledgements

* Thanks to [@linki](https://github.com/linki) [chaoskube](https://github.com/linki/chaoskube) for giving me the idea and having written something with a broader scope.
