# RSA密钥测试方法

## 下载官方示例
下载fabric-samples  
cd fabric-samples  
git checkout v1.4.3  
mkdir bin  

## 编译工具和生成docker镜像
cd fabric  
make clean  
make docker  
make tools-docker  

## 替换官方工具类
复制cryptogen、configtxgen 、configtxlator 到fabric-samples/bin目录下  

## 使用rsa证书启动fabric网络
cd fabric-samples/first-network  
./byfn.sh up  
