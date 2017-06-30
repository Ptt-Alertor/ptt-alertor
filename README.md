# Ptt-Alertor

<img align="right" src="https://raw.githubusercontent.com/liam-lai/ptt-alertor/master/logo.jpg">

[![Build Status](https://travis-ci.org/liam-lai/ptt-alertor.svg?branch=master)](https://travis-ci.org/liam-lai/ptt-alertor)
[![codecov](https://codecov.io/gh/liam-lai/ptt-alertor/branch/master/graph/badge.svg)](https://codecov.io/gh/liam-lai/ptt-alertor)
[![Go Report Card](https://goreportcard.com/badge/github.com/liam-lai/ptt-alertor)](https://goreportcard.com/report/github.com/liam-lai/ptt-alertor)
[![StackShare](https://img.shields.io/badge/tech-stack-0690fa.svg?style=flat)](https://stackshare.io/ptt-alertor/ptt-alertor)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## API

### Board

* GET /boards

* GET /boards/[board name]/articles

### User (Auth)

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

### Ptt

DMM, oas, bestpika, Zero0910, lucky0509, wbreeze, chang0206, lindo0130, hungys, gyman7788, tooilxui