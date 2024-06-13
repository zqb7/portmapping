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
    }
    #main {
        margin: 15% 0 15% 30%;
        height: 400px;
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
            ['{{.Port}}','{{.Network}}','{{.Status}}','{{.TargetHost}}','{{.TargetPort}}','{{.ClientNumber}}','{{.Desc}}'],
            {{ end }}
        ],
        rowHeaders: true,
        colHeaders: ['本地端口', '网络类型', '是否启用', '远程主机', '远程端口','连接数','说明'],
        autoWrapRow: true,
        autoWrapCol: true,
        columns: [
            {"data":0, "type":"numeric"},
            {"data":1, "type":"dropdown", source:["tcp"]},
            {"data":2, "type":"checkbox"},
            {"data":3, "type":"text"},
            {"data":4, "type":"numeric"},
            {"data":5, "type":"numeric",readOnly: true},
            {"data":6, "type":"text"},
        ],
        contextMenu: ['row_below', 'remove_row'],
        licenseKey: 'non-commercial-and-evaluation',
    });
    hot.addHook('afterChange', (row, e) => {
        console.log(row,e)
        let rowIndex = row[0][0]
        let before_value = row[0][2]
        let after_value = row[0][3]
        var req = new XMLHttpRequest();
        var url = `/action?index=${rowIndex}`
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
    hot.addHook("afterRemoveRow",(index, amount,physicalRows, source)=>{
        console.log(index, amount,physicalRows, source)
        var req = new XMLHttpRequest();
        req.open("GET", `/action?index=${index}&event=del`,true)
        req.send()
    })
</script>
</html>