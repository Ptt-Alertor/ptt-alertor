
if ! which terraform ; then
   echo "Install terraform package for your OS";
   exit
fi

if ! ls ./provider.tf ; then
   echo "Ask sys admin for this file for AWS access key";
   exit
fi 

echo "Start Terraform Testing..." 
terraform show

echo "Done!!"
