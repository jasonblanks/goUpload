

<html>
<head>
    <title>Register Upload Request</title><br>
    <table style="width:100%">
        <TR>
            <TH>HASH</TH>
            <TH>EXPIRE DATETIME</TH>
            <TH>ACCESSED</TH>
        </TR>
    {{range .}}
    <TR>
            <TD>{{.Sha1}}</TD>
            <td>{{.ExpireDate}}</td>
            <td>{{.Accessed}}</td>
    {{end}}
    </TR>
</head>
<body>