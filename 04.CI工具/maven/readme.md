#  Maven


1. maven 上传 

```azure

mvn deploy:deploy-file -Dmaven.test.skip=true -Dfile=./air-crypto-1.5.83-SNAPSHOT.jar -DrepositoryId=handpay-snapshots -DgroupId=cn.com.handpay -DartifactId=air-crypto -Dversion=1.5.83-SNAPSHOT -Dpackaging=jar -Durl=http://mvnrep.handpay.com.cn:8081/nexus/content/repositories/handpay-snapshots

```

 