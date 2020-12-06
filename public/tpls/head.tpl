{{define "head"}}
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Ptt Alertor - Ptt 即時文章通知，追蹤看板推文數、作者、關鍵字</title>
    <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
    <meta name="description" content="訂閱看板推文數作者關鍵字，即時通知 Ptt 最新文章">
    <meta name="keywords" content="Ptt, Ptt Alarm, Ptt Alert, FB Bot, Messenger Bot, Ptt Notification, Ptt 通知器, Ptt 追蹤">
    <meta name="author" content="Dino Lai, Liam Lai">
    <meta name="google-site-verification" content="oj4eRvcsK1aK3rqqBxn6piDte5-2sG9neQqcAnd8gTo" />
    <!-- Facebook Open Graph -->
    <meta property="og:url" content="https://pttalertor.dinolai.com/{{.URI}}" />
    <meta property="og:type" content="product" />
    <meta property="og:title" content="Ptt Alertor - Ptt 即時文章通知，追蹤看板推文數、作者、關鍵字" />
    <meta property="og:description" content="訂閱看板推文數作者關鍵字，即時通知 Ptt 最新文章" />
    <meta property="og:image" content="https://s3-us-west-2.amazonaws.com/ptt-alertor-2020-bucket/assets/alarmP.png" />

    <link rel="icon" type="image/png" href="https://s3-us-west-2.amazonaws.com/ptt-alertor-2020-bucket/assets/alarmP32x32.png">
    <link rel="apple-touch-icon" type="image/png" href="https://s3-us-west-2.amazonaws.com/ptt-alertor-2020-bucket/assets/alarmP.png">
    <meta name="msapplication-TileImage" content="https://s3-us-west-2.amazonaws.com/ptt-alertor-2020-bucket/assets/alarmP.png">
    <meta name="msapplication-TileColor" content="#FFFFFF"/>
    <link rel="alternate" href="https://pttalertor.dinolai.com/{{.URI}}" hreflang="zh-hant" />

    <!-- Bootstrap -->
    <!-- Latest compiled and minified CSS -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u"
        crossorigin="anonymous">

    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.3/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
    <link rel="stylesheet" href="assets/css/index.css">
    <script src="assets/js/google-analytics.js"></script>

</head>
{{end}}