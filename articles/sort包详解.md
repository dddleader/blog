---
title: "Go语言 sort 包详解：从基础到进阶"
date: "2024-01-20"
tags: 
  - "Go"
  - "标准库"
  - "算法"
  - "排序"
summary: "详细介绍了 Go 语言 sort 包的使用方法，包括基本排序功能、自定义排序实现以及二分查找等核心功能，适合想要深入了解 Go 排序机制的开发者。"
status: "published"
category: "技术教程"
cover: "/static/images/go-sort.jpg"
author: "您的名字"
---

# SORT
## 基本排序功能
```go
ints := []int{3, 1, 2}
sort.Ints(ints) // 结果: [1, 2, 3]

floats := []float64{3.1, 1.5, 2.3}
sort.Float64s(floats) // 结果: [1.5, 2.3, 3.1]

strings := []string{"banana", "apple", "cherry"}
sort.Strings(strings) // 结果: ["apple", "banana", "cherry"]

```
## 自定义排序
```go
type Person struct {
    Name string
    Age  int
}
people := []Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Charlie", 35},
}
// 按照 Age 升序排序
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})

//提供稳定排序
sort.SliceStable(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})

//自定义接口
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

sort.Sort(ByAge(people)) // 使用自定义排序规则排序

```
## 二分搜索
```go
arr := []int{1, 3, 5, 7, 9}
idx := sort.SearchInts(arr, 7) // 返回索引 3
```