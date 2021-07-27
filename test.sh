#!/bin/bash

for ARGUMENT in "$@"
do
    KEY=$(echo $ARGUMENT | cut -f1 -d=)
    VALUE=$(echo $ARGUMENT | cut -f2 -d=)   

    case "$KEY" in
            env) env_name=${VALUE} ;;
            ide) ide_name=${VALUE} ;;
            identity) identity_name=${VALUE} ;;
#	    url) url=${VALUE} ;;
            *)   
    esac    
done

#make a folder for the project and copy requirements
cd $HOME/Desktop/GenerateCode
mkdir $identity_name
cp Dockerfile $identity_name
cp -r code-server $identity_name 
cd $HOME/Desktop/GenerateCode/$identity_name

export AWS_ACCESS_KEY_ID=""
export AWS_SECRET_ACCESS_KEY=""
export AWS_SESSION_TOKEN=""


#download the code from s3 
#curl $url -o $identity_name.zip
aws s3 cp s3://codercom-code-server/$identity_name.zip .
unzip $identity_name.zip
rm -rf $identity_name.zip

#build image and push to ecr
docker build -t code-server .
docker tag code-server:latest public.ecr.aws/c9l5x7x2/code-server:$ide_name
docker push public.ecr.aws/c9l5x7x2/code-server:$ide_name
cd $HOME/Desktop/GenerateCode/

#create the coder-n folder with template
rm -rf $ide_name
mkdir $ide_name
cp template/* $ide_name
sed -i "s/coder_name/$ide_name/g" $ide_name/*

#make changes in gitops and push
rm -rf gitops-paas

cd $HOME/Desktop/GenerateCode
git clone git@github.com:nslhb/gitops-paas.git -b $env_name 

cp -r $ide_name gitops-paas/coder/

rm -rf $ide_name

echo '  - '$ide_name  >>  gitops-paas/coder/kustomization.yaml


cd $HOME/Desktop/GenerateCode/gitops-paas
git add .
#git commit -am $ide_name
#git push origin $env_name



