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
        margin: 15% 0 30% 30%;
        height: 400px;
    }
</style>
<body>
    <div id="main">
        <div id="myTable"></div>
    </div>
</body>
<script>
    function delRenderer(hotInstance, td, row, column, prop, value, cellProperties) {
        let btn = document.createElement('button')
        btn.innerHTML = 'del';
        td.innerText = '';
        td.appendChild(btn)
        return td;
    }

    function isEmptyRow(instance, row) {
        const rowData = instance.countRows();
        for (let i = 0, ilen = rowData.length; i < ilen; i++) {
            if (rowData[i] !== null) {
                return false;
            }
        }
        return true;
    }

    const templateValues = ['0', 'tcp', '1', '', '0','0',''];

    function defaultValueRenderer(instance, td, row, col, prop, value, cellProperties) {
        const args = arguments;
        if (args[5] === null && isEmptyRow(instance, row)) {
            args[5] = templateValues[col];
            td.style.color = '#999';
        } else {
            td.style.color = '';
        }
        Handsontable.renderers.TextRenderer.apply(this, args);
    }

    const container = document.querySelector('#myTable');
    const hot = new Handsontable(container, {
        data: [
            {{range .}}
            ['{{.Port}}','{{.Network}}','{{.Status}}','{{.TargetHost}}','{{.TargetPort}}','{{.ClientNumber}}','{{.Desc}}'],
            {{ end }}
        ],
        rowHeaders: true,
        colHeaders: ['本地端口', '网络类型', '是否启用', '远程主机', '远程端口','连接数','说明','action'],
        width: '100%',
        height: 'auto',
        autoWrapRow: true,
        autoWrapCol: true,
        columns: [
            {"data":0, "type":"numeric",readOnly: true},
            {"data":1, "type":"text", readOnly: true},
            {"data":2, "type":"checkbox",className: "htCenter"},
            {"data":3, "type":"text"},
            {"data":4, "type":"numeric"},
            {"data":5, "type":"numeric",readOnly: true},
            {"data":6, "type":"text"},
            {"data":7, renderer: delRenderer},
        ],
        rowHeights: 30,
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
    hot.addHook('afterOnCellMouseDown', (e, coords, TD)=>{
        if (coords.col != 7 ){
            return
        }
        var req = new XMLHttpRequest();
        var url = `/action?index=${coords.row}&event=del`
        req.open("GET", url,true)
        req.onreadystatechange=function(){
            if (req.status==200){
                hot.alter("remove_row",coords.row)
            }
        }
        req.send()
    })
</script>
</html>