
# Paths to packages
GO=$(shell which go)
ZIP=$(shell which zip)
AWS=$(shell which aws)

# Execution role. Please see Terraform repo for where this role
# gets created
AWS_ROLE="arn:aws:iam::${AWS_ACCOUNT_ID}:role/LambdaExecRole"

# Paths to locations
BUILD_DIR=build
FUNCTION_EXISTS=$(shell aws lambda get-function --function-name contactus 1>/dev/null 2>/dev/null; echo $$?)

# Build production by default
all: contactus

contactus: dependencies clean builddir
	@echo "Building contactus for Linux"
	@GOARCH=amd64 CGO_ENABLED=0 GOOS=linux ${GO} build -o ${BUILD_DIR}/contactus ./contactus
	@${ZIP} ${BUILD_DIR}/contactus.zip build/contactus
ifeq ($(FUNCTION_EXISTS),0)
	@aws lambda update-function-code --function-name contactus --zip-file fileb://./${BUILD_DIR}/contactus.zip	
else
	@aws lambda create-function --function-name contactus --runtime go1.x --zip-file fileb://./${BUILD_DIR}/contactus.zip --handler build/contactus --role ${AWS_ROLE}
endif

dependencies:
ifeq (,${GO})
        $(error "Missing go binary")
endif
ifeq (,${ZIP})
        $(error "Missing zip binary")
endif
ifeq (,${AWS})
        $(error "Missing aws binary")
endif
ifeq (,${AWS_ACCOUNT_ID})
        $(error "Missing AWS_ACCOUNT_ID environment variable")
endif

clean:
	@rm -fr ${BUILD_DIR}
	@${GO} mod tidy

builddir:
	@install -d ${BUILD_DIR}
