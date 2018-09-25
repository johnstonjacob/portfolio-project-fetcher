GOOS=linux go build -o main && zip main.zip main && aws lambda update-function-code --region us-east-2 --function-name portfolio-projects --zip-file fileb://./main.zip
