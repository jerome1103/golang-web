# GO-WEB

# 專案使用
* IDE: VSCode
* go 1.20
* MongoDB
* Redis
* Docker

## 專案API

### 呼叫查詢商品API

* [GET]查詢所有商品

```javascript
http://localhost:8000/api/item
```

* 回傳Json內容

```javascript
{
    "message": [
        {
            "Iid": 1,
            "Name": "掀蓋長夾 雙C琺瑯扣 粉色",
            "Price": 23,800,
            "Text": "附加配件｜原廠保卡、紙盒、塵袋、商品卡"
        },
		{
            "Iid": 2,
            "Name": "Cabas Navy 帆布包 黑色",
            "Price": 15,800,
            "Text": "主體配件｜皮革背帶、包包內袋"
        }
    ],
    "status": "success"
}
```

* [GET]查詢單筆商品

```javascript
http://localhost:8000/api/item/{id}
```

* 回傳Json內容

```javascript
{
    "message": [
        {
            "Iid": 1,
            "Name": "掀蓋長夾 雙C琺瑯扣 粉色",
            "Price": 23,800,
            "Text": "附加配件｜原廠保卡、紙盒、塵袋、商品卡"
        }
    ],
    "status": "success"
}
```

### 呼叫新增商品API

* [POST]新增商品

```javascript
http://localhost:8000/api/item/
```

* 代入的參數

```javascript
name:皮帶
price:2333
text:備註
```

* 回傳Json內容

```javascript
{
    "message": "6457c510064ed9c0757e3a47",
    "status": "success"
}

### 呼叫修改商品API

* [PUT]更新指定商品

```javascript
http://localhost:8000/api/item/{id}
```

* 回傳Json內容

```javascript
{
    "message": {
        "MatchedCount": 1,
        "ModifiedCount": 1,
        "UpsertedCount": 0,
        "UpsertedID": null
    },
    "status": "success"
}
```

### 呼叫刪除商品API

* [DELETE]刪除指定商品

```javascript
http://localhost:8000/api/item/{id}
```
* 回傳Json內容
```javascript
{
    "message": {
        "DeletedCount": 1
    },
    "status": "success"
}
```

### 呼叫商品排行榜

* [GET]查詢商品排行榜

```javascript
http://localhost:8000/api/rank
```
* 回傳Json內容
```javascript
{
    "message": [
        {
            "Score": 40,
            "Member": "LV"
        },
        {
            "Score": 50,
            "Member": "Dior"
        },
        {
            "Score": 60,
            "Member": "Prada"
        }
    ],
    "status": "success"
}
```

### 更新商品排行榜

* [POST]更新商品排行榜

```javascript
http://localhost:8000/api/rank
```
* 代入的參數

```javascript
name:Prada
score:99
```

* 回傳Json內容
```javascript
{
    "message": [
        {
            "Score": 99,
            "Member": "Prada"
        }
    ],
    "status": "success"
}
```
