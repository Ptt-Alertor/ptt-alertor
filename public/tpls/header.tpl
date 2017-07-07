{{define "header"}}
<div class="header">
    <nav>
        <ul class="nav nav-pills pull-right">
            <li role="presentation"><a href="/top">TOP 100</a></li>
            {{if eq .URI "line"}}
            <li role="presentation"><a href="/messenger">Messenger</a></li>
            {{else}}
            <li role="presentation"><a href="/line">LINE</a></li>
            {{end}}
        </ul>
    </nav>
    <h3 class="text-muted">Ptt Alertor</h3>
</div>
<hr>
{{end}}