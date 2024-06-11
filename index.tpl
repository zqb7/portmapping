<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>portmapping</title>
    <script type="text/javascript" src="https://cdn.jsdelivr.net/npm/handsontable/dist/handsontable.full.min.js"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/handsontable/dist/handsontable.full.min.css" />
</head>
<style>
    body {
        margin: 0 0;
        text-align: center;
    }
    #main {
        margin: 15% auto;
        width: 80%;
        text-align: center;
    }
</style>
<body>
    <div id="main">
        <div id="myTable"></div>
    </div>
</body>
<script>
    const container = document.querySelector('#myTable');
    const hot = new Handsontable(container, {
        data: [
            {{range .}}
            ['{{.Port}}','{{.Network}}','{{.Status}}','{{.TargetIP}}','{{.TargetPort}}','{{.ClientNumber}}','{{.Desc}}'],
            {{end }}
        ],
        rowHeaders: true,
        colHeaders: ['本地端口', '网络类型', '是否启用', '远程主机', '远程端口','连接数', '说明'],
        width: '100%',
        height: 'auto',
        columns: [
            {"data":0, "type":"numeric",readOnly: true},
            {"data":1, "type":"text",readOnly: true},
            {"data":2, "type":"checkbox",className: "htCenter"},
            {"data":3, "type":"text"},
            {"data":4, "type":"numeric"},
            {"data":5, "type":"numeric",readOnly: true},
            {"data":6, "type":"text"},
        ],
        licenseKey: 'non-commercial-and-evaluation',
    });
    hot.addHook('afterChange', (row, e) => {
        console.log(row,e)
        var req = new XMLHttpRequest();
        var url = `/action?index=${row[0][0]}`
        var event = ""
        var desc = ""
        if (row[0][1] == 2){
            if (row[0][3] == true){
                event ="start"
            }else {
                event ="stop"
            }
        } else if (row[0][1] == 6) {
            if (row[0][2] != row[0][3]) {
                event ="desc"
                desc = row[0][3]
            }
        }
        url+=`&event=${event}`
        if (event =="desc" ){
            url +=`&v=${row[0][3]}`
        }
        if (event != ""){
            req.open("GET", url,true)
            req.onreadystatechange=function(){

            }
            req.send()
        }

    })
</script>
</html>