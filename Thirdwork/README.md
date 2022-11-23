## 2185011259_KuangShouxu 

## 基本题目
完成了在minikube的部署

过程：
  先将实验一在docker下生成镜像，再在Depolment.yaml配置文件和server.yaml配置文件下生成POD和将网络的接口暴露，在本地打开client，检验验证正确
  仍然能正常工作
  ![image](https://user-images.githubusercontent.com/64403824/203555310-8dde585e-48d4-40a3-bd03-ba61d4baab8e.png)
  
  kubectl port-forward deployment myserver  1235:1235 在端口转发下连接外部端口，以下是在shell的验证
![image](https://user-images.githubusercontent.com/64403824/203556494-fc0d6d52-b067-424a-881e-bcb1e0f3eeec.png)
