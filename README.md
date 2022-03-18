A simple bank completed with golang+Postgres+Docker. it's a learningProject
from <a href="https://github.com/techschool/simplebank">原项目</a>

mockgen: 数据库模拟程序。可以用于测试操作数据库代码

```mockgen -package mockdb -destination 输出文件路径 数据库接口的包路径 数据库操作接口
   mockgen -package mockdb -destination db/mock/store.go github.com/0RAJA/Bank/db/sqlc Store
```

JWT 的弃用

1. 可选择的加密算法参差不齐，有些已经失效
2. 黑客可以将非对称加密方式更改为对称加密，然后使用公钥加密JWT然后欺骗服务器

PASETO 的使用

1. 算法强大而稳定，使用者只需要考虑其版本即可。
2. PASETO 也存在两套加密算法。本地使用对称加密。外部使用非对称。
3. PASETO 对整个token进行加密和验证,无法伪造算法头

PASETO 分为四个部分 本地:

1. 版本号 v2
2. 使用场景
3. 有效载荷
    1. 数据信息和到期时间
    2. nonce 用于加密和消息认证过程中
    3. 消息认证标签 用于验证加密消息和与其关联的未加密消息(版本号，使用场景和页脚)
4. 公共信息(仅用base64编码)

外部:

1. 版本号
2. 使用场景
3. 有效载荷
    1. 数据信息和到期时间(base64编码)
    2. 使用私钥进行加密的数字证书用于校验真实性
