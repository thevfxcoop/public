# public

Repository for public information pertaining to vfx.coop

  * The `docs` directory contains the static website, which is served with github pages
  * The `lambda` directory contains AWS lamba functions which are used for any dynamic
    elements of the website.

## Deploying the Lambda function

The Lambda function requires the following in order to deploy:

  * Make, Go compiler, Zip and AWS command-line tools
  * The environment variable AWS_ACCOUNT_ID should contain your AWS account
  * A role "LambdaExecRole" is created for executing Lambda functions

In the root directory, do the following:

```bash
[bash] cd lambda
[bash] make
Building contactus for Linux
Adding: build/contactus (deflated 48%)
{
    "FunctionName": "contactus",
    "FunctionArn": "...",
    ...
}
```

The JSON return provides details of the uploaded Lambda function. The triggering of the
function is not covered here, there is some Terraform which creates an API gateway 
detailed elsewhere.


