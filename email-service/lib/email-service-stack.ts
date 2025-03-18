import { Construct } from "constructs"
import { Stack, StackProps } from "aws-cdk-lib"
import { LambdaProps } from "./helper/env.helper"
import { EmailService } from "./lambda/email-service.construct"
import { GoEmailServiceSqs } from "./sqs/sqs.construct"

export interface EmailServiceStackProps extends StackProps {
  readonly environment: string
  readonly environmentUpperCase: string
  readonly emailServiceLambda: LambdaProps
}

export class EmailServiceStack extends Stack {
  constructor(scope: Construct, id: string, props: EmailServiceStackProps) {
    super(scope, id, props)

    const postgresConfig = {
      host: "localhost",
      port: 5432,
      user: "postgres",
      password: "postgres",
      db: "email_service",
    }

    const emailServiceSqs = new GoEmailServiceSqs(this, `${props.stackName}${GoEmailServiceSqs.name}`, {
      queueName: "go-email-service-queue",
      visibilityTimeout: 300,
      retentionPeriod: 4,
    })

    new EmailService(this, `${props.stackName}${EmailService.name}Lambda`, {
      environment: props.environment,
      environmentUpperCase: props.environmentUpperCase,
      emailServiceLambda: props.emailServiceLambda,
      postgresConfig: postgresConfig,
      sqs: emailServiceSqs.queue,
    })
  }
}
