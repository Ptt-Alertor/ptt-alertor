{{define "header"}}
<div class="header">
    <nav>
        <ul class="nav nav-pills pull-right">
            <li role="presentation"><a href="/top">TOP 100</a></li>

            <li class="hidden-xs" role="presentation"><a href="/line">LINE</a></li>
            <li class="hidden-xs" role="presentation"><a href="/messenger">Messenger</a></li>
            <li class="hidden-xs" role="presentation"><a href="/telegram">Telegram</a></li>

            <li role="presentation" class="dropdown visible-xs-block">
                <a href="#" class="dropdown-toggle" type="button" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                    Platform<span class="caret"></span>
                </a>
                <ul class="dropdown-menu dropdown-menu-right">
                    <li><a href="/line">LINE</a></li>
                    <li><a href="/messenger">Messenger</a></li>
                    <li><a href="/telegram">Telegram</a></li>
                </ul>
            </li>
        </ul>
    </nav>
    <h3 class="text-muted">Ptt Alertor</h3>
</div>
<hr>
{{end}}