# Development

This doc explains how to set up a development environment so you can get started contributing to `KubeFin dashboard`.

## Prerequisites

Follow the instructions below to set up your development environment. Once you meet these requirements, you can make changes and deploy your own version of KubeFin!

### Install requirements

You must install these tools:
1. [`nodejs`](https://nodejs.org/en): React is used to render the UI.
2. [`KubeFin`](https://kubefin.dev/docs/installation/kubectl/): KubeFin needs to be installed in at least one K8s cluster.

## Local development

### Check out your fork

To check out this repository:

1. Create your own
   [fork of this repo](https://help.github.com/articles/fork-a-repo/)
1. Clone it to your machine:
   ```shell
   git clone https://github.com/${YOUR_GITHUB_USERNAME}/kubefin.git
   cd kubefin/dashboard
   git remote add upstream https://github.com/kubefin/kubefin.git
   git remote set-url --push upstream no_push
    ```

_Adding the `upstream` remote sets you up nicely for regularly
[syncing your fork](https://help.github.com/articles/syncing-a-fork/)._

Once you reach this point you are ready to do a full build and deploy as described below.

### Install dependencies

Run following command to install dependencies:
```sh
cd dashboard/web
npm install
```

### Start dashboard

After you've [set up your development environment](#prerequisites), please config the API endpoint in file `kubefin/dashboard/web/src/common/network/http-common.js`. For example, my API endpoint is `http://192.168.1.3:8080`:
```
import axios from "axios";

export default axios.create({
  baseURL: "http://192.168.1.3:8080/api/v1",
});

```

Run following command and `dashboard` can be launched with the code you cloned:
```sh
cd web
PORT=3001 npm start --host 0.0.0.0
```

Now you can access the dashboard at `http://localhost:3001` with your browser(Chrome or Edge).

### Iterating

Once you make changes to the code-base, the page will be refreshed automatically.
