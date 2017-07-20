{{define "script"}}
<!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
<!-- Include all compiled plugins (below), or include individual files as needed -->
<!-- Latest compiled and minified JavaScript -->
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa"
    crossorigin="anonymous"></script>
    {{if .Count}}
    <script src="https://cdn.jsdelivr.net/countupjs/1.8.5/countUp.min.js"></script>
    <script>
        $(function () {
            var url = "{{.WSHost}}/ws";
            var ws = new WebSocket(url);
            var counterUps = [];
            var spans = document.getElementById("counter").querySelectorAll(".label");
            spans.forEach(function (span) {
                var countup = new CountUp(span, 0, parseInt(span.textContent));
                countup.start();
                counterUps.push(countup);
            })
            ws.onmessage = function (event) {
                var digits = event.data.split("");
                counterUps.forEach(function (counterUp, i) {
                    counterUp.update(digits[i]);
                })
            }
        })
    </script>
    {{end}}
{{end}}