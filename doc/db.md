# 退群杯后端 数据库结构文档

## 用户 user

- _id `ObjectId` 用户 ID
- username `string` 用户名
- password `string` 密码（加密）
- email `string` 邮箱地址
- is_email_verified `bool` 邮箱是否验证
- last_question `ObjectId` 最后做的题目，用于实现继续答题
- start_time `int` 最后做题开始时间
- unlocked_scene `Array<ObjectId>` 已解锁剧情列表

## 题目 question

- _id `ObjectId` 题目 ID
- title `string` 标题
- desc `string` 描述
- type `int` 类型：1 选择，2 上传 PDF
- option `Array<string>` 选项文本（仅选择题）
- true_option `int` 正确选项索引
- full_point `float` 总分
- author `string` 出题人
- audio `string` 音频 URL（仅听力）
- time_limit `int` 时限，单位为秒
- next_scene `Array<#1>` 下一剧情列表

next_scene 内的对象：

- scene `ObjectId` 剧情 ID
- option `string` 选项文本

## 剧情 scene

- _id `ObjectId` 剧情 ID
- from_question `ObjectId` 上一题目
- next_question `ObjectId` 下一题目
- title `string` 剧情标题
- text `string` 剧情文本

## 提交 submission

- _id `ObjectId` 提交 ID
- author_id `ObjectId` 答题人 ID
- question_id `ObjectId` 题目 ID
- time `int` 提交时间
- file `string` 提交文件 ID（仅上传 PDF）
- option `int` 选项索引（仅选择题）
- point `float` 提交获得分数，-1 为未评分

## 学科 subject

- _id `ObjectId` 学科 ID
- abbr `string` 学科简写
- name `string` 学科名称
- start_scene `ObjectId` 学科初始剧情

## 文件 file

- _id `ObjectId` 文件 ID
- filename `string` 实际文件名
- user `ObjectId` 上传用户 ID
- time `int` 上传时间