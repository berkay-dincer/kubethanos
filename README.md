<p align="center"><img src ="https://github.com/berkay-dincer/kubethanos/blob/master/kubethanos.png" width="40%" align="center" alt="chaoskube"></p>

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
--healthcheck // Listens this endpoint for healtcheck
--interval // Interval between killing pods
--dry-run // If true, print out the pod names without actually killing them.
--percentage // percentage of pods to kill. Default is 0.5 
--debug // Enable debug logging.
```

* Configure kubernetes readiness & liveliness probes to `/healthz` endpoint.

## Other similar projects

* [chaosmonkey](https://github.com/Netflix/chaosmonkey)
* [chaoskube](https://github.com/linki/chaoskube)
* [kube-monkey](https://github.com/asobti/kube-monkey)
* [PowerfulSeal](https://github.com/bloomberg/powerfulseal)
* [fabric8's chaos monkey](https://fabric8.io/guide/chaosMonkey.html)
* [k8aos](https://github.com/AlexsJones/k8aos)
* [Cthulhu](https://github.com/xmatters/cthulhu-chaos-testing)

## Acknowledgements

* Thanks to [@linki](https://github.com/linki) [chaoskube](https://github.com/linki/chaoskube) for giving me the idea and having written something with a broader scope.

## Disclaimer

* You are responsible for your actions. If you break things in production while using this software I cannot help you to restore the damage caused.  

## Contribute

Any contributions are welcome!

## License

This project licensed under [MIT](LICENSE).
