STACK = gotestaws
LFUNCTION = "HelloWorldFunction"
RFUNCTION = "gotestaws-HelloWorldFunction-xxre27wBUdFp"

.PHONY: build deploy validate destroy

invokelocal: build
	sam local invoke $(LFUNCTION) -e event.json

invokeremote:
	aws lambda invoke --function-name $(RFUNCTION) \
		--invocation-type Event \
		--cli-binary-format raw-in-base64-out --payload '{"hello": "John Smith"}' /dev/stdout

tail:
	sam logs --tail --stack-name $(STACK)

deploy: build
	SAM_CLI_TLEMETRY=0 sam deploy --resolve-s3 --stack-name $(STACK) \
		 --no-confirm-changeset --no-fail-on-empty-changeset --capabilities CAPABILITY_IAM --disable-rollback

build:
	sam build --use-container

validate:
	aws cloudformation validate-template --template-body file://template.yml

destroy:
	aws cloudformation delete-stack --stack-name $(STACK)
