list-stacks:
	npx cdk list

synth:
	npm run build && npx cdk synth

deploy-dev-email-service:
	npm run build && npx cdk deploy DevEmailServiceStack -c EnvName=dev

deploy-prod-email-service:
	npm run build && npx cdk deploy ProdEmailServiceStack -c EnvName=prod
