package main

var TEMPLATE_INDEX = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>GitLab build status</title>
    <link href="//unpkg.com/purecss@0.6.2/build/base-min.css" rel="stylesheet">
    <link href="//unpkg.com/purecss@0.6.2/build/grids-min.css" rel="stylesheet">
    <link href="//unpkg.com/purecss@0.6.2/build/grids-responsive-min.css" rel="stylesheet">
    <style>
        html, body { height: 100%;}
        body {font-size: 18px; color: #d8d9da; background-color: #1f1d1d;}
        .pure-g {margin: 1.5rem 1.5rem 0;}
        .button { display: inline-block; height: 1.3rem; padding: 0 .4rem; color: #555; text-align: center; font-size: .55rem; font-weight: 200; line-height: 1.2rem; text-transform: uppercase; text-decoration: none; white-space: nowrap; margin:.03rem; background-color: transparent; border-radius: 7px; border: 1px solid #bbb; cursor: pointer; box-sizing: border-box; }
        .button-success,.button-success:visited {color: #7fb36d; border-color: #7fb36d;}
        .button-running,.button-running:visited {color: #6ed0e0; border-color: #6ed0e0;}
        .button-pending,.button-pending:visited {color: #ebb939; border-color: #ebb939;}
        .button-created,.button-created:visited {display: none;}
        .button-failed,.button-failed:visited {color: #fff; background-color: #e34d42; border-color: #e34d42;}
        .button-canceled,.button-canceled:visited {color: #d8d9da; border-color: #d8d9da;}
        .meta {font-weight:200; font-size: .8rem;}
        a,a:hover,a:active,a:visited,.button:hover {color: #95999b; text-decoration: none;}
        .statuses {padding-left: 2rem;}
        .legend { text-align: center; margin-bottom: 2rem; }
        .legend>.button { width: 40%; }
        body > .container {height: auto; min-height:100%; box-sizing: border-box; padding-bottom: 5rem; }
        .container>.pure-g {margin-top: 0;}
        .footer { height: 5rem; margin-top: -5rem; }
    </style>
    <script src="//ajax.googleapis.com/ajax/libs/jquery/3.1.0/jquery.min.js"></script>
</head>
<body>
    <div class="container">
        <div class="pure-g">
        {{ range . }}
            <div id='{{ .Id }}' class="pure-u-1 pure-u-lg-1-2 pure-u-xl-1-3 project_card">
                <div class="pure-g">
                    <div class='pure-u-1-3'>
                        <b>{{ .Name }}</b>
                    </div>
                    <div class='pure-u-2-3'>
                        <style>
                            .spinner { margin: auto; width: 35px; height: 30px;}
                            .spinner > div { margin: 1px; background-color: #ccc; height: 100%; width: 5px; display: inline-block; -webkit-animation: sk-stretchdelay 1.2s infinite ease-in-out; animation: sk-stretchdelay 1.2s infinite ease-in-out; }
                            .spinner .rect2 { -webkit-animation-delay: -1.1s; animation-delay: -1.1s; }
                            .spinner .rect3 { -webkit-animation-delay: -1.0s; animation-delay: -1.0s; }
                            .spinner .rect4 { -webkit-animation-delay: -0.9s; animation-delay: -0.9s; }
                            .spinner .rect5 { -webkit-animation-delay: -0.8s; animation-delay: -0.8s; }
                            @-webkit-keyframes sk-stretchdelay { 0%, 40%, 100% { -webkit-transform: scaleY(0.4) }  20% { -webkit-transform: scaleY(1.0) } }
                            @keyframes sk-stretchdelay { 0%, 40%, 100% { transform: scaleY(0.4); -webkit-transform: scaleY(0.4);}  20% { transform: scaleY(1.0); -webkit-transform: scaleY(1.0); }}
                        </style>
                        <div class="spinner"><div class="rect1"></div><div class="rect2"></div><div class="rect3"></div><div class="rect4"></div><div class="rect5"></div></div>
                    </div>
                </div>
            </div>
        {{end}}
        </div>
    </div>
    <div class="pure-g footer">
        <div class="pure-u-1">
            <hr noshade>
        </div>
        <div class="pure-u-1-5 legend">
            <a class='button button-success'>success</a>
        </div>
        <div class="pure-u-1-5 legend">
            <a class='button button-running'>running</a>
        </div>
        <div class="pure-u-1-5 legend">
            <a class='button button-pending'>pending</a>
        </div>
        <div class="pure-u-1-5 legend">
            <a class='button button-failed'>failed</a>
        </div>
        <div class="pure-u-1-5 legend">
            <a class='button button-canceled'>canceled</a>
        </div>
    </div>
    <script type="text/javascript">
        function updateStatus() {
            $(".project_card").each(function(){
                var id = this.id;
                $.get("/status/?id=" + id).done(function( data ) {
                    $("#" + id).html(data);
                    $("#" + id).show();
                }).fail(function() {
                    $("#" + id).hide();
                });
            });
            setTimeout(updateStatus, 10000);
        };

        updateStatus();
        setTimeout(function(){location=""}, 10*60*1000);

    </script>
</body>
</html>
`

var TEMPLATE_STATUS = `
<div class="pure-g">
    <div class='pure-u-1-3'>
        <b>{{ .Project.Name }}</b> {{ with .Coverage }}= {{ . }}%{{end}}<br>
        <p class="meta"><a target="_blank" href="{{ .Url }}/{{ .Project.Name }}/commit/{{ .Commit.Id }}" >#{{ .Commit.Id | Short }}: {{ .Commit.Message | Title }}</a><br>
        <a target="_blank" href="{{ .Url }}/u/{{ .Commit.Author_name }}">@{{ .Commit.Author_name }}</a>
        </p>
    </div>
    <div class='pure-u-2-3'>
        <div class="statuses">
        {{ range .Builds }}
            <a target="_blank" class='button button-{{ .Status }}' href="{{ $.Url }}/{{ $.Project.Name }}/builds/{{ .Id }}">{{ .Name }}</a>
        {{end}}
        </div>
    </div>
</div>
`
