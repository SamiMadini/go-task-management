# Email Service Lambda CDK Deployment

This directory contains the AWS CDK code to deploy the Email Service as a Lambda function.

## Prerequisites

- Node.js 14.x or later
- AWS CLI configured with appropriate credentials
- AWS CDK CLI installed (`npm install -g aws-cdk`)

## Setup

1. Install dependencies:

```bash
npm install
```

2. Build the TypeScript code:

```bash
npm run build
```

3. Bootstrap CDK (if you haven't done this before in your AWS account/region):

```bash
cdk bootstrap
```

## Deployment

To deploy the Email Service Lambda:

```bash
cdk deploy
```

This will:

- Build the Go Lambda function
- Package it for Lambda deployment
- Create the Lambda function in AWS
- Set up the necessary IAM permissions for SES

## Useful CDK Commands

- `npm run build` compile typescript to js
- `npm run watch` watch for changes and compile
- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `cdk synth` emits the synthesized CloudFormation template
