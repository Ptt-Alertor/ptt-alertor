{{define "counter"}}
    {{if .Count}}
    <div id="counter-board" class="well well-sm">
        <i class="fa fa-bell fa-fw fa-2x" aria-hidden="true"></i>
        <span>&nbsp;已送出</span>
        <span id="counter">
            {{range .Count}}
                {{if eq . ","}}
                <span>,</span>
                {{else}}
                <span class="label label-default">{{.}}</span>
                {{end}}
            {{end}}
        </span>
        <span>則通知</span>
    </div>
    {{end}}
{{end}}