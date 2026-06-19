# Spec 库索引

本目录是 specflow 项目的 spec 库根。spec 是项目长期沉淀的行为规约、编码规范、架构决策等知识，区别于任务级的 prd.md / implement.md（一次性）。

`specflow-update-spec` skill 会把任务中值得沉淀的经验固化到这里。

## Pre-Development Checklist

动手写代码前，AI 必须确认以下事项（由 `specflow-before-dev` skill 强制）：

- [ ] 已读取当前任务的 prd.md 与 implement.md
- [ ] 已读取 implement.jsonl manifest 中声明的所有 spec 文件
- [ ] 已确认当前任务 status 为 in_progress（非 planning / completed）
- [ ] 已确认要改的文件在 relatedFiles 或 jsonl manifest 范围内
- [ ] 不执行 git commit / push / merge（交给 finish-work）
- [ ] 遇到 prd 不明确的点，先停止并报告，不臆断

## 目录结构（按需创建）

```
spec/
├── index.md            # 本文件，全局索引
├── backend/            # 后端相关 spec
│   └── index.md
├── frontend/           # 前端相关 spec
│   └── index.md
└── architecture/       # 架构决策记录（ADR）
    └── index.md
```

## 索引占位

> 下方各 section 在 spec 增加后由 AI 维护。新增 spec 文件时，在此追加一行 `- [相对路径](相对路径) — 一句话摘要`。

### Backend

（暂无）

### Frontend

（暂无）

### Architecture

（暂无）
