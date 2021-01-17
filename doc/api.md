# 退群杯后端 API 文档

## 数据交换格式

### 身份验证

在 Header 中加入 `Authorization` 字段进行验证，将获取的 JWT 令牌作为 Bearer Token 加入该字段的值，例如：

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U
```

以下接口中，标题带有 `*` 标记的为需要身份验证的接口。

### 响应格式

响应使用 JSON 格式，例如：

```json
{
    "success": true,
    "error": null,
    "data": {
        // ...
    }
}
```

### URL 前缀

文档中所有接口 URL 都包含前缀 `/api/v1`。

## 用户相关

### 获取登录公钥 GET /user/public_key?email={邮箱地址}

为了保证密码安全，登录时用 RSA 加密密码传输，获取一次公钥有效期 15 分钟。

响应：

```json
{
    "public_key": "base64 编码的 1024 位 RSA 公钥"
}
```

### 登录 GET /user/token?email={邮箱地址}&password={加密后的密码}

获取 JWT 令牌，调用前使用获取公钥接口。

响应：

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
}
```

### 注册 POST /user

新建用户，调用前使用获取公钥接口。

请求：

```json
{
    "email": "123@qq.com",
    "username": "张三",
    "password": "加密后的密码"
}
```

响应：空

### 请求修改密码 GET /user/reset_mail?email={邮箱地址}

向邮箱发送一封带有修改密码链接的邮件，若用户不存在则不发送邮件且不提供响应。

响应：空

### (\*)修改密码 PUT /user/password?reset_id={修改密码ID}

修改密码，调用前使用获取公钥接口。用户登录后可不加修改密码 ID。

请求：

```json
{
    "password": "加密后的密码"
}
```

响应：空

### \*获取解锁剧情 GET /user/:id/unlocked_scene

获取用户自己所有已解锁剧情。

响应：

```json
{
    "scene": [
        {
            "_id": "剧情ID",
            "title": "剧情标题",
            "text": "剧情内容",
            "next_question": "下一问题ID"
        }
    ]
}
```

### \*获取提交 GET /user/:id/submission

获取用户自己所有的提交。

响应：

```json
{
    "submission": [
        {
            "_id": "提交ID",
            "time": 1610880000, // 提交时间
            "file": "提交文件ID",
            "option": 0, // 提交选项索引
            "point": 5.0 // 该题得分，-1 为未评分
        }
    ]
}
```

## 题目相关

### \*获取学科列表 GET /subject

获取所有学科列表。

响应：

```json
{
    "subject": [
        {
            "abbr": "简称",
            "name": "名称",
            "start_scene": "初始剧情ID"
        }
    ]
}
```

### \*获取剧情信息 GET /scene/:id

获取指定剧情信息，仅可获取已解锁或可解锁剧情。

响应：

```json
{
    "text": "剧情文本",
    "next_question": "下一题目 ID"
}
```

### \*获取题目信息 GET /question/:id

获取指定题目信息，仅可获取已解锁或可解锁题目。

响应：

```json
{
    "title": "标题",
    "desc": "描述",
    "type": 1, // 类型：1 选择，2 上传 PDF
    "option": ["选项1", "选项2"],
    "true_option": 0, // 正确选项索引
    "full_point": 5.0, // 总分
    "author": "出题人",
    "audio": "音频URL",
    "time_limit": 300, // 时间限制，单位为秒
    "next_scene": [
        {
            "scene": "剧情ID",
            "option": "选项文本"
        }
    ]
}
```

### \*提交题目解答 POST /qustion/:id/submission

提交题目解答，仅可提交当前题目。

请求：

```json
{
    "option": 0, // 作答选项索引（仅选择题）
    "file": "文件 ID" // 提交 PDF 文件 ID（仅上传 PDF）
}
```

响应：

```json
{
    "_id": "提交 ID"
}
```

## 其他

### \*上传文件 POST /file

上传 PDF 文件用，不支持其他格式。

请求：multipart/form-data 格式的文件，key 设置为 `file`

响应：

```json
{
    "_id": "上传的文件 ID"
}
```

### \*获取文件 GET /file/:id

获取已上传的 PDF 文件，仅可获取用户自己上传的文件。

响应：文件
