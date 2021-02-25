# 退群杯后端 数据库结构文档

## 用户 user

- _id `ObjectId` 用户 ID
- username `string` 用户名
- password `string` 密码（加密）
- email `string` 邮箱地址
- last_question `ObjectId` 最后做的题目
- last_scene `ObjectId` 最后看的剧情
- start_time `int` 最后做题开始时间
- unlocked_scene `Array<ObjectId>` 已解锁剧情列表
- finished_question `Array<ObjectId>` 已做答问题列表
- complete_count `int` 完成周目数

## 题目 question

- _id `ObjectId` 题目 ID
- title `string` 标题
- desc `string` 描述
- statement `string` 题面
- sub_question `Array<#1>` 题目包含的所有子题信息
- author `string` 出题人
- audio `string` 音频 URL（仅听力）
- time_limit `int` 时限，单位为秒
- next_scene_text `string` 下一剧情列表展示的文本
- next_scene `Array<#2>` 下一剧情列表

subquestion 内的对象：

- type `int` 类型：1 上传 PDF，2 选择
- desc `string` 子题描述
- option `Array<string>` 选项文本（仅选择题）
- true_option `Array<int>` 正确选项下标（仅选择题）
- full_point `float` 满分值
- part_point `float` 多选未选全分值（仅选择题）

next_scene 内的对象：

- scene `ObjectId` 剧情 ID
- option `string` 选项文本

## 剧情 scene

- _id `ObjectId` 剧情 ID
- from_question `ObjectId` 上一题目
- next_question `ObjectId` 下一题目
- title `string` 剧情标题
- text `string` 剧情文本
- bgm `string` BGM 音频文件 URL

## 提交 submission

- _id `ObjectId` 提交 ID
- user `ObjectId` 答题人 ID
- question `ObjectId` 题目 ID
- time `int` 提交时间
- file `Array<string>` 提交文件 ID（仅上传 PDF）
- option `Array<int>` 选项索引（仅选择题）
- point `float` 提交获得分数，-1 为未评分
- answer_time `int` 答题用时，单位为秒
- is_time_out `bool` 是否为超时提交

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