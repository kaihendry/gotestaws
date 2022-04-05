# Testing AWS Lambda functions in Go

Learning example of how best to "business logic".

https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/

# Improvement ideas

- Use [net.IP](https://pkg.go.dev/net#IP)
- Avoid testify dependency?
- How to mock the http.Get to test error handling?
- Does it make any sense to test the handler function?

# Javascript

Since I work in Javascript, I want to apply the same ideas to Javascript in https://github.com/kaihendry/aws-jest-lambda
