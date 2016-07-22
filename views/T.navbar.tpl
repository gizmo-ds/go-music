{{define "navbar"}}</head>
<body>
<nav class="navbar navbar-inverse navbar-fixed-top" role="navigation">
   <div class="container">
      <div class="navbar-header">
         <button type="button" class="navbar-toggle" data-toggle="collapse"
            data-target="#example-navbar-collapse">
            <span class="sr-only">切换导航</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
         </button>
         <a class="navbar-brand" href="/">{{.Name}}</a>
      </div>
      <div class="collapse navbar-collapse" id="example-navbar-collapse">
         <ul class="nav navbar-nav">
            <li{{if .IsHome}} class="active"{{end}}><a href="/">首页</a></li>
            <li class="dropdown">
               <a class="dropdown-toggle" data-toggle="dropdown" role="button" >网易云音乐 <span class="caret"></span></a>
               <ul class="dropdown-menu">
                  <li><a href="/song">音乐</a></li>
                  <li><a href="/list">歌单</a></li>
                  <li><a href="/album">专辑</a></li>
               </ul>
            </li>
            <li{{if .IsXiami}} class="active"{{end}}><a href="/xiami">虾米音乐</a></li>
            <li{{if .IsKugou}} class="active"{{end}}><a href="/kugou">酷狗音乐</a></li>
            <li{{if .IsMessage}} class="active"{{end}}><a href="/message">留言板</a></li>
         </ul>
      </div>
   </div>
</nav>{{end}}
