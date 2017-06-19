# \[推薦\][自製] PTT Alertor PTT新文章即時通知

[#1PF_IETE@EZsoft](!https://www.ptt.cc/bbs/EZsoft/M.1497363598.A.74E.html)

## 軟體名稱：

PTT Alertor

## 軟體資訊：

LINE Bot, LINE Notify, Facebook Messenger Bot

## 軟體功能：

新增看板及作者或關鍵字，即時通知最新文章。
可一次新增多看板多關鍵字或多作者，目前設定為 10 秒通知一次。

Line Demo:
https://media.giphy.com/media/l0Iy28oboQbSw6Cn6/giphy.gif

Messenger Demo:
https://media.giphy.com/media/3ohzdF6vidM6I49lQs/giphy.gif

## 軟體特色：

有常常錯過版上樂透的經驗嗎？想第一時間知道八卦版的最新消息？想買東西？搶打工？

## 實用範例：

* LINE：

  新增 Instant_Mess 抽獎,免費

  新增 Lifeismoney line

* Stock：

  新增 stock 標的

  新增 tech_job 新聞

* 賭徒：

  新增 lol,nba,baseball,tennis 樂透

* 資訊恐慌：

  新增 gossiping 爆卦

* 找打工、工作：

  新增 part-time 全國,台中

  新增 tech_job,soft_job 徵才

* 買東西、合購：

  新增 buytogether 褲，裙，鞋，外套

  新增 nb-shopping,macshop macbook

  新增 drama-ticket 售&杰倫

* 追星：

  新增 gossiping,movie 金城武，結衣

* 追蹤作者：

  新增作者 gossiping ffaarr, mayaman

## 進階比對

* 同時出現: '&'

  新增 drama-ticket 售&杰倫
  (標題同時出現售和杰倫的才會通知)

* 除了B以外的所有文章: '!'

  新增 gossiping !問卦
  (除了問卦以外的文章全部通知)

* 出現A裡的文章排除B: '&!'

  新增 gossiping 柯文哲&!閒聊
  (標題出現柯文哲且不是問卦的才會通知)

* 正規表示式: 'regexp:'

  新增 hardwaresale regexp:\\[賣/(台中|台北)/.\*\\]?\(RAM|ram\)+.*
  (看得懂的.. 就看得懂 Orz)

## 下載連結：

Messenger: http://m.me/pttalertor
LINE: https://line.me/R/ti/p/%40vxl5146r

## 官方網站：

http://pttalertor.dinolai.com

## 粉絲團：

http://facebook.com/pttalertor

## 使用感想：

身家都賭光了，但是又靠著在某些版第一時間幫到人又賺回來了。

## 附註：

一開始是因為我和我弟都喜歡賭 LoL 的樂透，但常常忘記跟到，所以有了這個專案。

本來只想追樂透後來發現其實很多版可以用，一直改良便分享給大家試用看看。

## 常見問題：

* 為什麼要分成 LINE Bot 和 LINE Notify，不能合在一起嗎？

  因為若要使用 LINE Bot 內建的推播功能的話，要付一筆不小的月租費。

  LINE Notify 是免費的推播服務，嘗試看看行得通就將兩個結合在一起了。

* 會取得使用者 Facebook or LINE 個人資料嗎？

  系統只會存 Facebook or LINE 回傳給我的專屬於此 Bot 的 id，用來發送訊息。

  此 id 也無法直接連至你的 Facebook or LINE 頁面，請不用擔心資料問題。

* 若不想再收到通知怎麼做？
  1. 可以刪除所有追蹤清單
  1. 封鎖 PTT Alertor 帳號。LINE Notify 也將取消推播。

若有問題與建議麻煩推文、粉絲團或LINE首頁留言，希望大家使用愉快，謝謝。