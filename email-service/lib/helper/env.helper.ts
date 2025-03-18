"use strict"

import { Node } from "constructs"
import { RetentionDays } from "aws-cdk-lib/aws-logs"

export enum LogLevel {
  ERROR = "error",
  WARN = "warn",
  INFO = "info",
  HTTP = "http",
  VERBOSE = "verbose",
  DEBUG = "debug",
  SILLY = "silly",
}

export enum Stages {
  DEV = "dev",
  PROD = "prod",
}

export interface EnvProps {
  node: Node
  prodName?: string
}

export enum SsmParamPathPrefix {
  DEV = "dev",
  PROD = "prod",
}

export interface LambdaProps {
  memorySize: number
  timeout: number
  retryAttempts: number
  logRetention: RetentionDays
  logLevel: LogLevel
}

export class Env {
  readonly envNameParameterName = "EnvName"
  readonly defaultEnvName = "dev"
  readonly environment: string = Stages.DEV
  readonly environmentUpperCase: string
  readonly prodName?: string
  readonly isDev: boolean
  readonly isProd: boolean
  readonly ssmParamPathPrefix: string
  readonly configurationSetName: string

  constructor(props: EnvProps) {
    this.environment = this.resolveEnvName(props.node)
    this.isDev = this.resolveIsDev()
    this.isProd = this.resolveIsProd()
    this.environmentUpperCase = this.resolveEnvNameUpperCase()
    this.ssmParamPathPrefix = this.resolvePrefixForSsmParamPath()
  }

  private resolveEnvName(node: Node): string {
    return node.tryGetContext(this.envNameParameterName) || this.defaultEnvName
  }

  private resolveIsDev(): boolean {
    return this.environment !== this.prodName
  }

  private resolveIsProd(): boolean {
    return this.environment === Stages.PROD
  }

  private resolveEnvNameUpperCase(): string {
    return this.environment[0].toUpperCase() + this.environment.slice(1)
  }

  private resolvePrefixForSsmParamPath(): string {
    switch (this.environment) {
      case Stages.DEV:
        return SsmParamPathPrefix.DEV
      case Stages.PROD:
        return SsmParamPathPrefix.PROD
      default:
        throw new Error("Unknown env name")
    }
  }
}

export interface EmailServiceEnvInterface {
  stackName: string
  stackPrefixDescription: string
  emailServiceLambda: LambdaProps
}

export class EmailServiceEnv extends Env implements EmailServiceEnvInterface {
  stackName: string
  stackPrefixDescription: string
  emailServiceLambda: LambdaProps

  constructor(props: EnvProps) {
    super(props)
    this.stackName = `${this.environmentUpperCase}EmailService`
    this.stackPrefixDescription = `Email Service of ${this.environment}`
    this.initVarsForStack()
    Object.freeze(this)
  }

  private initVarsForStack(): void {
    switch (this.environment) {
      case Stages.DEV:
        this.setVarValuesForDev()
        break
      case Stages.PROD:
        this.setVarValuesForProd()
        break
      default:
        throw new Error("Unknown env name")
    }
  }

  private setVarValuesForDev(): void {
    this.emailServiceLambda = {
      memorySize: 512,
      timeout: 1,
      retryAttempts: 1,
      logRetention: RetentionDays.ONE_MONTH,
      logLevel: LogLevel.INFO,
    }
  }

  private setVarValuesForProd(): void {
    this.emailServiceLambda = {
      memorySize: 1024,
      timeout: 5,
      retryAttempts: 2,
      logRetention: RetentionDays.ONE_YEAR,
      logLevel: LogLevel.INFO,
    }
  }
}
