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

export AWS_ACCESS_KEY_ID="ASIA27Z34IUCG2PNF2VZ"
export AWS_SECRET_ACCESS_KEY="L0Qfr9cTTkWm4fAcPvJGkvqbUyTHCNePLigY9S1Z"
export AWS_SESSION_TOKEN="IQoJb3JpZ2luX2VjEEYaCXVzLWVhc3QtMSJHMEUCIC1we+o0k5Ha6oVWcNsfoFUJLVLnxfjya/qvupYqfWonAiEAkT9LnEUBi2gX+S0P1yW+fqFKlSq4TmTrt2Z7+XwuIDYqlgMIPxAAGgw3NTU1MDI5NTc4MjgiDHe2r1lYAcLR9WvrdCrzAoGX7wlm+6UiPOLUBMKsGiGWm91x7D9tdcRAeaZroztGIqChoUf18dPWpLtHmi13mK9Ko8UtB9+x6nCZlZcmmGFjVG5mj0zgXENFn1P98lVoJ9q8HABeL/Xk6aRSgFNWu1rLYusqDOgaQRvYGwmsXj8t8HYtk5t8742OeFKl7o8sqAynhF0Xp6drqzI2i4r5IKyacPCkqV8QD2cd09UkfYy7oaI0KCnYJfIKhTO/rOVmuCrGVaHmnNBaFgW2liJiN2xrVI02uGzQfNTvICcQ1Lu8NO5CSVrJg11bXBGN3Mx60HizC56jWNFvAzKQj2txDy5MXGiC6rot6ArF7dOFIKHzlYaqL1ZU1jyqdtJx/QRAmjEZo1VLH7RFVcKsOvbZRryVr71Br2+0bJjKOovfbDumM76TC9t5U5LFspIyFCnPPxDeZtEwsk/TiSN8sa0tZidkd+5XC6L3DLBp8MfmO/cKwpmwIBlJ1O8tu++fdLON/B+7MNi3xIcGOqYBWUsRa4WJ4lPxG7dCuPiBD+qIKrBt43NkyfW+0c9jOr7CtNPkd4vU5NBaLKrOwKY1mYj/Qp0A63ZhQT2mH0AeV8yIgMtszFt+FwT6z9m6FlTV5h/I0uzdD5XJ/+ufl/ZPvYK9VEp9v2GYFqQfiqRSMgYW+1boQ82giIndYUyvoiusO+5wh1liBCW4V6VZJ9ScFW45IuNPRp+RDykMGwE4C+qkfa2N7w=="


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



