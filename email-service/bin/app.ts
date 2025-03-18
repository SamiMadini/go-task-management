#!/usr/bin/env node
import "source-map-support/register"
import * as cdk from "aws-cdk-lib"
import { EmailServiceStack } from "../lib/email-service-stack"
import { EmailServiceEnv } from "../lib/helper/env.helper"

const app = new cdk.App()

const emailServiceEnv = new EmailServiceEnv({
  node: app.node,
})

const tags = {
  environment: emailServiceEnv.environment,
  project: "go-task-manager-email-service",
}

new EmailServiceStack(app, `${emailServiceEnv.environmentUpperCase}${EmailServiceStack.name}`, {
  ...emailServiceEnv,
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION || "us-east-1",
  },
  tags,
})
