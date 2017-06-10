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
    ; 新增單項
    Send, 新增 gossiping 問卦{Enter}
    Sleep, 2000

    ; 新增混合有空白
    Send, 新增 gossiping 問卦,爆卦, ＦＢ{Enter}
    Sleep, 2000

    ; 新增混合逗號
    Send, 新增 gossiping 新聞,公告，協尋{Enter}
    Sleep, 2000

    Send, 清單{Enter}
    Sleep, 2000

    ; 新增反面
    ; 新增版名錯誤
    Send, 新增 gossping 問卦{Enter}
    Sleep, 2000

    ; 新增多個版名
    Send 新增 gossiping,lol 問卦{Enter}
    Sleep, 2000

    ; 新增逗號開頭
    Send 新增 gossiping ,逗號開頭{Enter}
    Sleep, 2000

    ; 新增逗號結尾
    Send 新增 gossiping 逗號結尾，{Enter}
    Sleep, 2000

    Send, 清單{Enter}
    Sleep, 2000

    ; 刪除反面
    ; 刪除版名錯誤
    Send, 刪除 gossping 問卦{Enter}
    Sleep, 2000

    ; 刪除多個版名
    Send 刪除 gossiping,lol 問卦{Enter}
    Sleep, 2000

    Send, 清單{Enter}
    Sleep, 2000

    ; 刪除正面
    Send, 刪除 gossiping 新聞， 公告{Enter}
    Sleep, 2000

    Send, 刪除 gossiping 爆卦,ＦＢ{Enter}
    Sleep, 2000

    Send, 刪除 gossiping 協尋{Enter}
    Sleep, 2000

    Send, 清單{Enter}
    Sleep, 2000

Escape::
ExitApp
Return