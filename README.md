# go-config
基于golang的配置管理
## 功能
* 客户端，服务端采用websocket链接，避免poll
* 支持使用git作为配置版本管理
* 支持git branch，tag, commitId作为获取配置条件
* 支持仓库路由
* 支持动态仓库配置，参见"仓库模型"
* 支持按需推送配置更新
* *计划支持 提供webhook api
* *计划支持 可在浏览器页面直接修改仓库文件功能


## 支持
为了适应不同的系统，项目采用gox进行交叉编译，并使用upx对可执行文件进行压缩，这些可执行文件都存放于项目bin文件夹下，根据需要选择合适的可执行文件，目前支持的系统：
* MacOS 32/64
* Linux 32/64
* Windows 32/64

## 使用
### client
```
./client_{os}_{arch} configClient.yml
```
configClient.yml
```yaml
#server端地址
server: localhost:5337
#如果Server不可用或其它原因导致连接断开，多久会进行重试，单位：秒
tick: 5
#要获取配置的应用信息
app:
  - name: app01   #应用名称
    label: test   #对应git的branch，tag，commitId
    profile: dev  #profile支持，设置此参数会只从服务端获取application-{profile}之类的文件
                  #多个profile使用逗号","分隔
    homePath:     #该应用的配置需要存放的本地目录
      - /Users/liolay/configRepo
#  - name: app02   
#    label: master      
#    homePath:          
#      - /Users/liolay/configRepo
#


```

### server
```
./server_{os}_{arch} configServer.yml
```
configServer.yml
```yaml
#本地存储配置仓库的目录
homePath: /Users/liolay/haha
#启动端口
port: 5337
#默认仓库地址
defaultRepo: git@github.com:liolay/go-config-repo-all.git
#有访问仓库权限的sshKey
sshKey: |-
        -----BEGIN RSA PRIVATE KEY-----
        MIIEpAIBAAKCAQEAyB+gq1xjj4HCe4hPLfA3a9y4pXVhrJfQQheHR1Mj9Y2eBJTm
        4jKtHwYWOxVbUahmim+vUuVia5nH5LYA9c/GZxUKrhHwND0hF7q7kKRMboT2w/2J
        JRPkVC3T2I02ptADTyFCZLNUViviFm6JoVEnEZAwlBDakvBQBqYgQstrSE2mH8wN
        niM7T7U6gPChpgPVpuTOBNbsgfag7bgAUVl56i3IFvht+N0LXKuNU0TTuMUq2+lR
        i+iha86OtkSYZqA34SXpoL2cyQJUikMkPsjXsxfpq1f35AbPC4FUfNnnN/NRBD/2
        jOHHQSHcMOq0Y/SapCXgxa4urlLkCOra1lC5hQIDAQABAoIBAQCeK5E7nzv5cp+a
        L3QVZOUI1V0DOTFHzn2Fnz8GeonTTGj2ShHp+g+mk5MCg7C3a5gQFpHFvRL65IJ/
        G/LKVbwEQTc9uWPWhfIf5TDV82WNfH3lDgBVU9GFTus/Hu1xDrtu0WS+XpZrvSdm
        f1s8Kv3r/cDHZkK7HEDD4I1i/Y//hinEOFbZglgmgGCwF3aZUdfj0wq1t1CnVJsj
        P3c+rbxtvBLPblizHunqMiNHM1DTYqK5sW92yyMvy6ZbN3fSntcEdUNLWl61uWh1
        E0eeg+GSuQoY/STFNqeTCqPju+gW8YR+Qix1R2n/AMo/D3HXcKxZCxWSD5zY8bH0
        C/5oZerBAoGBAOdLjhQb2L1xU0V7vHWIHsWeOOTaDu1IhratY+detY/D0VCbnTXx
        CZQmCgwT88/iPLVbEg3tsBBlVRYdTSOwtGk4Dvis83Izs1kU4RV8AI3+pZ74qA6O
        8V3xrvaIrnJWHX0M078EAOFnj9zvkMFQnovnSyli+GSzPAreHVAmOMANAoGBAN1/
        uojTaFg+cCfD6mQACo8DMN2CbWxJjXo4cpGv2HSYcjrVK+knjCdFhEnquLa0lQ+u
        TxvNTFDhPDnEXFUegZwD+ge2oDI7wP/qAJUpgU53W12OvaRs5gMHKyU8AHWrMgYk
        PvEwZIOlEc7ceyWzKCyt7nZwPHiteRsv+QgSS4lZAoGAUH27sQXL1ImWmAyqliBL
        zSv10raMEUl3ECWhKciM2L4lnq649Cew1Ky0PGXJKGQsClTqIIzCA8Kv7KU/zhbV
        gfRvSV0uz2RsmqioeAiSTNf8nSkdmwtltfLAl60TQFj1pCoNmmDzSX3308RPFOdQ
        dZGFV57IoIq7b3DCtLzIbRUCgYB8ZFAYqUlPTXllC6SlllRXrn4R2D6lcsUuX2cQ
        JEYWbMqx+aeYX+pY37SEYnpruQyBau3oeioives5sen8r44wVRdkn45lx6MC1aKQ
        ImgI7gT0jMY6AiJGjw8O8Rx8+LC2PELQ5tF8EQboOnA6YtvsA54JC80aJKn/t7hO
        bR/YuQKBgQDik5Es/Y2NM0JSKQwy7myGP4g5SadZd5Pv4r89OY/8mEOhKvHlaxki
        M2o2is6RiFTKn/aFWXl0/QKXyupq35jsesq78y29zoUuph7J5PYcLgFF2t49q3B4
        1Pn/53pMethW8zhphp3fM77vBkUu/NcusttbQjtdLEI7dkR/kqEc4fQ==
        -----END RSA PRIVATE KEY-----
#仓库        
route:
  - pattern: #匹配模式，支持"*"匹配，匹配规则是针对于 client端app节点的name和profile属性的，即：{name}/{profile}
      - "*/*"
    repo: git@github.com:liolay/${app}.git #仓库地址支持动态参数，这些参数来自于client端app节点的name和profile属性
                                           # ${app}     -> client.app.name
                                           # ${profile} -> client.app.profile
  - pattern:
      - "n*/pro*"
    repo: git@github.com:liolay/${app}-${profile}.git

```
  仓库路由和动态仓库为组织仓库模型提供了更好的灵活性，目前支持3种仓库模型，参见"仓库模型"，
在配置动态仓库时如果不遵循仓库模型，将会得到意想不到的结果. 

  作为快速理解配置和仓库模型的对应关系，你可以理解为：
  
* 模式一：仓库地址不含有变量
* 模式二：仓库地址含有${app}
* 模式三：仓库地址含有${app}，${profile}

`关于为什么不在仓库地址中加入${label}:个人认为，label应当仅作为版本管理使用，在仓库地址中加入改变量会影响到版本控制`

### 刷新配置
server提供了如下endpoint可以方便刷新特定范围的配置：
<table>
    <tr>
        <th >endpoint</th>
        <th >作用</th>
    </tr>
        <tr>
            <td>/refresh</td>
            <td>刷新所有配置</td>
        </tr>
  <tr>
            <td>/refresh/:app</td>
            <td>刷新指定应用配置</td>
        </tr>   
  <tr>
            <td>/refresh/:app/:profile</td>
            <td>刷新指定应用，指定profile配置</td>
        </tr>   
  <tr>
            <td>/refresh/:app/:profile/:label</td>
            <td>刷新指定应用，指定profile,指定label配置</td>
        </tr>                     
</table>


