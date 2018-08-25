package main

import (
	"html/template"
	"os"
)

type html struct {
	resultSet
}

func (out *html) print() {
	t := template.New("a")
	t.Parse(templateStr)
	t.Execute(os.Stdout, map[string]interface{}{"data": out.searches})
}

const templateStr string = `
<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <title>template</title>
    <style>
        body {
            width: 1260px;
            margin: auto;
        }

        nav {
			width: 1260px;
            font-size: 2em;
            position: fixed;
            padding: 5px;
            background-color: #cddc39;
        }

        nav:not(:hover) {
            height: 1.3em;
            overflow-y: hidden;
        }

        .nav {
            color: white;
        }

        nav a {
            text-decoration: none;
        }

        nav a:visited {
            color: white;
            text-decoration: none;
        }

        .nav+.nav:before {
            content: "|";
        }

        .entry {
            display: inline-block;
            width: 120px;
            height: 200px;
        }

        img {
            max-width: 100px;
            max-height: 100px;
        }

        .info {
            display: inline;
        }

        .info p {
            width: 100px;
            display: inline-block;
            margin: 0;
        }

        .searches {
            padding: 4em 5px 30px 5px;
        }

        h1 {
            width: 400px;
            margin: 0;
            padding: 0 0 10px 0;
            color: #9c27b0;
        }
    </style>
</head>

<body>
    <nav>
        {{ range .data }}
        <a class='nav' href='#{{ .Keyword }}'>
            {{ .Keyword }}
        </a>
        {{ end }}
    </nav>
    {{ range .data }}
    <div id='{{ .Keyword }}' class='searches'>
        <h1 class='keyword'> {{ .Keyword }} </h1>
        <div class='entries'>
            {{ range .Entries }}
            <a href='{{ .Href }}' target='_blank'>
                <div class='entry'>
                    <img src='{{ .Img }}' />
                    <div class='info'>
                        <p class='time'> {{ .Time }} </p>
                        <p class='title'> {{ .Title }} </p>
                        <p class='price'> {{ .Price }} </p>
                    </div>
                </div>
            </a>
            {{ end }}
        </div>

    </div>
    {{ end }}
</body>

</html>
`
