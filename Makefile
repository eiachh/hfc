run-tests:
	export KUBECONFIG=~/.kube/config && \
	env | grep -i kube
	robot test/tests