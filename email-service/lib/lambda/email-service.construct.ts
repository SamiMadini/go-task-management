import { LambdaProps } from "../helper/env.helper"
import { CfnOutput, Duration, DockerImage } from "aws-cdk-lib"
import * as lambda from "aws-cdk-lib/aws-lambda"
import { Construct } from "constructs"
import { EmailServiceStackProps } from "../email-service-stack"
import { IFunction } from "aws-cdk-lib/aws-lambda"
import { Queue } from "aws-cdk-lib/aws-sqs"
import { SqsEventSource } from "aws-cdk-lib/aws-lambda-event-sources"

export interface PostgresConfigInterface {
  readonly host: string
  readonly port: number
  readonly user: string
  readonly password: string
  readonly db: string
}

export interface EmailServiceProps extends EmailServiceStackProps {
  readonly postgresConfig: PostgresConfigInterface
  readonly emailServiceLambda: LambdaProps
  readonly sqs: Queue
}

export class EmailService extends Construct {
  readonly function: IFunction

  constructor(scope: Construct, id: string, props: EmailServiceProps) {
    super(scope, id)

    this.function = new lambda.Function(this, `${id}Name`, {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: "bootstrap",
      architecture: lambda.Architecture.ARM_64,
      code: lambda.Code.fromAsset("./src", {
        bundling: {
          image: DockerImage.fromRegistry("golang:1.24"),
          command: [
            "bash",
            "-c",
            "GOCACHE=/tmp go mod tidy && " + "GOCACHE=/tmp GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o /asset-output/bootstrap",
          ],
        },
      }),
      timeout: Duration.seconds(props.emailServiceLambda.timeout),
      memorySize: props.emailServiceLambda.memorySize,
      retryAttempts: props.emailServiceLambda.retryAttempts,
      logRetention: props.emailServiceLambda.logRetention,
      environment: {
        POSTGRES_HOST: props.postgresConfig.host,
        POSTGRES_PORT: props.postgresConfig.port.toString(),
        POSTGRES_USER: props.postgresConfig.user,
        POSTGRES_PASSWORD: props.postgresConfig.password,
        POSTGRES_DB: props.postgresConfig.db,
      },
    })

    const sqsEventSource = new SqsEventSource(props.sqs, {
      batchSize: 1,
      maxBatchingWindow: Duration.seconds(10),
      enabled: true,
    })

    this.function.addEventSource(sqsEventSource)

    new CfnOutput(this, "EmailServiceLambdaArn", {
      value: this.function.functionArn,
      description: `Arn for Email Service lambda for ${props.environment} env.`,
    })
  }
}
