; win+s 開始測試
#s::
    ; 簡單指令正面
    Send, 指令{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    Send, 排行{Enter}
    Sleep, 1500

    ; 簡單指令反面
    Send, 指{Enter}
    Sleep, 1500

    ; 新增正面
    Send, 新增 gossiping 逗號結尾,{Enter}
    Sleep, 1500

    Send, 新增 gossiping ，逗號開頭{Enter}
    Sleep, 1500

    Send, 新增 gossiping 卦{Enter}
    Sleep, 1500

    ; 新增單項
    Send, 新增 gossiping 問卦{Enter}
    Sleep, 1500

    ; 新增內含*符號
    Send, 新增 gossiping 內含,*,星號{Enter}
    Sleep, 1500

    ; 新增混合有空白
    Send, 新增 gossiping 問卦,爆卦, ＦＢ{Enter}
    Sleep, 1500

    ; 新增混合逗號
    Send, 新增 gossiping 新聞,公告，協尋{Enter}
    Sleep, 1500

    ; 新增多版多關鍵字
    Send, 新增 lol,nba,baseball，tennis 樂透，閒聊{Enter}
    Sleep, 1500

    ; 新增同時出現的關鍵字
    Send, 新增 gossiping 新聞&柯文哲{Enter}
    Sleep, 1500

    ; 新增不要的關鍵字
    Send, 新增 gossiping !問卦{Enter}
    Sleep, 1500

    ; 新增有此關鍵字但是除了不要的
    Send, 新增 gossiping 柯文哲&!問卦{Enter}
    Sleep, 1500

    ; 新增 regexp
    Send, 新增 gossiping regexp:{^}\[問卦\]{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    ; 新增反面
    ; 新增版名錯誤
    Send, 新增 gossping 問卦{Enter}
    Sleep, 1500

    ; 新增版名逗號開頭
    Send, 新增 ,gossiping 逗號開頭{Enter}
    Sleep, 1500

    ; 新增版名逗號結尾
    Send, 新增 gossiping， 逗號結尾{Enter}
    Sleep, 1500

    ; 新增版名中間空白字元
    Send, 新增 gossiping， nba 問卦{Enter}
    Sleep, 1500

    ; 新增錯誤 regexp
    Send, 新增 gossiping regexp:{^}[{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    ; 刪除反面
    ; 刪除版名錯誤
    Send, 刪除 gossping 問卦{Enter}
    Sleep, 1500

    ; 刪除版名逗號開頭
    Send, 刪除 ，gossiping 逗號開頭{Enter}
    Sleep, 1500

    ; 刪除版名逗號結尾
    Send, 刪除 gossiping, 逗號結尾{Enter}
    Sleep, 1500

    ; 刪除版名中間空白字元
    Send, 刪除 gossiping， nba 問卦{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    ; 刪除正面
    ; 刪除多版多關鍵字
    Send, 刪除 lol,nba,baseball，tennis 樂透，閒聊{Enter}
    Sleep, 1500

    Send, 刪除 gossiping 新聞， 公告{Enter}
    Sleep, 1500

    Send, 刪除 gossiping 爆卦,ＦＢ{Enter}
    Sleep, 1500

    Send, 刪除 gossiping 協尋{Enter}
    Sleep, 1500

    Send, 刪除 gossiping 逗號結尾,{Enter}
    Sleep, 1500

    Send, 刪除 gossiping ，逗號開頭{Enter}
    Sleep, 1500

    Send, 刪除 gossiping 卦{Enter}
    Sleep, 1500

    ; 刪除同時出現的關鍵字
    Send, 刪除 gossiping 新聞&柯文哲{Enter}
    Sleep, 1500

    ; 刪除不要的關鍵字
    Send, 刪除 gossiping !問卦{Enter}
    Sleep, 1500

    ; 刪除有此關鍵字但是除了不要的
    Send, 刪除 gossiping 柯文哲&!問卦{Enter}
    Sleep, 1500

    ; 刪除 regexp
    Send, 刪除 gossiping regexp:{^}\[問卦\]{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    ; 新增作者
    Send, 新增作者 gossiping ffaarr{Enter}
    Sleep, 1500

    Send, 新增作者 gossiping ffaarr,obov{Enter}
    Sleep, 1500

    Send, 新增作者 lol,boy-girl sumade,mrp{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    ; 新增反面
    Send, 新增作者 gossiping obov{Enter}
    Sleep, 1500

    Send, 新增作者 gossping ffaarr{Enter}
    Sleep, 1500

    Send, 新增作者 ,gossiping ffaarr{Enter}
    Sleep, 1500

    Send, 新增作者 gossiping, ffaarr{Enter}
    Sleep, 1500

    Send, 新增作者 stock 抄底王{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    ; 刪除作者
    ; 刪除反面
    Send, 刪除作者 gossping ffaarr{Enter}
    Sleep, 1500

    ; 刪除逗號前
    Send, 刪除作者 ,gossiping ffaarr{Enter}
    Sleep, 1500

    ; 刪除刪除逗號後
    Send, 刪除作者 gossiping, ffaarr{Enter}
    Sleep, 1500

    Send, 刪除作者 stock 抄底王{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    ; 刪除正面
    ; 刪除作者
    Send, 刪除作者 gossiping ffaarr{Enter}
    Sleep, 1500

    Send, 刪除作者 gossiping ffaarr,obov{Enter}
    Sleep, 1500

    Send, 刪除作者 lol,boy-girl sumade,mrp{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    ; 推文數
    ; 新增正面
    Send, 新增推文數 gossiping 10{Enter}
    Sleep, 1500

    ; 多板
    Send, 新增推文數 gossiping,beauty 10{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    ; 新增反面
    Send, 新增推文數 gossping 101{Enter}
    Sleep, 1500

    Send, 新增推文數 gossping -10{Enter}
    Sleep, 1500

    Send, 新增推文數 gossping abc{Enter}
    Sleep, 1500

    Send, 新增推文數 gossping 10{Enter}
    Sleep, 1500

    Send, 新增推文數 ,gossiping 10{Enter}
    Sleep, 1500

    Send, 新增推文數 gossiping, 10{Enter}
    Sleep, 1500

    Send, 清單 {Enter}
    Sleep, 1500

    ; 推文數
    ; 新增正面
    Send, 新增噓文數 gossiping 10{Enter}
    Sleep, 1500

    ; 多板
    Send, 新增噓文數 gossiping,beauty 10{Enter}
    Sleep, 1500

    Send, 清單 {Enter}
    Sleep, 1500

    ; 新增反面
    Send, 新增噓文數 gossping 101{Enter}
    Sleep, 1500

    Send, 新增噓文數 gossping -10{Enter}
    Sleep, 1500

    Send, 新增噓文數 gossping abc{Enter}
    Sleep, 1500

    Send, 新增噓文數 gossping 10{Enter}
    Sleep, 1500

    Send, 新增噓文數 ,gossiping 10{Enter}
    Sleep, 1500

    Send, 新增噓文數 gossiping, 10{Enter}
    Sleep, 1500

    Send, 清單 {Enter}
    Sleep, 1500

    ; 歸零
    Send, 新增推文數 gossiping 0{Enter}
    Sleep, 1500

    Send, 新增噓文數 gossiping 0{Enter}
    Sleep, 1500

    Send, 清單 {Enter}
    Sleep, 1500

    ; 推文
    ; 新增推文
    Send, 新增推文 https://www.ptt.cc/bbs/EZsoft/M.1497363598.A.74E.html{Enter}
    Sleep, 1500

    ; 新增反面
    Send, 新增推文 www.ptt.cc/bbs/EZsoft/M.1497363598.A.74E.html{Enter}
    Sleep, 1500

    Send, 新增推文 https://www.ptt.cc/bbs/EZsoft/M.1497363598.A.74E.html 賣，買{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

    ; 刪除推文
    Send, 刪除推文 https://www.ptt.cc/bbs/EZsoft/M.1497363598.A.74E.html{Enter}
    Sleep, 1500

    Send, 清單{Enter}
    Sleep, 1500

Escape::
ExitApp
Return

; Manual Testing

Send, 新增 lol,nba,tennis,baseball,gossiping 樂透，閒聊,戰況,姆斯,炸裂{Enter}
Send, 刪除 ** 樂透{Enter}
Send, 刪除 ** 閒聊,戰況{Enter}
Send, 刪除 lol *{Enter}
Send, 刪除 nba,gossping *{Enter}
Send, 刪除 ** *{Enter}

Send, 新增作者 lol,nba,tennis,baseball,gossiping chodino,obov,ffaarr,sumade,boyo{Enter}
Send, 刪除作者 ** chodino{Enter}
Send, 刪除作者 ** obov,ffaarr{Enter}
Send, 刪除作者 lol *{Enter}
Send, 刪除作者 nba,gossping *{Enter}
Send, 刪除作者 ** *{Enter}