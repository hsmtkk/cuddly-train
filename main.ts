// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0
import { Construct } from "constructs";
import { App, TerraformStack, CloudBackend, NamedCloudWorkspace } from "cdktf";
import * as google from '@cdktf/provider-google';

const project = 'cuddly-train';
const region = 'asia-northeast1';
const repository = 'cuddly-train';

class MyStack extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new google.provider.GoogleProvider(this, 'google', {
      project,
      region,
    });

    new google.artifactRegistryRepository.ArtifactRegistryRepository(this, 'artifact_registry', {
      format: 'docker',
      location: region,
      repositoryId: 'registry',
    });

    new google.cloudbuildTrigger.CloudbuildTrigger(this, 'build_trigger', {
      filename: 'cloudbuild.yaml',
      github: {
        owner: 'hsmtkk',
        name: repository,
        push: {
          branch: 'main',
        },
      },
    });

    /*
    new google.containerCluster.ContainerCluster(this, 'my_cluster', {
      name: 'my-cluster',
      enableAutopilot: true,
      location: region,
      // To avoid this error
      // Error 400: Max pods constraint on node pools for Autopilot clusters should be 32., badRequest
      ipAllocationPolicy: {},
    });
    */
  }
}

const app = new App();
const stack = new MyStack(app, "cuddly-train");
new CloudBackend(stack, {
  hostname: "app.terraform.io",
  organization: "hsmtkkdefault",
  workspaces: new NamedCloudWorkspace("cuddly-train")
});
app.synth();
