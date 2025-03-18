import { Construct } from "constructs"
import { Stack, StackProps } from "aws-cdk-lib"
import { LambdaProps } from "./helper/env.helper"
import { EmailService } from "./lambda/email-service.construct"

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

    new EmailService(this, `${props.stackName}${EmailService.name}Lambda`, {
      environment: props.environment,
      environmentUpperCase: props.environmentUpperCase,
      emailServiceLambda: props.emailServiceLambda,
      postgresConfig: postgresConfig,
    })
  }
}
