# Ptt-Alertor

<img align="right" src="https://raw.githubusercontent.com/meifamily/ptt-alertor/master/logo.jpg">

[![Build Status](https://travis-ci.org/meifamily/ptt-alertor.svg?branch=master)](https://travis-ci.org/meifamily/ptt-alertor)
[![codecov](https://codecov.io/gh/meifamily/ptt-alertor/branch/master/graph/badge.svg)](https://codecov.io/gh/meifamily/ptt-alertor)
[![Go Report Card](https://goreportcard.com/badge/github.com/meifamily/ptt-alertor)](https://goreportcard.com/report/github.com/meifamily/ptt-alertor)
[![StackShare](https://img.shields.io/badge/tech-stack-0690fa.svg?style=flat)](https://stackshare.io/ptt-alertor/ptt-alertor)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## API

### Board

* GET /boards

* GET /boards/[board name]/articles

* GET /boards/[board name]/articles/[article code]

### Keyword

* GET /keyword/boards

### Author

* GET /author/boards

### PushSum

* GET /pushsum/boards

### Articles

* GET /articles

### User (Auth)

* GET /users

* GET /users/[account]

* POST /users

```json
{
    "profile":{
        "account": "sample",
        "email":"sample@mail.com"
    },
    "subscribes":[
        {
            "board":"gossiping",
            "keywords":["問卦","爆卦","公告"]
        },
        {
            "board":"lol",
            "keywords":["閒聊"]
        }
    ]
}
```

* PUT /users/[account]

```json
{
    "profile":{
        "account": "sample",
        "email":"sample@mail.com"
    },
    "subscribes":[]
}
```

## Credits

### Real Life

Rose Li, Aries Huang, Scott Kao, Amy Li

### Ptt

DMM, oas, bestpika, Zero0910, lucky0509, wbreeze, chang0206, lindo0130, hungys, gyman7788, tooilxui, myamyakoko, whkuo