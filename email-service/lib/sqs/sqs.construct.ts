import { Duration, StackProps } from "aws-cdk-lib"
import { Queue, QueueEncryption } from "aws-cdk-lib/aws-sqs"
import { Construct } from "constructs"

export interface GoEmailServiceSqsProps extends StackProps {
  queueName: string
  visibilityTimeout: number
  retentionPeriod: number
  deadLetterQueue?: {
    maxReceiveCount: number
    queue: Queue
  }
}

export class GoEmailServiceSqs extends Construct {
  public readonly queue: Queue

  constructor(scope: Construct, id: string, props: GoEmailServiceSqsProps) {
    super(scope, id)

    this.queue = new Queue(this, `${id}Construct`, {
      queueName: props.queueName,
      visibilityTimeout: Duration.seconds(props.visibilityTimeout),
      retentionPeriod: Duration.days(props.retentionPeriod),
      encryption: QueueEncryption.SQS_MANAGED,
      deadLetterQueue: props.deadLetterQueue,
    })
  }
}
