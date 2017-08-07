#s::
    ; 正面
    Send, add -k ptt ezsoft{Enter}
    Sleep, 1500

    Send add -k ptt&alertor ezsoft{Enter}
    Sleep, 1500

    Send, add -a chodino ezsoft{Enter}
    Sleep, 1500

    Send, add -p 10 ezsoft{Enter}
    Sleep, 1500

    Send, add -b 10 ezsoft{Enter}
    Sleep, 1500

    Send, add -k ptt -a chodino -p 10 -b 10 ezsoft,gossiping{Enter}
    Sleep, 1500

    Send, list{Enter}
    Sleep, 1500

    ; 反面
    Send add -a 你好 ezsoft{Enter}
    Sleep, 1500

    Send add -k 123 ezezez{Enter}
    Sleep, 1500

    Send add -p 111 ezsoft{Enter}
    Sleep, 1500

    Send add -p abc ezsoft{Enter}
    Sleep, 1500

    Send add -b abc ezsoft{Enter}
    Sleep, 1500

    Send, list{Enter}
    Sleep, 1500

Escape::
ExitApp
Return