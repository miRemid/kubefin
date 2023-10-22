# KubeFin Dashboard

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/kubefin/kubefin/blob/main/LICENSE)
[![Releases](https://img.shields.io/github/release/kubefin/kubefin/all.svg)](https://github.com/kubefin/kubefin/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/kubefin/kubefin-agent)](https://hub.docker.com/r/kubefin/kubefin-agent)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/7952/badge)](https://www.bestpractices.dev/en/projects/7952)

## Introduction

KubeFin Dashboard is a general purpose, web-based UI for KubeFin. It enables users to view resource usage and costs from various dimensions across multiple clusters.

<img src="../docs/resources/cost-insights-all-clusters.png" width="45%"> <img src="../docs/resources/cost-insights-single-cluster.png" width="45%">

## Getting Started

To install the latest KubeFin release in primary cluster from the official manifest, run the following command.

```shell
kubectl apply -f https://github.com/kubefin/kubefin/releases/latest/download/kubefin.yaml
```
Once your KubeFin has been installed, wait for the pod to be ready and port forward with:

```shell
kubectl port-forward -nkubefin svc/kubefin-cost-analyzer-service --address='0.0.0.0' 8080 3000
```
To verify that the dashboard and server are running, you may access the KubeFin dashboard at http://localhost:3000.

For more installation method, please refer to the KubeFin documentation.

Community, discussion, contribution, and support

## Documentation

Full documentation is available on the [KubeFin website](https://kubefin.dev).

## Community

We want your contributions and suggestions! One of the easiest ways to contribute is to participate in discussions on the Github Issues/Discussion, chat on IM or the bi-weekly community calls. For more information on the community engagement, developer and contributing guidelines and more, head over to the KubeFin community repo.

## Contact Us

Reach out with any questions you may have and we'll make sure to answer them as soon as possible!

- Slack: [KubeFin Slack](https://kubefin.slack.com)
- Wechat Group (*Chinese*): Broker wechat to add you into the user group.

  <img src="https://kubefin.dev/img/kubefin-assistant.jpg" width="200" />

## Contributing

Check out [DEVELOPMENT](./DEVELOPMENT.md) to see how to develop with KubeFin dashboard.

## Report Vulnerability

Security is a first priority thing for us at KubeFin. If you come across a related issue, please send email to [security@kubefin.dev](security@kubefin.dev).

## Code of Conduct

KubeFin adopts [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).
