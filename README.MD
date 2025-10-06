# Distributed System – Product API

This project deploys a simple **Go-based REST API** to AWS ECS using **Terraform** and **Docker**.  
The service exposes endpoints to retrieve and manage product data.

These instructions explain exactly how to deploy and test the system on **any machine**.



## Instructions

### Prepare Credentials

Retrieve you temporary credentials from Learner's Lab.
Enter your configuration when prompted:
```
aws configure
```

And set your session token:
```
aws configure set aws_session_token <YOUR-TEMP-SESSION-TOKEBN>
```

### Apply Infrastructure
```
cd terraform
terraform init
terraform apply -auto-approve
```

Once complete, Terraform will output:
ecs_cluster_name
ecs_service_name

### Step 3: Build and Push Docker Image
```
cd src
docker buildx build --platform linux/amd64 \
  -t <YOUR_ECR_REPO_URI>:latest \
  --push .
```

Then redeploy by re-running:
```
cd ../terraform
terraform apply -auto-approve
```

### Send Requests
Make sure you are in terraform folder

Get public ip address
```
aws ec2 describe-network-interfaces \
--network-interface-ids $(
    aws ecs describe-tasks \
    --cluster $(terraform output -raw ecs_cluster_name) \
    --tasks $(
        aws ecs list-tasks \
        --cluster $(terraform output -raw ecs_cluster_name) \
        --service-name $(terraform output -raw ecs_service_name) \
        --query 'taskArns[0]' --output text
    ) \
    --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" \
    --output text
) \
--query 'NetworkInterfaces[0].Association.PublicIp' \
--output text
```

### Test API Endpoints
```
curl http://<PUBLIC-IP-ADDRESS>:8080/products
```


## Clean Up
```
terraform destroy -auto-approve
```