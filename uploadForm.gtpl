<html>
<head>
    <title>Upload file</title>
</head>
<body>
<form enctype="multipart/form-data" action="http://127.0.0.1:8080/upload/{{.token}}" method="post">
    Your Firstname:
    <input type="text" name="userFirstname"><br>
    Your Lastname:
    <input type="text" name="userLastname"><br>
    Your email:
    <input type="text" name="userEmail"><br>
    Details about upload:
    <input type="text" name="userReason"><br>
    <input type="file" name="uploadfile" />
    //<input type="hidden" name="token" value="{{.token}}"/>
    <input type="submit" value="upload" />
</form>
</body>
</html>