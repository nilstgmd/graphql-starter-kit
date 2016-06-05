# Deployment

The deployment of this example project is based on Azure Container Service (ACS) with Mesos cluster management.

## Creating a Resource Group

Assuming you have an Azure account and the CLI installed, get an access token and set the CLI into `arm` mode.

```sh
> azure login -u user@domain.com
Password:

> azure config mode arm
```

To start with ACS you have to create a resource group. A resource group is similar to AWS' CloudFormation stack.
For this example we use a basic group configuration, which is defined in `azuredeploy.json`. This will create a resource
group in North Europe with the name `ace-test`.

```sh
> azure group create -n "acs-test" "deployment-1" -l "North Europe" --template-uri ./deployment/azuredeploy.json
```

During the creating you will be ask for a `dnsPrefixName` and the `sshPublicKey`. The public key needs to be given in the form of
`ssh-rsa AAAAB3Nz..xXSYn myuser@mymachine`.

The deployment will take some time, since it creates all the resources needed for a Mesos cluster. After it is done, you need to connect to the Mesos cluster.

## Connecting to Mesos

The easiest way to work with the newly deployed Mesos cluster is to open a ssh tunnel. That will enable you to work with the Mesos cluster in ACS like it would be a local installation.

```sh
> sudo ssh -i ~/.ssh/acs_key -L 80:localhost:80 -f -N azureuser@YOUR_DNS_PREFIXmgmt.northeurope.cloudapp.azure.com -p 2200
```

That will forward port `80` of the Mesos master to `localhost` port `80`. Make sure that there is no application running on port `80`, you use `sudo` for reserved ports and to include the `-i` option pointing to the matching private key.

If successful you can connect to Mesos, Marathon and DC/OS via your browser:
* DC/OS: http://localhost/#/dashboard/
* Marathon: http://localhost/marathon
* Mesos: http://localhost/mesos/#/

## Deploy Container into Marathon

To run the GraphQL-Server and both databases, you have to deploy a task group into Marathon using the API.

First, make sure no app is running:

```sh
> curl localhost/marathon/v2/apps
{"apps":[]}
```

Now, you can POST the `marathon-compose.json` to the `groups` API.

```sh
> curl -X POST -H "Content-Type: application/json" http://localhost:80/marathon/v2/groups -d@marathon-compose.json
{"version":"2016-06-04T20:53:54.516Z","deploymentId":"07bb968e-b454-4b95-8d60-80fcad57f9d9"}
```

To track the progress of the deployment you can use the [Marathon Web UI](http://localhost/marathon).
