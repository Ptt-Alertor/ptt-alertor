; win+s 開始測試
#s::
    ; 簡單指令正面
    Send, 指令{Enter}
    Sleep, 2000

    Send, 清單{Enter}
    Sleep, 2000

    ; 簡單指令反面
    Send, 指{Enter}
    Sleep, 2000

    ; 新增正面
    Send, 新增 gossiping 逗號結尾,{Enter}
    Sleep, 2000

    Send, 新增 gossiping ，逗號開頭{Enter}
    Sleep, 2000

    Send, 新增 gossiping 卦{Enter}
    Sleep, 2000

    ; 新增單項
    Send, 新增 gossiping 問卦{Enter}
    Sleep, 2000

    ; 新增混合有空白
    Send, 新增 gossiping 問卦,爆卦, ＦＢ{Enter}
    Sleep, 2000

    ; 新增混合逗號
    Send, 新增 gossiping 新聞,公告，協尋{Enter}
    Sleep, 2000

    ; 新增多版多關鍵字
    Send, 新增 lol,nba, baseball，tennis 樂透，閒聊{Enter}
    Sleep, 2000

    Send, 清單{Enter}
    Sleep, 2000

    ; 新增反面
    ; 新增版名錯誤
    Send, 新增 gossping 問卦{Enter}
    Sleep, 2000

    ; 新增版名逗號開頭
    Send, 新增 ,gossiping 逗號開頭{Enter}
    Sleep, 2000

    ; 新增版名逗號結尾
    Send, 新增 gossiping， 逗號結尾{Enter}
    Sleep, 2000

    ; 新增版名中間空白字元
    Send, 新增 gossiping， nba 問卦{Enter}
    Sleep, 2000

    Send, 清單{Enter}
    Sleep, 2000

    ; 刪除反面
    ; 刪除版名錯誤
    Send, 刪除 gossping 問卦{Enter}
    Sleep, 2000

    ; 刪除版名逗號開頭
    Send, 刪除 ，gossiping 逗號開頭{Enter}
    Sleep, 2000

    ; 刪除版名逗號結尾
    Send, 刪除 gossiping, 逗號結尾{Enter}
    Sleep, 2000

    ; 刪除版名中間空白字元
    Send, 刪除 gossiping， nba 問卦{Enter}
    Sleep, 2000

    Send, 清單{Enter}
    Sleep, 2000

    ; 刪除正面
    ; 刪除多版多關鍵字
    Send, 刪除 lol,nba, baseball，tennis 樂透，閒聊{Enter}
    Sleep, 2000

    Send, 刪除 gossiping 新聞， 公告{Enter}
    Sleep, 2000

    Send, 刪除 gossiping 爆卦,ＦＢ{Enter}
    Sleep, 2000

    Send, 刪除 gossiping 協尋{Enter}
    Sleep, 2000

    Send, 刪除 gossiping 逗號結尾,{Enter}
    Sleep, 2000

    Send, 刪除 gossiping ，逗號開頭{Enter}
    Sleep, 2000

    Send, 刪除 gossiping 卦{Enter}
    Sleep, 2000

    Send, 清單{Enter}
    Sleep, 2000

Escape::
ExitApp
Return